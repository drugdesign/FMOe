package cpf

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	errors "github.com/pkg/errors"
)

// Version is abinit-mp check point file (CPF) Version type
type Version int

const (
	// Ver7_2 is Ver.7.2 CPF
	Ver7_2 Version = 72
	// Ver4_201MIZUHO is Ver.4.201 (MIZUHO)
	Ver4_201MIZUHO = 4201
	//
)

// Cpf is CPF file
type Cpf struct {
	Version  Version
	NumAtoms int
	NumFrags int

	AtomIndices     []int
	AtomElements    []string
	AtomTypes       []string
	AtomResNames    []string
	AtomResIndices  []int
	AtomFragIndices []int
	AtomX           []float64
	AtomY           []float64
	AtomZ           []float64
	AtomHFMulliken  []float64
	AtomMP2Mulliken []float64
	AtomHFNBO       []float64
	AtomMP2NBO      []float64
	AtomHFRESP      []float64
	AtomMP2RESP     []float64
	AtomChainID     []string
	AtomInsCode     []string

	FragBondNumbers []int
	FragBondSelfs   []int
	FragBondOthers  []int

	DimerDistances []float64
	DimerES        []float64
	DimerDI        []float64
	DimerEX        []float64
	DimerCT        []float64
}

type cpfParser struct {
	scanner *bufio.Scanner
	result  Cpf
}

// ParseCpf parse cpf file
func ParseCpf(reader io.Reader) (*Cpf, error) {
	parser := cpfParser{scanner: bufio.NewScanner(reader)}
	return parser.parse()
}

// ErrNilPointerReciever indicate nil pointer reciever error
var ErrNilPointerReciever = errors.New("nil pointer reciever")

func (cpf *cpfParser) scan() (string, error) {
	if cpf == nil {
		return "", ErrNilPointerReciever
	}
	r := cpf.scanner.Scan()
	if r {
		return cpf.scanner.Text(), nil
	}

	return "", cpf.scanner.Err()
}

// UnknownCPFVersion error
type UnknownCPFVersion struct{ Version string }

func (err *UnknownCPFVersion) Error() string {
	return fmt.Sprintf("unknown CPF version: %s", err.Version)
}

func (cpf *cpfParser) parseVersion() error {
	line, err := cpf.scan()
	if err != nil {
		return err
	}

	if strings.HasPrefix(line, "CPF Ver.7.2") {
		cpf.result.Version = Ver7_2
	} else if strings.HasPrefix(line, "CPF Ver.4.201 (MIZUHO)") {
		cpf.result.Version = Ver4_201MIZUHO
	} else {
		return &UnknownCPFVersion{Version: line}
	}

	return nil
}

func (cpf *cpfParser) parseNumAtomsAndNumFrags() error {
	line, err := cpf.scan()
	if err != nil {
		return err
	}

	if na, err := intField(line, 0, 5); err == nil {
		cpf.result.NumAtoms = na
	} else {
		return err
	}

	if nf, err := intField(line, 5, 10); err == nil {
		cpf.result.NumFrags = nf
	} else {
		return err
	}
	return nil
}

