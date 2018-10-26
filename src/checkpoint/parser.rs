use super::types::{self, CheckPointFile};

use std::error;
use std::fmt;
use std::io::{self, BufRead, BufReader, Read};
use std::iter;
use std::num::{ParseFloatError, ParseIntError};
use std::result;

macro_rules! fixed_width_part {
    ((), $input:expr, $i:expr) => {{
        let _ = $input.by_ref().take($i).collect::<String>();
        ()
    }};
    (RawString, $input:expr, $i:expr) => {
        $input.by_ref().take($i).collect::<String>()
    };
    (String, $input:expr, $i:expr) => {
        $input.by_ref().take($i).collect::<String>().trim()
    };
    ($t:tt, $input:expr, $i:expr) => {
        $input
            .by_ref()
            .take($i)
            .collect::<String>()
            .trim()
            .parse::<$t>()?
    };
}

macro_rules! fixed_width {
    ($input:expr, $($i:expr => $t:tt),*) => {
        {
            let mut line = $input.chars();
            let result: Result<_> = Ok(($(fixed_width_part!($t, line, $i),)*));
            result
        }
    };
}

#[derive(Debug)]
pub enum Error {
    Io { error: io::Error },
    ParseInt { error: ParseIntError },
    ParseFloat { error: ParseFloatError },
    InvalidFormat(String),
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            Error::Io { error } => write!(f, "IO Error: {:}", error),
            Error::ParseInt { error } => write!(f, "ParseInt Error: {:}", error),
            Error::ParseFloat { error } => write!(f, "ParseFloat Error: {:}", error),
            Error::InvalidFormat(error) => write!(f, "Invalid Format: {:}", error),
        }
    }
}

impl error::Error for Error {
    fn cause(&self) -> Option<&dyn error::Error> {
        match self {
            Error::Io { error } => Some(error),
            Error::ParseInt { error } => Some(error),
            Error::ParseFloat { error } => Some(error),
            Error::InvalidFormat(_) => None,
        }
    }
}

impl From<ParseIntError> for Error {
    fn from(e: ParseIntError) -> Self {
        Error::ParseInt { error: e }
    }
}

impl From<ParseFloatError> for Error {
    fn from(e: ParseFloatError) -> Self {
        Error::ParseFloat { error: e }
    }
}

impl From<io::Error> for Error {
    fn from(e: io::Error) -> Self {
        Error::Io { error: e }
    }
}

pub type Result<T> = result::Result<T, Error>;

impl CheckPointFile {
    pub fn parse<R: Read>(inner: R) -> Result<CheckPointFile> {
        CPFParser::new(inner).parse()
    }
}

impl types::SplittedBond {
    fn parse(line: &str) -> Result<Self> {
        let (bda, baa) = fixed_width!(line, 5 => i32, 5 => i32)?;
        Ok(types::SplittedBond { bda: bda, baa: baa })
    }
}

struct CPFParser<R> {
    file: BufReader<R>,
    num_frags: usize,
    buffer: String,
}

impl<R: Read> CPFParser<R> {
    fn new(inner: R) -> Self {
        CPFParser {
            file: BufReader::new(inner),
            num_frags: 0,
            buffer: String::new(),
        }
    }

