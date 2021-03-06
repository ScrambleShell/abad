package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/NeowayLabs/abad"
	"github.com/NeowayLabs/abad/cmd/abad/cli"
)

func repl() error {

	cli, err := cli.NewCli(os.Stdin, os.Stdout)
	if err != nil {
		return err
	}

	cli.Repl()
	return nil
}

func eval(codepath string) error {
	code, err := ioutil.ReadFile(codepath)
	if err != nil {
		return err
	}
	abadjs, err := abad.NewAbad()
	if err != nil {
		return err
	}
	_, err = abadjs.EvalFile(filepath.Base(codepath), string(code))
	return err
}

func main() {
	var execute string
	var help bool

	flag.BoolVar(&help, "help", false, "prints usage")
	flag.StringVar(&execute, "e", "", "execute code")
	flag.Parse()

	if help {
		fmt.Println("Abad: the bad JS interpreter")
		flag.PrintDefaults()
		return
	}

	if execute != "" {
		abadjs, err := abad.NewAbad()
		abortonerr(err)
		_, err = abadjs.Eval(execute)
		abortonerr(err)
		return
	}

	if len(flag.Args()) == 0 {
		abortonerr(repl())
		return
	}

	filepath := flag.Args()[0]
	abortonerr(eval(filepath))
}

func abortonerr(err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
