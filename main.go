package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Printf("%+v\n", errors.New("There is only one argument"))
		os.Exit(1)
	}

	filePath, err := filepath.Abs(filepath.Clean(flag.Args()[0]))
	if err != nil {
		fmt.Printf("%+v\n", errors.Wrap(err, ""))
		os.Exit(1)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("%+v\n", errors.Wrap(err, ""))
		os.Exit(1)
	}

	workDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("%+v\n", errors.Wrap(err, ""))
		os.Exit(1)
	}

	workDir, err = filepath.Abs(workDir)
	if err != nil {
		fmt.Printf("%+v\n", errors.Wrap(err, ""))
		os.Exit(1)
	}
	defer os.Chdir(workDir)

	os.Chdir(filepath.Dir(filePath))

	loader := openapi3.NewSwaggerLoader()
	loader.IsExternalRefsAllowed = true

	openapi, err := loader.LoadSwaggerFromFile(filePath)
	if err != nil {
		fmt.Printf("%+v\n", errors.Wrap(err, ""))
		os.Exit(1)
	}

	if err := openapi.Validate(context.Background()); err != nil {
		fmt.Printf("%+v\n", errors.Wrap(err, ""))
		os.Exit(1)
	}
}
