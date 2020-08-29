package generator

import (
	"fmt"
	"regexp"
	"strings"
)

type TemplateField struct {
	//{{FieldName}}
	Name string
	// {{structName}}.{{FieldName}}
	FullName string
	Type     string
}

type TemplateFieldValidator struct {
	Field TemplateField
	// int, max, min
	Name string
	// arg1, arg2, arg3
	Args string
}

type TemplateValidateStruct struct {
	Type       string
	Alias      string
	Validators []TemplateFieldValidator
}
type TemplateData struct {
	Package string
	Imports map[string]string

	ValidationTargets []TemplateValidateStruct
}

// генерирует вывод
type structureValidatorGenerator struct {
}

func NewStructureValidatorGenerator() *structureValidatorGenerator {
	return &structureValidatorGenerator{}
}

func (g *structureValidatorGenerator) Generate(file *validationFile) (*TemplateData, error) {

	// collect validators
	var tmplData = &TemplateData{
		Package:           file.PkgName,
		Imports: map[string]string{
			"fmt" : "fmt",
		},
		ValidationTargets: []TemplateValidateStruct{},
	}

	for _, item := range file.Structures {
		tmplData.ValidationTargets = append(tmplData.ValidationTargets, TemplateValidateStruct{
			Type:       item.Name,
			Alias:      item.Alias,
			Validators: g.extractValidators(item),
		})
	}

	return tmplData, nil
}

func (g *structureValidatorGenerator) extractValidators(validationItem *validationStruct) []TemplateFieldValidator {
	var collectValidators []TemplateFieldValidator

	for _, field := range validationItem.fields {
		// parse tag
		for _, tag := range field.Tags {
			var args string

			// escape strings
			if tag.Name == "in" && (field.Type == "string" ||  field.Type == "[]string"){
				args = fmt.Sprintf("[]string{\"%s\"}", strings.Join(tag.Args, "\", \""))
			}else if tag.Name == "in" && (field.Type == "int" ||  field.Type == "[]int"){
				args = fmt.Sprintf("[]int{%s}", strings.Join(tag.Args, ","))
			}else if tag.Name == "regexp" {
				args = fmt.Sprintf("\"%s\"", regexp.QuoteMeta(tag.Args[0]))
			} else{
				args = strings.Join(tag.Args, ",")
			}

			collectValidators = append(collectValidators, TemplateFieldValidator{
				Field: TemplateField{
					Name:     field.Name,
					FullName: fmt.Sprintf("%s.%s", validationItem.Alias, field.Name),
					Type:     field.Type,
				},
				Name: tag.Name,
				Args: args,
			})
		}

	}

	return collectValidators
}