    fn parse(mut self) -> Result<CheckPointFile> {
        let (version, comment) = self.parse_version()?;
        let (num_atoms, num_frags) = fixed_width!(self.next_line()?, 5 => usize, 5 => usize)?;
        self.num_frags = num_frags;

        let mut atoms = Vec::with_capacity(num_atoms);
        let mut atom_indices: Vec<Vec<i32>> = vec![vec![]; num_frags];
        for i in 0..num_atoms {
            let atom = self.next_line().and_then(types::Atom::parse)?;
            atom_indices[atom.frag_index as usize - 1].push(i as i32);
            atoms.push(atom);
        }

        let num_electrons = self.parse_size16(num_frags)?;
        let num_attached = self.parse_size16(num_frags)?;
        let num_splitted_bonds = num_attached.iter().sum::<i32>();
        let splitted_bonds = (0..num_splitted_bonds)
            .map(|_| self.next_line().and_then(types::SplittedBond::parse))
            .collect::<Result<Vec<_>>>()?;

        let distances = (0..self.len_dimers())
            .map(|_| {
                self.next_line().and_then(|line| {
                    line.split_whitespace()
                        .skip(2)
                        .next()
                        .ok_or(Error::InvalidFormat("distance".to_string()))
                        .and_then(|s| s.parse::<f64>().map_err(Into::into))
                })
            }).collect::<Result<Vec<_>>>()?;

        let mut hf_dipoles = Vec::with_capacity(num_frags);
        let mut mp2_dipoles = Vec::with_capacity(num_frags);
        for _ in 0..num_frags {
            match self
                .next_line()?
                .split_whitespace()
                .map(|s| s.parse::<f64>().map_err(Into::into))
                .collect::<Result<Vec<_>>>()
            {
                Ok(ref v) if v.len() == 6 => {
                    hf_dipoles.push(types::Vec3::new(v[0], v[1], v[2]));
                    mp2_dipoles.push(types::Vec3::new(v[3], v[4], v[5]));
                }
                _ => Err(Error::InvalidFormat("dipole".to_string()))?,
            }
        }
        let level = self.next_line()?.trim().to_string();
        let electronic_state = self.next_line()?.trim().to_string();
        let basisset = self.next_line()?.trim().to_string();
        let approx_params = match self
            .next_line()?
            .split_whitespace()
            .map(|s| s.parse::<f64>().map_err(Into::into))
            .collect::<Result<Vec<_>>>()
        {
            Ok(ref v) if v.len() == 3 => Ok(types::ApproxParams {
                ao_population: v[0],
                point_charge: v[1],
                dimer_es: v[2],
            }),
            _ => Err(Error::InvalidFormat("approx".to_string())),
        }?;

        let nuclear_repulsion = self.next_line()?.trim().parse::<f64>()?;
        let total_electronic = self.next_line()?.trim().parse::<f64>()?;
        let total_energy = self.next_line()?.trim().parse::<f64>()?;

        let mut m_repls = Vec::with_capacity(num_frags);
        let mut m_eles = Vec::with_capacity(num_frags);
        let mut m_mp2s = Vec::with_capacity(num_frags);
        let mut m_mp3s = Vec::with_capacity(num_frags);
        let mut m_aos = Vec::with_capacity(num_frags);
        let mut m_mos = Vec::with_capacity(num_frags);
        for _ in 0..num_frags {
            let (m_repl, m_ele, m_mp2, m_mp3, m_ao, m_mo) = fixed_width!(self.next_line()?, 24 => f64, 24 => f64, 24 => f64, 24 => f64, 12 => i32, 12 => i32)?;
            m_repls.push(m_repl);
            m_eles.push(m_ele);
            m_mp2s.push(m_mp2);
            m_mp3s.push(m_mp3);
            m_aos.push(m_ao);
            m_mos.push(m_mo);
        }

        let mut d_repls = Vec::with_capacity(self.len_dimers());
        let mut d_eles = Vec::with_capacity(self.len_dimers());
        let mut d_ess = Vec::with_capacity(self.len_dimers());
        let mut d_mp2s = Vec::with_capacity(self.len_dimers());
        let mut d_scs_mp2s = Vec::with_capacity(self.len_dimers());
        let mut d_mp3s = Vec::with_capacity(self.len_dimers());
        let mut d_scs_mp3s = Vec::with_capacity(self.len_dimers());
        let mut d_hf_bsses = Vec::with_capacity(self.len_dimers());
        let mut d_mp2_bsses = Vec::with_capacity(self.len_dimers());
        let mut d_scs_mp2_bsses = Vec::with_capacity(self.len_dimers());
        let mut d_mp3_bsses = Vec::with_capacity(self.len_dimers());
        let mut d_scs_mp3_bsses = Vec::with_capacity(self.len_dimers());
        let mut d_exs = Vec::with_capacity(self.len_dimers());
        let mut d_cts = Vec::with_capacity(self.len_dimers());
        let mut d_qijs = Vec::with_capacity(self.len_dimers());
        for _ in 0..self.len_dimers() {
            let (
                d_repl,
                d_ele,
                d_es,
                d_mp2,
                d_scs_mp2,
                d_mp3,
                d_scs_mp3,
                d_hf_bsse,
                d_mp2_bsse,
                d_scs_mp2_bsse,
                d_mp3_bsse,
                d_scs_mp3_bsse,
                d_ex,
                d_ct,
                d_qij,
            ) = fixed_width!(self.next_line()?, 24 => f64, 24 => f64, 24 => f64, 24 => f64, 24 => f64, 24 => f64, 24 => f64, 24 => f64, 24 => f64, 24 => f64, 24 => f64, 24 => f64, 24 => f64, 24 => f64, 24 => f64)?;

            d_repls.push(d_repl);
            d_eles.push(d_ele);
            d_ess.push(d_es);
            d_mp2s.push(d_mp2);
            d_scs_mp2s.push(d_scs_mp2);
            d_mp3s.push(d_mp3);
            d_scs_mp3s.push(d_scs_mp3);
            d_hf_bsses.push(d_hf_bsse);
            d_mp2_bsses.push(d_mp2_bsse);
            d_scs_mp2_bsses.push(d_scs_mp2_bsse);
            d_mp3_bsses.push(d_mp3_bsse);
            d_scs_mp3_bsses.push(d_scs_mp3_bsse);
            d_exs.push(d_ex);
            d_cts.push(d_ct);
            d_qijs.push(d_qij);
        }

        Ok(CheckPointFile {
            version: version,
            comments: comment,
            num_atoms: num_atoms as i32,
            num_frags: num_frags as i32,
            atoms: atoms,
            monomers: types::Monomers {
                atom_indices: atom_indices,
                num_electrons: num_electrons,
                num_attached: num_attached,
                hf_dipoles: hf_dipoles,
                mp2_dipoles: mp2_dipoles,

                nuclear_repulsion: m_repls,
                electronic: m_eles,
                mp2: m_mp2s,
                mp3: m_mp3s,
                ao: m_aos,
                mo: m_mos,
            },
            splitted_bonds: splitted_bonds,
            dimers: types::Dimers {
                distances: types::DimerData::new(distances),
                nuclear_repulsion: types::DimerData::new(d_repls),
                electronic: types::DimerData::new(d_eles),
                electrostatic: types::DimerData::new(d_ess),
                mp2: types::DimerData::new(d_mp2s),
                scs_mp2: types::DimerData::new(d_scs_mp2s),
                mp3: types::DimerData::new(d_mp3s),
                scs_mp3: types::DimerData::new(d_scs_mp3s),
                hf_bsse: types::DimerData::new(d_hf_bsses),
                mp2_bsse: types::DimerData::new(d_mp2_bsses),
                scs_mp2_bsse: types::DimerData::new(d_scs_mp2_bsses),
                mp3_bsse: types::DimerData::new(d_mp3_bsses),
                scs_mp3_bsse: types::DimerData::new(d_scs_mp3_bsses),
                exchange_repulsion: types::DimerData::new(d_exs),
                charge_transfer: types::DimerData::new(d_cts),
                qij: types::DimerData::new(d_qijs),
            },
            level: level,
            electronic_state: electronic_state,
            basisset: basisset,
            approx_params: approx_params,
            system_energy: types::SystemEnergy {
                nuclear_repulsion: nuclear_repulsion,
                electronic: total_electronic,
                total: total_energy,
            },
        })
    }

