package main

import (
	"github.com/philopon/fmoe/cpf2svl/cpf"
	"github.com/philopon/fmoe/cpf2svl/svlwriter"
)

func writeCpf(writer *svlwriter.SVLWriter, cpf *cpf.Cpf) error {
	if err := writer.WriteInt([]int{int(cpf.NumAtoms)}); err != nil {
		return err
	}
	if err := writer.WriteInt([]int{int(cpf.NumFrags)}); err != nil {
		return err
	}

	if err := writer.WriteInt(cpf.AtomIndices); err != nil {
		return err
	}
	if err := writer.WriteToken(cpf.AtomElements); err != nil {
		return err
	}
	if err := writer.WriteToken(cpf.AtomTypes); err != nil {
		return err
	}
	if err := writer.WriteToken(cpf.AtomResNames); err != nil {
		return err
	}
	if err := writer.WriteInt(cpf.AtomResIndices); err != nil {
		return err
	}
	if err := writer.WriteInt(cpf.AtomFragIndices); err != nil {
		return err
	}
	if err := writer.WriteFloat(cpf.AtomX); err != nil {
		return err
	}
	if err := writer.WriteFloat(cpf.AtomY); err != nil {
		return err
	}
	if err := writer.WriteFloat(cpf.AtomZ); err != nil {
		return err
	}
	if err := writer.WriteFloat(cpf.AtomHFMulliken); err != nil {
		return err
	}
	if err := writer.WriteFloat(cpf.AtomMP2Mulliken); err != nil {
		return err
	}
	if err := writer.WriteFloat(cpf.AtomHFNBO); err != nil {
		return err
	}
	if err := writer.WriteFloat(cpf.AtomMP2NBO); err != nil {
		return err
	}
	if err := writer.WriteFloat(cpf.AtomHFRESP); err != nil {
		return err
	}
	if err := writer.WriteFloat(cpf.AtomMP2RESP); err != nil {
		return err
	}
	if err := writer.WriteToken(cpf.AtomChainID); err != nil {
		return err
	}
	if err := writer.WriteToken(cpf.AtomInsCode); err != nil {
		return err
	}

	if err := writer.WriteInt(cpf.FragBondNumbers); err != nil {
		return err
	}
	if err := writer.WriteInt(cpf.FragBondSelfs); err != nil {
		return err
	}
	if err := writer.WriteInt(cpf.FragBondOthers); err != nil {
		return err
	}

	if err := writer.WriteFloat(cpf.DimerDistances); err != nil {
		return err
	}
	if err := writer.WriteFloat(cpf.DimerES); err != nil {
		return err
	}
	if err := writer.WriteFloat(cpf.DimerDI); err != nil {
		return err
	}
	if err := writer.WriteFloat(cpf.DimerEX); err != nil {
		return err
	}
	if err := writer.WriteFloat(cpf.DimerCT); err != nil {
		return err
	}

	return nil
}
