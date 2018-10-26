extern crate byteorder;
#[cfg(test)]
#[macro_use]
extern crate lazy_static;

mod checkpoint;

use checkpoint::CheckPointFile;
use std::env;
use std::error::Error;
use std::fs;
use std::io;

fn main() -> Result<(), Box<Error>> {
    let cpf_path = env::var("CPF_PATH")?;
    let svl_path = env::var("SVL_PATH")?;
    let cpf = {
        let src = fs::File::open(cpf_path)?;
        CheckPointFile::parse(src)?
    };
    {
        let mut dst = io::BufWriter::new(fs::File::create(svl_path)?);
        cpf.serialize(&mut dst)?;
    }
    Ok(())
}
