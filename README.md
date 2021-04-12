# FMOE

![fmotool](./img/fmotool.gif)

## Installation

1. Clone or download this repository

    ```
    git clone https://github.com/drugdesign/FMOe.git $SOME_DIRECTORY
    ```

    or

    download source code from github release(https://github.com/drugdesign/FMOe/releases/latest)


2. Write the code running loader.svl to `$HOME/moefiles/start.svl`

    for example:

    ```$HOME/moefiles/start.svl
    local function main []
        run '$FMOE_INSTALL_PATH/loader.svl';
    endfunction
    ```

## Features

### Custom Templates for ajf/sh files

FMOe generates ajf/sh files from templates in templates directory.

#### Variables
- {{{BASENAME}}}
- {{{TOTAL_CHARGE}}}
- {{{NUM_FRAGS}}}
- {{{BASIS_SET}}}
- {{{ABINITMP_FRAGMENT}}}


