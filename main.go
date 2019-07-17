package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	var (
		rootDir string
	)

	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("%+v", errors.Wrap(err, ""))
		os.Exit(1)
	}

	flag.StringVar(&rootDir, "rootDir", rootDir, "Root directory when reading OpenAPI file")
	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	rootPath, err := filepath.Abs(rootDir)
	if err != nil {
		fmt.Printf("%+v", errors.Wrap(err, ""))
		os.Exit(1)
	}

	filePath, err := filepath.Abs(flag.Args()[0])
	if err != nil {
		fmt.Printf("%+v", errors.Wrap(err, ""))
		os.Exit(1)
	}

	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		fmt.Printf("%+v", errors.Wrap(err, ""))
		os.Exit(1)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
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

	loc, err := url.Parse(filePath)
	if err != nil {
		fmt.Printf("%+v", errors.Wrap(err, ""))
		os.Exit(1)
	}

	loader := openapi3.NewSwaggerLoader()
	loader.IsExternalRefsAllowed = true

	openapi, err := loader.LoadSwaggerFromDataWithPath(data, loc)
	if err != nil {
		fmt.Printf("%+v", errors.Wrap(err, ""))
		os.Exit(1)
	}

	if err := openapi.Validate(context.Background()); err != nil {
		fmt.Printf("%+v", errors.Wrap(err, ""))
		os.Exit(1)
	}

	showInfo(openapi)

	for path, pathItem := range openapi.Paths {
		for k, v := range pathItem.ExtensionProps.Extensions {
			if k == "$ref" {
				str, err := json.MarshalToString(v)
				if err != nil {
					fmt.Printf("%+v", errors.Wrap(err, ""))
					os.Exit(1)
				}

				log.Println(path, k, str)
			}
		}
	}
}

func showInfo(openapi *openapi3.Swagger) {
	log.Printf("%-15s %s", "OpenAPI:", openapi.OpenAPI)
	log.Printf("%-15s %s", "Title:", openapi.Info.Title)
	if openapi.Info.Description != "" {
		log.Printf("%-15s %s", "Description:", openapi.Info.Description)
	}
	if openapi.Info.TermsOfService != "" {
		log.Printf("%-15s %s", "TermsOfService:", openapi.Info.TermsOfService)
	}
	if openapi.Info.Contact != nil {
		if openapi.Info.Contact.Name != "" {
			log.Printf("%-15s %s", "Contact->Name:", openapi.Info.Contact.Name)
		}
		if openapi.Info.Contact.Email != "" {
			log.Printf("%-15s %s", "Contact->Email:", openapi.Info.Contact.Email)
		}
		if openapi.Info.Contact.URL != "" {
			log.Printf("%-15s %s", "Contact->URL:", openapi.Info.Contact.URL)
		}
	}
	if openapi.Info.License != nil {
		if openapi.Info.License.Name != "" {
			log.Printf("%-15s %s", "License->Name:", openapi.Info.License.Name)
		}
		if openapi.Info.License.URL != "" {
			log.Printf("%-15s %s", "License->URL:", openapi.Info.License.URL)
		}
	}
	log.Printf("%-15s %s", "Version:", openapi.Info.Version)
}
