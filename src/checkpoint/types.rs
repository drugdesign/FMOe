#[derive(Debug)]
pub struct Vec3<T> {
    pub x: T,
    pub y: T,
    pub z: T,
}

impl<T> Vec3<T> {
    pub fn new(x: T, y: T, z: T) -> Self {
        Vec3 { x: x, y: y, z: z }
    }
}

#[derive(Debug)]
pub struct Monomers {
    pub atom_indices: Vec<Vec<i32>>,
    pub num_electrons: Vec<i32>,
    pub num_attached: Vec<i32>,
    pub hf_dipoles: Vec<Vec3<f64>>,
    pub mp2_dipoles: Vec<Vec3<f64>>,

    pub nuclear_repulsion: Vec<f64>,
    pub electronic: Vec<f64>,
    pub mp2: Vec<f64>,
    pub mp3: Vec<f64>,
    pub ao: Vec<i32>,
    pub mo: Vec<i32>,
}

#[derive(Debug)]
pub struct DimerData<T> {
    pub data: Vec<T>,
}

impl<T> DimerData<T> {
    pub fn new(v: Vec<T>) -> Self {
        DimerData { data: v }
    }
}

#[derive(Debug)]
pub struct ApproxParams {
    pub ao_population: f64,
    pub point_charge: f64,
    pub dimer_es: f64,
}

#[derive(Debug)]
pub struct SystemEnergy {
    pub nuclear_repulsion: f64,
    pub electronic: f64,
    pub total: f64,
}

#[derive(Debug)]
pub struct Dimers {
    pub distances: DimerData<f64>,
    pub nuclear_repulsion: DimerData<f64>,
    pub electronic: DimerData<f64>,
    pub electrostatic: DimerData<f64>,
    pub mp2: DimerData<f64>,
    pub scs_mp2: DimerData<f64>,
    pub mp3: DimerData<f64>,
    pub scs_mp3: DimerData<f64>,
    pub hf_bsse: DimerData<f64>,
    pub mp2_bsse: DimerData<f64>,
    pub scs_mp2_bsse: DimerData<f64>,
    pub mp3_bsse: DimerData<f64>,
    pub scs_mp3_bsse: DimerData<f64>,
    pub exchange_repulsion: DimerData<f64>,
    pub charge_transfer: DimerData<f64>,
    pub qij: DimerData<f64>,
}

#[derive(Debug)]
pub struct SplittedBond {
    pub bda: i32,
    pub baa: i32,
}

#[derive(Debug, Clone)]
pub struct Atom {
    pub index: i32,
    pub element: String,
    pub atom_type: String,
    pub res_name: String,
    pub res_index: i32,
    pub frag_index: i32,
    pub x: f64,
    pub y: f64,
    pub z: f64,
    pub hf_mulliken: f64,
    pub mp2_mulliken: f64,
    pub hf_nbo: f64,
    pub mp2_nbo: f64,
    pub hf_esp: f64,
    pub mp2_esp: f64,
}

#[derive(Debug)]
pub struct CheckPointFile {
    pub version: String,
    pub comments: String,
    pub num_atoms: i32,
    pub num_frags: i32,
    pub atoms: Vec<Atom>,
    pub monomers: Monomers,
    pub splitted_bonds: Vec<SplittedBond>,
    pub dimers: Dimers,
    pub level: String,
    pub electronic_state: String,
    pub basisset: String,
    pub approx_params: ApproxParams,
    pub system_energy: SystemEnergy,
}