    fn next_line<'a>(&'a mut self) -> Result<&'a str> {
        self.buffer.clear();
        match self.file.read_line(&mut self.buffer)? {
            0 => panic!("unexpected EOF"),
            _ => return Ok(&self.buffer),
        }
    }

    fn parse_version(&mut self) -> Result<(String, String)> {
        let line = self.next_line()?.chars();
        let mut line = line.skip(8);
        let version = line.by_ref().take_while(|c| c != &' ').collect();
        let comment = line.take_while(|c| c != &'\n' && c != &'\r').collect();

        Ok((version, comment))
    }

    fn parse_size16(&mut self, num_frags: usize) -> Result<Vec<i32>> {
        let mut result = Vec::with_capacity(num_frags);
        for size in iter::repeat(16)
            .take(num_frags / 16)
            .chain(vec![num_frags % 16])
        {
            let mut line = self.next_line()?.chars();
            for _ in 0..size {
                result.push(
                    line.by_ref()
                        .take(5)
                        .collect::<String>()
                        .trim()
                        .parse::<i32>()?,
                )
            }
        }

        Ok(result)
    }

    fn len_dimers(&self) -> usize {
        self.num_frags * (self.num_frags - 1) / 2
    }
}

impl types::Atom {
    fn parse(line: &str) -> Result<Self> {
        let (
            index,
            _,
            element,
            _,
            atom_type,
            _,
            res_name,
            _,
            res_index,
            _,
            frag_index,
            _,
            x,
            y,
            z,
            hf_mul,
            mp2_mul,
            hf_nbo,
            mp2_nbo,
            hf_esp,
            mp2_esp,
        ) = fixed_width!(line,
                5 => i32, 1 => (),
                2 => RawString, 1 => (),
                4 => RawString, 1 => (),
                3 => RawString, 1 => (),
                4 => i32, 1 => (),
                4 => i32, 1 => (),
                12 => f64, 12 => f64, 12 => f64,
                12 => f64, 12 => f64, 12 => f64,
                12 => f64, 12 => f64, 12 => f64
            )?;
        Ok(types::Atom {
            index: index,
            element: element,
            atom_type: atom_type,
            res_name: res_name,
            res_index: res_index,
            frag_index: frag_index,
            x: x,
            y: y,
            z: z,
            hf_mulliken: hf_mul,
            mp2_mulliken: mp2_mul,
            hf_nbo: hf_nbo,
            mp2_nbo: mp2_nbo,
            hf_esp: hf_esp,
            mp2_esp: mp2_esp,
        })
    }
}
