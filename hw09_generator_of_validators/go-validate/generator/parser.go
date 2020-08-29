package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

const (
	validateTagName = "validate"
	tagSeparator    = "|"
	tagArgsSeparator    = ":"
)

type validationTag struct {
	Name string
	Args []string
}

type validationField struct {
	Name string
	Type string
	Tags []validationTag
}
type validationStruct struct {
	Name  string
	Alias string

	fields []*validationField
}

type validationFile struct {
	PkgName    string
	Structures []*validationStruct
}

type FileParser struct {
	filePath string
}

func NewFileParser(filePath string) *FileParser {
	return &FileParser{filePath: filePath}
}

func (p *FileParser) Parse() (*validationFile, error) {
	fset := token.NewFileSet()
	fileAst, err := parser.ParseFile(fset, p.filePath, nil, parser.ParseComments)

	if err != nil {
		return nil, err
	}

	var fileTarget = &validationFile{
		// grab package Name
		PkgName: fileAst.Name.Name,
	}

	// grab fields
	ast.Inspect(fileAst, func(node ast.Node) bool {
		// get spec for grab struct Name and type
		spec, ok := node.(*ast.TypeSpec)
		if !ok {
			return true
		}

		if spec.Type == nil {
			return true
		}

		// grab fields
		nodeType, ok := spec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		// skip empty struct
		if nodeType.Fields == nil || len(nodeType.Fields.List) == 0 {
			return true
		}

		var validationFields []*validationField
		for _, field := range nodeType.Fields.List {
			typeName, ok := p.extractFieldTypeName(field.Type)
			if !ok {
				continue
			}

			if field.Tag != nil{
				validationFields = append(validationFields, &validationField{
					Name: field.Names[0].Name,
					Type: typeName,
					Tags: p.extractTags(field.Tag.Value),
				})
			}
		}

		fileTarget.Structures = append(fileTarget.Structures, &validationStruct{
			Name:   spec.Name.Name,
			Alias:  strings.ToLower(spec.Name.Name[:1]),
			fields: validationFields,
		})

		return false
	})

	return fileTarget, nil
}

func (p *FileParser) extractTags(rawTags string) []validationTag {
	strTags, ok := reflect.StructTag(strings.Trim(rawTags, "`")).Lookup(validateTagName)
	if !ok {
		return nil
	}

	var tags []validationTag

	// validator:args -> struct{Name, Args}
	for _, tag := range strings.Split(strTags, tagSeparator) {
		tagParts := strings.Split(tag, tagArgsSeparator)
		if  len(tagParts) > 0{
			tags = append(tags, validationTag{
				Name: tagParts[0],
				Args: tagParts[1:],
			})
		}
	}

	return tags
}

func (p *FileParser) extractFieldTypeName(rawType ast.Expr) (string, bool) {
	var typeName string

	if fieldType, ok := rawType.(*ast.Ident); ok {
		typeName = fieldType.Name
	}

	// array --> []Type
	if fieldType, ok := rawType.(*ast.ArrayType); ok {
		name, ok := fieldType.Elt.(*ast.Ident)
		if !ok {
			return "", false
		}
		typeName = "[]" + name.Name
	}

	return typeName, true
}

