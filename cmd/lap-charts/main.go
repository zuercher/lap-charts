package main

import (
	"flag"
	"fmt"
	"os"

	"zuercher.us/lapcharts"
)

func errorf(f string, args ...any) {
	fmt.Fprintf(os.Stdout, f, args...)
	os.Exit(1)
}

func main() {
	options := &lapcharts.Options{}

	options.ConfigureFlags()
	flag.Parse()

	if err := options.Validate(); err != nil {
		errorf("ERROR: invalid options: %s\n", err.Error())
	}

	if err := lapcharts.Generate(options); err != nil {
		errorf("ERROR: %s\n", err.Error())
	}
}
