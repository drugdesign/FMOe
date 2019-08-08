package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	flags "github.com/jessevdk/go-flags"
	"github.com/philopon/fmoe/cpf2svl/cpf"
	"github.com/philopon/fmoe/cpf2svl/svlwriter"
)

type options struct {
	CpfPath string `short:"i" long:"input" description:"input cpf file" env:"CPF_PATH"`
	SvlPath string `short:"o" long:"output" description:"output moe binary file" env:"SVL_PATH"`
	JSON    bool   `short:"j" long:"json" description:"json output"`
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

	if !opts.JSON && opts.SvlPath == "" {
		return optionParseFailed, fmt.Errorf("the required flag `-o, --output' was not specified")
	}

	var file *os.File
	if opts.CpfPath == "" {
		file = os.Stdin
	} else {
		var err error
		file, err = os.Open(opts.CpfPath)
		if err != nil {
			return ioError, err
		}
	}
	defer file.Close()

	var input io.Reader
	if filepath.Ext(opts.CpfPath) == ".gz" {
		var err error
		input, err = gzip.NewReader(file)
		if err != nil {
			return ioError, err
		}
	} else {
		input = file
	}

	cpf, err := cpf.ParseCpf(input)
	if err != nil {
		return parseError, err
	}

	var output *os.File
	if opts.SvlPath == "" {
		output = os.Stdout
	} else {
		var err error
		output, err = os.Create(opts.SvlPath)
		if err != nil {
			return ioError, err
		}
	}
	defer output.Close()

	if opts.JSON {
		enc := json.NewEncoder(output)
		enc.Encode(cpf)

	} else {
		w := svlwriter.NewSVLWriter(output)
		if err := writeCpf(&w, cpf); err != nil {
			return ioError, err
		}
		if err := w.Flush(); err != nil {
			return ioError, err
		}
	}

	return ok, nil
}

func main() {
	code, err := mainProcess()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
	os.Exit(code)
}
