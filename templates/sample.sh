#!/bin/sh

export OMP_NUM_THREADS=8

#------- ABINIT-MP option -------#
ABINIT_DIR=ABINIT-MP
BINARY_NAME=abinitmp_smp
MKINP_NAME=mkinp.py

#------- Input option -------#
FILE_NAME={{{BASENAME}}}
AJF_NAME=$FILE_NAME.ajf
OUT_NAME=$FILE_NAME.log

#------- Program execution -------#
$ABINIT_DIR/$MKINP_NAME < $AJF_NAME > $FILE_NAME.new.ajf
mpiexec $ABINIT_DIR/${BINARY_NAME} < $FILE_NAME.ajf > $OUT_NAME

