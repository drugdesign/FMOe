package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aymerick/raymond"
	flags "github.com/jessevdk/go-flags"
)

type options struct {
	InputPath  string `short:"i" long:"input" description:"input template file" env:"TEMPLATE_PATH"`
	OutputPath string `short:"o" long:"output" description:"output file" env:"OUTPUT_PATH"`
}

func mainProcess() error {
	var opts options
	if _, err := flags.ParseArgs(&opts, os.Args); err != nil {
		return err
	}
	var input *os.File
	if opts.InputPath == "" {
		input = os.Stdin
	} else {
		if f, err := os.Open(opts.InputPath); err == nil {
			input = f
		} else {
			return err
		}
	}

	var output *os.File
	if opts.OutputPath == "" {
		output = os.Stdout
	} else {
		if f, err := os.Create(opts.OutputPath); err == nil {
			output = f
		} else {
			return err
		}
	}

	b, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}
	tpl := string(b)

	rawEnvs := os.Environ()
	envs := make(map[string]string, len(rawEnvs))

	for _, rawEnv := range rawEnvs {
		if i := strings.IndexRune(rawEnv, '='); i >= 0 {
			envs[rawEnv[:i]] = rawEnv[i+1:]
		}
	}

	result, err := raymond.Render(tpl, envs)
	if err != nil {
		return err
	}
	fmt.Fprintf(output, result)

	return nil
}

func main() {
	if err := mainProcess(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
