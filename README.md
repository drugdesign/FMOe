# FMOE

![fmotool](./img/fmotool.gif)

## Features

FMOe is a GUI-based utility library that supports the preprocessing and visualization of molecular structures for Fragment Molecular Orbital (FMO) calculations.  
It runs within the Molecular Operating Environment (MOE) and is designed to generate input files for ABINIT-MP and assist in the interpretation of its results.

### Key Features

- **Support for fragmentation and fragment merging**  
  Automatically divides biomolecules into fragments based on amino acid or nucleotide units.  
  Also supports nonstandard residues, cyclic peptides, and coordination complexes.  
  Users can manually assign fragments or merge functionally related units, such as metal ions and surrounding residues.

- **Template-based generation of ajf/sh input files**  
  Automatically creates ABINIT-MP input files with user-defined templates, inserting values such as total charge, number of fragments, and basis set.

- **Visualization and analysis of FMO results**  
  Graphical tools are available for inspecting fragment interaction energies (IFIE) and PIEDA-derived energy components.

- **Easy integration with MOE**  
  FMOe modules can be automatically loaded at MOE startup by adding a loader script to `start.svl`.

## Installation

1. Clone or download this repository

    ```
    git clone https://github.com/drugdesign/FMOe.git $SOME_DIRECTORY
    ```

    or

    download source code from github release(https://github.com/drugdesign/FMOe/releases/latest)


2. Write the code running loader.svl to `$HOME/moefiles/start.svl`

    for example:

    ```svl
    local function main []
        run '$FMOE_INSTALL_PATH/loader.svl';
    endfunction
    ```

    ```svl
    # if you want to set the templates path
    local function main []
        run ['$FMOE_INSTALL_PATH/loader.svl', [templates: '/path/to/your/templates/directory']];
    endfunction
    ```
    
### Custom Templates for ajf/sh files

FMOe generates ajf/sh files from templates in templates directory.

#### Variables
- {{{BASENAME}}}
- {{{TOTAL_CHARGE}}}
- {{{NUM_FRAGS}}}
- {{{BASIS_SET}}}
- {{{ABINITMP_FRAGMENT}}}

## Citation

If you find FMOe useful for your work, please cite the following article:

> Hirotomo Moriwaki, Yusuke Kawashima, Chiduru Watanabe*, Kikuko Kamisaka, Yoshio Okiyama, Kaori Fukuzawa, Teruki Honma.  
> *FMOe: Preprocessing and Visualizing Package of the Fragment Molecular Orbital Method for Molecular Operating Environment and Its Applications in Covalent Ligand and Metalloprotein Analyses*.  
> *Journal of Chemical Information and Modeling*, **64**(18), Application Note, September 5, 2024.  
> [https://doi.org/10.1021/acs.jcim.4c01169](https://doi.org/10.1021/acs.jcim.4c01169)