func (cpf *cpfParser) parseAtoms() error {
	cpf.result.AtomIndices = make([]int, cpf.result.NumAtoms)
	cpf.result.AtomElements = make([]string, cpf.result.NumAtoms)
	cpf.result.AtomTypes = make([]string, cpf.result.NumAtoms)
	cpf.result.AtomResNames = make([]string, cpf.result.NumAtoms)
	cpf.result.AtomResIndices = make([]int, cpf.result.NumAtoms)
	cpf.result.AtomFragIndices = make([]int, cpf.result.NumAtoms)
	cpf.result.AtomX = make([]float64, cpf.result.NumAtoms)
	cpf.result.AtomY = make([]float64, cpf.result.NumAtoms)
	cpf.result.AtomZ = make([]float64, cpf.result.NumAtoms)
	cpf.result.AtomHFMulliken = make([]float64, cpf.result.NumAtoms)
	cpf.result.AtomMP2Mulliken = make([]float64, cpf.result.NumAtoms)
	cpf.result.AtomHFNBO = make([]float64, cpf.result.NumAtoms)
	cpf.result.AtomMP2NBO = make([]float64, cpf.result.NumAtoms)
	cpf.result.AtomHFRESP = make([]float64, cpf.result.NumAtoms)
	cpf.result.AtomMP2RESP = make([]float64, cpf.result.NumAtoms)
	cpf.result.AtomChainID = make([]string, cpf.result.NumAtoms)
	cpf.result.AtomInsCode = make([]string, cpf.result.NumAtoms)

	for i := 0; i < cpf.result.NumAtoms; i++ {
		line, err := cpf.scan()
		if err != nil {
			return err
		}
		if v, err := intField(line, 0, 5); err == nil {
			cpf.result.AtomIndices[i] = v
		} else {
			return err
		}
		if v, err := slice(line, 6, 8); err == nil {
			cpf.result.AtomElements[i] = v
		} else {
			return err
		}
		if v, err := slice(line, 9, 13); err == nil {
			cpf.result.AtomTypes[i] = v
		} else {
			return err
		}
		if v, err := slice(line, 14, 17); err == nil {
			cpf.result.AtomResNames[i] = v
		} else {
			return err
		}
		if v, err := intField(line, 18, 22); err == nil {
			cpf.result.AtomResIndices[i] = v
		} else {
			return err
		}
		if v, err := intField(line, 23, 27); err == nil {
			cpf.result.AtomFragIndices[i] = v
		} else {
			return err
		}
		if v, err := floatField(line, 28, 40); err == nil {
			cpf.result.AtomX[i] = v
		} else {
			return err
		}
		if v, err := floatField(line, 40, 52); err == nil {
			cpf.result.AtomY[i] = v
		} else {
			return err
		}
		if v, err := floatField(line, 52, 64); err == nil {
			cpf.result.AtomZ[i] = v
		} else {
			return err
		}
		if v, err := floatField(line, 64, 76); err == nil {
			cpf.result.AtomHFMulliken[i] = v
		} else {
			return err
		}
		if v, err := floatField(line, 76, 88); err == nil {
			cpf.result.AtomMP2Mulliken[i] = v
		} else {
			return err
		}
		if v, err := floatField(line, 88, 100); err == nil {
			cpf.result.AtomHFNBO[i] = v
		} else {
			return err
		}
		if v, err := floatField(line, 100, 112); err == nil {
			cpf.result.AtomMP2NBO[i] = v
		} else {
			return err
		}
		if v, err := floatField(line, 112, 124); err == nil {
			cpf.result.AtomHFRESP[i] = v
		} else {
			return err
		}
		if v, err := floatField(line, 124, 136); err == nil {
			cpf.result.AtomMP2RESP[i] = v
		} else {
			return err
		}
		if v, err := slice(line, 137, 138); IsStringOutOfRange(err) {
			cpf.result.AtomChainID[i] = " "
		} else if err == nil {
			cpf.result.AtomChainID[i] = v
		} else {
			return err
		}
		if v, err := slice(line, 139, 140); IsStringOutOfRange(err) {
			cpf.result.AtomInsCode[i] = " "
		} else if err == nil {
			cpf.result.AtomInsCode[i] = v
		} else {
			return err
		}
	}
	return nil
}

func (cpf *cpfParser) skip(lines int) error {
	for i := 0; i < lines; i++ {
		_, err := cpf.scan()
		if err != nil {
			return err
		}
	}
	return nil
}

func (cpf *cpfParser) skipFragElectrons() error {
	lines := cpf.result.NumFrags / 16
	if cpf.result.NumFrags%16 > 0 {
		lines++
	}
	return cpf.skip(lines)
}

func (cpf *cpfParser) parseFragBondNumbers() error {
	cpf.result.FragBondNumbers = make([]int, 0, cpf.result.NumFrags)

	lines := cpf.result.NumFrags / 16
	rests := cpf.result.NumFrags % 16
	for l := 0; l < lines; l++ {
		line, err := cpf.scan()
		if err != nil {
			return err
		}

		for j := 0; j < 16; j++ {
			if v, err := intField(line, 0, 5); err == nil {
				cpf.result.FragBondNumbers = append(cpf.result.FragBondNumbers, v)
				line = line[5:]
			} else {
				return err
			}
		}
	}

	if rests == 0 {
		return nil
	}

	line, err := cpf.scan()
	if err != nil {
		return err
	}
	for j := 0; j < rests; j++ {
		if v, err := intField(line, 0, 5); err == nil {
			cpf.result.FragBondNumbers = append(cpf.result.FragBondNumbers, v)
			line = line[5:]
		} else {
			return err
		}
	}
	return nil
}

