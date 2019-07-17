package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func main() {
	var (
		rootDir string
	)

	workDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("%+v", errors.Wrap(err, ""))
		os.Exit(1)
	}

	flag.StringVar(&rootDir, "rootDir", workDir, "Root directory when reading OpenAPI file")
	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	filePath := flag.Args()[0]

	if _, err := os.Stat(filePath); err != nil {
		fmt.Printf("%+v", errors.Wrap(err, ""))
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("%+v", errors.Wrap(err, ""))
		os.Exit(1)
	}

	var version struct {
		Swagger string `json:"swagger,omitempty" yaml:"swagger,omitempty"`
		OpenAPI string `json:"openapi,omitempty" yaml:"openapi,omitempty"`
	}

	switch filepath.Ext(filePath) {
	case ".json":
		if err := json.Unmarshal(data, &version); err != nil {
			fmt.Printf("%+v", errors.Wrap(err, ""))
			os.Exit(1)
		}
	case ".yml", ".yaml":
		if err := yaml.Unmarshal(data, &version); err != nil {
			fmt.Printf("%+v", errors.Wrap(err, ""))
			os.Exit(1)
		}
	default:
		fmt.Printf("%+v", errors.New("Unsupported file extension"))
		os.Exit(1)
	}

	switch {
	case version.Swagger == "" && version.OpenAPI == "":
		fmt.Printf("%+v", errors.New("Invalid data"))
		os.Exit(1)
	case version.Swagger != "" && version.OpenAPI != "":
		fmt.Printf("%+v", errors.New("Invalid data"))
		os.Exit(1)
	case version.Swagger != "":
		// TODO: Support Swagger/OpenAPI 2.0
		fmt.Printf("%+v", errors.New("Unsupported Swagger/OpenAPI 2.0"))
		os.Exit(1)
	}

	log.Println(version)
}
