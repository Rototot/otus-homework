package main

import (
	"flag"
	"github.com/Rototot/otus-homework/hw09_generator_of_validators/go-validate/generator"
	"log"
	"os"
)

var source string
var out string

func main() {
	flag.StringVar(&source, "source", "models/models.go", "")
	flag.StringVar(&out, "out", "models/model_generated.go", "")
	flag.Parse()

	// parse
	parseResult, err := generator.NewFileParser(source).Parse()

	// generate and render

	// generate
	output, err := os.Create(out)
	if err != nil {
		log.Fatalln(err)
	}
	defer output.Close()

	renderData, err := generator.NewStructureValidatorGenerator().Generate(parseResult)
	if err != nil {
		log.Fatalln(err)
	}

	err = generator.Render(output, renderData)
	if err != nil {
		log.Fatalln(err)
	}
}