func (cpf *cpfParser) getFragBonds() int {
	i := 0
	for _, num := range cpf.result.FragBondNumbers {
		i += num
	}
	return i
}

func (cpf *cpfParser) parseFragBonds(fragBonds int) error {
	cpf.result.FragBondSelfs = make([]int, fragBonds)
	cpf.result.FragBondOthers = make([]int, fragBonds)

	for i := 0; i < fragBonds; i++ {
		line, err := cpf.scan()
		if err != nil {
			return err
		}

		if v, err := intField(line, 0, 5); err == nil {
			cpf.result.FragBondOthers[i] = v
		} else {
			return err
		}
		if v, err := intField(line, 5, 10); err == nil {
			cpf.result.FragBondSelfs[i] = v
		} else {
			return err
		}
	}
	return nil
}

// MissingFields error
type MissingFields struct {
	Index  int
	String string
}

func (err *MissingFields) Error() string {
	return fmt.Sprintf("missing %d-th field in %s", err.Index, err.String)
}

func (cpf *cpfParser) parseDimerDistances(numDimers int) error {
	cpf.result.DimerDistances = make([]float64, numDimers)

	for i := 0; i < numDimers; i++ {
		line, err := cpf.scan()
		if err != nil {
			return err
		}
		fs := strings.Fields(line)
		if len(fs) < 2 {
			return &MissingFields{Index: 2, String: line}
		}

		if v, err := strconv.ParseFloat(strings.TrimSpace(fs[2]), 64); err == nil {
			cpf.result.DimerDistances[i] = v
		} else {
			return err
		}
	}
	return nil
}

func (cpf *cpfParser) parseDimers(numDimers int) error {
	cpf.result.DimerES = make([]float64, numDimers)
	cpf.result.DimerDI = make([]float64, numDimers)
	cpf.result.DimerEX = make([]float64, numDimers)
	cpf.result.DimerCT = make([]float64, numDimers)

	var exStart, ctStart int
	fieldWidth := 24
	if cpf.result.Version == Ver4_201MIZUHO {
		exStart = 288
		ctStart = 312
	} else if cpf.result.Version == Ver7_2 {
		exStart = 336
		ctStart = 360
	}
	exEnd := exStart + fieldWidth
	ctEnd := ctStart + fieldWidth

	for i := 0; i < numDimers; i++ {
		line, err := cpf.scan()
		if err != nil {
			return err
		}
		if v, err := floatField(line, 48, 72); err == nil {
			cpf.result.DimerES[i] = v
		} else {
			return err
		}

		if v, err := floatField(line, 72, 96); err == nil {
			cpf.result.DimerDI[i] = v
		} else {
			return err
		}

		if v, err := floatField(line, exStart, exEnd); err == nil {
			cpf.result.DimerEX[i] = v
		} else {
			return err
		}
		if v, err := floatField(line, ctStart, ctEnd); err == nil {
			cpf.result.DimerCT[i] = v
		} else {
			return err
		}
	}
	return nil
}

func (cpf *cpfParser) parse() (*Cpf, error) {
	if err := cpf.parseVersion(); err != nil {
		return nil, errors.Wrap(err, "parse version")
	}
	if err := cpf.parseNumAtomsAndNumFrags(); err != nil {
		return nil, errors.Wrap(err, "parse number of atoms and numbber of fragments")
	}

	if err := cpf.parseAtoms(); err != nil {
		return nil, errors.Wrap(err, "parse atoms")
	}

	if err := cpf.skipFragElectrons(); err != nil {
		return nil, errors.Wrap(err, "skip fragment electrons")
	}

	if err := cpf.parseFragBondNumbers(); err != nil {
		return nil, errors.Wrap(err, "parse fragment bond numbers")
	}

	fragBonds := cpf.getFragBonds()
	if err := cpf.parseFragBonds(fragBonds); err != nil {
		return nil, errors.Wrap(err, "parse fragment bonds")
	}

	numDimers := (cpf.result.NumFrags * (cpf.result.NumFrags - 1)) / 2

	if err := cpf.parseDimerDistances(numDimers); err != nil {
		return nil, errors.Wrap(err, "parse dimer distances")
	}

	// (dipole moment + monomers) + information
	if err := cpf.skip(2*cpf.result.NumFrags + 7); err != nil {
		return nil, errors.Wrap(err, "skip informations")
	}

	if err := cpf.parseDimers(numDimers); err != nil {
		return nil, errors.Wrap(err, "parse dimers")
	}

	return &cpf.result, nil
}
