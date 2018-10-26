use super::types;
use byteorder::{BigEndian, WriteBytesExt};
use std::io;

pub trait WriteSvl: io::Write {
    fn write_size(&mut self, size: usize) -> io::Result<()> {
        self.write_u32::<BigEndian>(size as u32)
    }

    fn write_int(&mut self, v: &[i32]) -> io::Result<()> {
        self.write(&[2])?;
        self.write_size(v.len())?;
        for a in v.iter() {
            self.write_i32::<BigEndian>(*a)?;
        }
        Ok(())
    }
    fn write_float(&mut self, v: &[f64]) -> io::Result<()> {
        self.write(&[3])?;
        self.write_size(v.len())?;
        for a in v.iter() {
            self.write_f64::<BigEndian>(*a)?;
        }
        Ok(())
    }

    fn write_token(&mut self, v: &[&str]) -> io::Result<()> {
        self.write(&[4])?;
        self.write_size(v.len())?;
        for a in v.iter() {
            self.write_size(a.len())?;
            self.write(&a.bytes().collect::<Vec<_>>())?;
        }
        Ok(())
    }
}

impl<T: io::Write> WriteSvl for T {}

impl types::Atom {
    fn serialize<W: WriteSvl>(&self, w: &mut W) -> io::Result<()> {
        w.write(&[5])?;
        w.write_size(15)?;
        w.write_int(&[self.index])?;
        w.write_token(&[&self.element])?;
        w.write_token(&[&self.atom_type])?;
        w.write_token(&[&self.res_name])?;
        w.write_int(&[self.res_index])?;
        w.write_int(&[self.frag_index])?;
        w.write_float(&[self.x])?;
        w.write_float(&[self.y])?;
        w.write_float(&[self.z])?;
        w.write_float(&[self.hf_mulliken])?;
        w.write_float(&[self.mp2_mulliken])?;
        w.write_float(&[self.hf_nbo])?;
        w.write_float(&[self.mp2_nbo])?;
        w.write_float(&[self.hf_esp])?;
        w.write_float(&[self.mp2_esp])
    }
}

impl types::Vec3<f64> {
    fn serialize<W: WriteSvl>(&self, w: &mut W) -> io::Result<()> {
        w.write_float(&[self.x, self.y, self.z])
    }
}

impl types::SplittedBond {
    fn serialize<W: WriteSvl>(&self, w: &mut W) -> io::Result<()> {
        w.write_int(&[self.bda, self.baa])
    }
}

impl types::CheckPointFile {
    pub fn serialize<W: WriteSvl>(&self, w: &mut W) -> io::Result<()> {
        w.write_token(&[&self.version])?;
        w.write_token(&[&self.comments])?;
        w.write_int(&[self.num_atoms])?;
        w.write_int(&[self.num_frags])?;

        w.write(&[5])?;
        w.write_size(self.atoms.len())?;
        for atom in self.atoms.iter() {
            atom.serialize(w)?;
        }

        w.write(&[5])?;
        w.write_size(self.monomers.atom_indices.len())?;
        for indices in self.monomers.atom_indices.iter() {
            w.write_int(&indices.iter().map(|v| v + 1).collect::<Vec<_>>())?;
        }
        w.write_int(&self.monomers.num_electrons)?;
        w.write_int(&self.monomers.num_attached)?;

        w.write(&[5])?;
        w.write_size(self.monomers.hf_dipoles.len())?;
        for dipole in self.monomers.hf_dipoles.iter() {
            dipole.serialize(w)?;
        }
        w.write(&[5])?;
        w.write_size(self.monomers.mp2_dipoles.len())?;
        for dipole in self.monomers.mp2_dipoles.iter() {
            dipole.serialize(w)?;
        }

        w.write_float(&self.monomers.nuclear_repulsion)?;
        w.write_float(&self.monomers.electronic)?;
        w.write_float(&self.monomers.mp2)?;
        w.write_float(&self.monomers.mp3)?;
        w.write_int(&self.monomers.ao)?;
        w.write_int(&self.monomers.mo)?;

        w.write(&[5])?;
        w.write_size(self.splitted_bonds.len())?;
        for bond in self.splitted_bonds.iter() {
            bond.serialize(w)?;
        }

        w.write_float(&self.dimers.distances.data)?;
        w.write_float(&self.dimers.nuclear_repulsion.data)?;
        w.write_float(&self.dimers.electronic.data)?;
        w.write_float(&self.dimers.electrostatic.data)?;
        w.write_float(&self.dimers.mp2.data)?;
        w.write_float(&self.dimers.scs_mp2.data)?;
        w.write_float(&self.dimers.mp3.data)?;
        w.write_float(&self.dimers.scs_mp3.data)?;
        w.write_float(&self.dimers.hf_bsse.data)?;
        w.write_float(&self.dimers.mp2_bsse.data)?;
        w.write_float(&self.dimers.scs_mp2_bsse.data)?;
        w.write_float(&self.dimers.mp3_bsse.data)?;
        w.write_float(&self.dimers.scs_mp3_bsse.data)?;
        w.write_float(&self.dimers.exchange_repulsion.data)?;
        w.write_float(&self.dimers.charge_transfer.data)?;
        w.write_float(&self.dimers.qij.data)?;

        w.write_token(&[&self.level])?;
        w.write_token(&[&self.electronic_state])?;
        w.write_token(&[&self.basisset])?;

        w.write_float(&[self.approx_params.ao_population])?;
        w.write_float(&[self.approx_params.point_charge])?;
        w.write_float(&[self.approx_params.dimer_es])?;

        w.write_float(&[self.system_energy.nuclear_repulsion])?;
        w.write_float(&[self.system_energy.electronic])?;
        w.write_float(&[self.system_energy.total])?;

        Ok(())
    }
}

#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn test_write_float() {
        let mut d: Vec<u8> = Vec::with_capacity(128);
        d.write_float(&[1f64]).unwrap();
        assert_eq!(d, &[3, 0, 0, 0, 1, 63, 240, 0, 0, 0, 0, 0, 0]);

        d.clear();
        d.write_float(&[1f64, 2f64]).unwrap();
        assert_eq!(
            d,
            &[3, 0, 0, 0, 2, 63, 240, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0]
        );
    }
    #[test]
    fn test_write_token() {
        let mut d: Vec<u8> = Vec::with_capacity(128);
        d.write_token(&["foo", "bar"]).unwrap();
        assert_eq!(
            d,
            &[4, 0, 0, 0, 2, 0, 0, 0, 3, 102, 111, 111, 0, 0, 0, 3, 98, 97, 114]
        );
    }

}
