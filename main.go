package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func getInFile(arg1 string) (io.ReadCloser, error) {
	if arg1 == "" {
		return os.Stdin, nil
	}
	f, err := os.Open(arg1)
	return f, err
}

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{""}
	}
	in, err := getInFile(args[0])
	if err != nil {
		return err
	}

	d, err := ioutil.ReadAll(in)
	in.Close()
	if err != nil {
		return err
	}

	d, err = edit(d)
	if err != nil {
		return err
	}

	if args[0] == "" {
		_, err = os.Stdout.Write(d)
	} else {
		err = ioutil.WriteFile(args[0], d, os.ModePerm) // file already exists
	}
	return err
}
