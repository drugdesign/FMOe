# FMOE

![fmotool](./img/fmotool.gif)

## installation

1. clone or download the repository

    ```
    git clone https://github.com/drugdesign/FMOE.git $SOME_DIRECTORY
    ```

    or

    download source code from github release(https://github.com/drugdesign/FMOe/releases/latest)


2. write the code running loader.svl to `$HOME/moefiles/start.svl`

    for example:

    ```$HOME/moefiles/start.svl
    local function main []
        run '$FMOE_INSTALL_PATH/loader.svl';
    endfunction
    ```
