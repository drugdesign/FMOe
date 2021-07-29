package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

type options struct {
	AutoFrag string `short:"i" long:"input" description:"abinitmp autofrag log" env:"AUTOFRAG_PATH"`
	SvlBin   string `short:"o" long:"output" description:"auto frag svl binary data" env:"SVLBIN_PATH"`
}

type SVLWriter struct {
	writer *bufio.Writer
}

func NewSVLWriter(writer io.Writer) SVLWriter {
	return SVLWriter{writer: bufio.NewWriter(writer)}
}

func (w *SVLWriter) Flush() error {
	return w.writer.Flush()
}

func (w *SVLWriter) WriteSize(v uint32) error {
	return binary.Write(w.writer, binary.BigEndian, v)
}

func (w *SVLWriter) WriteInt(vals []int) error {
	if err := w.writer.WriteByte(2); err != nil {
		return err
	}
	if err := w.WriteSize(uint32(len(vals))); err != nil {
		return err
	}
	for _, val := range vals {
		if err := binary.Write(w.writer, binary.BigEndian, int32(val)); err != nil {
			return err
		}
	}
	return nil
}

type AutoFrag struct {
	BDA []int
	BAA []int
}

type AutoFragParser struct {
	scanner *bufio.Scanner
	result  AutoFrag
}

func (p *AutoFragParser) scan() (string, error) {
	if p == nil {
		return "", errors.New("error")
	}
	r := p.scanner.Scan()
	if r {
		return p.scanner.Text(), nil
	}
	return "", errors.New("error")
}

func (p *AutoFragParser) parse() (*AutoFrag, error) {
	for {
		line, err := p.scan()
		if err != nil {
			return &p.result, errors.New("error")
		}
		if strings.Contains(line, "Frag.   Bonded Atom  Proj.") {
			break
		}
	}
	for {
		line, err := p.scan()
		if err != nil {
			return &p.result, errors.New("error")
		}
		if len(line) < 13 {
			break
		}
		if _, err := strconv.Atoi(strings.TrimSpace(line[0:13])); err != nil {
			break
		} else {
			if v, err := strconv.Atoi(strings.TrimSpace(line[14:21])); err == nil {
				p.result.BDA = append(p.result.BDA, v)
			} else {
				return &p.result, errors.New("error")
			}
			if v, err := strconv.Atoi(strings.TrimSpace(line[22:27])); err == nil {
				p.result.BAA = append(p.result.BAA, v)
			} else {
				return &p.result, errors.New("error")
			}
		}
	}
	return &p.result, nil
}

func ParseAutoFrag(reader io.Reader) (*AutoFrag, error) {
	parser := AutoFragParser{scanner: bufio.NewScanner(reader)}
	return parser.parse()
}

const (
	ok                int = 0
	optionParseFailed     = 1
	ioError               = 2
	parseError            = 3
)

func mainProcess() (int, error) {
	var opts options
	if _, err := flags.ParseArgs(&opts, os.Args); err != nil {
		return optionParseFailed, err
	}
	var file *os.File
	if opts.AutoFrag == "" {
		return ioError, errors.New("input isn't set")
	} else {
		file, _ = os.Open(opts.AutoFrag)
	}
	defer file.Close()
	var autofrag, _ = ParseAutoFrag(file)
	var output *os.File
	if opts.SvlBin == "" {
		output = os.Stdout
	} else {
		var err error
		output, err = os.Create(opts.SvlBin)
		if err != nil {
			return ioError, err
		}
	}
	defer output.Close()
	var writer = NewSVLWriter(output)
	writer.WriteInt(autofrag.BDA)
	writer.WriteInt(autofrag.BAA)
	writer.Flush()
	return ok, nil
}

func main() {
	code, err := mainProcess()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
	os.Exit(code)
}
