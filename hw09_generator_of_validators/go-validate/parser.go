package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
	"strings"
)

var (
	supportedTypes = map[string]struct{}{
		"string": {},
		"int":    {},
	}
)

func parse(fn string) (ParsedData, error) {
	parsedData := ParsedData{}
	set := token.NewFileSet()
	node, err := parser.ParseFile(set, fn, nil, parser.ParseComments)
	if err != nil {
		return parsedData, err
	}
	// fill package names
	parsedData.PackageName = node.Name.Name

	// find custom user types
	customKindTypes := parseCustomKindType(node.Decls)
	// structs parse begin
	structParams := make(map[string][]FieldParams)
	for _, f := range node.Decls {
		g, ok := f.(*ast.GenDecl)
		if !ok {
			continue
		}
	FIELDS_LOOP:
		for _, spec := range g.Specs {
			currType, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue FIELDS_LOOP
			}
			currStruct, ok := currType.Type.(*ast.StructType)
			if !ok {
				continue FIELDS_LOOP
			}
			// parse structs validation params
			structField := getValidationParams(currStruct, customKindTypes)

			if len(structField) > 0 {
				structParams[currType.Name.Name] = structField
			}
		}
	}
	parsedData.Structures = structParams
	return parsedData, nil
}

func getValidationParams(currStruct *ast.StructType, customKindTypes map[string]string) (fieldParams []FieldParams) {
	for _, field := range currStruct.Fields.List {
		var fieldParamsItem FieldParams
		// check tag
		var validationTag string
		if field.Tag != nil {
			tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
			validationTag = strings.TrimSpace(tag.Get(parseTagName))
		}
		if validationTag == "" {
			continue
		}
		// field name
		fieldParamsItem.Name = field.Names[0].Name

		// parse validation rule types
		tags := strings.Split(validationTag, "|")
		for _, tagItem := range tags {
			tagChunks := strings.Split(tagItem, ":")
			if len(tagChunks) != 2 {
				// unexpected tag format
				continue
			}
			switch tagChunks[0] {
			case "regexp":
				fieldParamsItem.Regexp = Node{true, tagChunks[1]}
			case "len":
				lenValue, _ := strconv.Atoi(tagChunks[1])
				fieldParamsItem.Len = Node{true, lenValue}
			case "min":
				minValue, _ := strconv.Atoi(tagChunks[1])
				fieldParamsItem.Min = Node{true, minValue}
			case "max":
				maxValue, _ := strconv.Atoi(tagChunks[1])
				fieldParamsItem.Max = Node{true, maxValue}
			case "in":
				fieldParamsItem.Enum = Node{true, strings.Split(tagChunks[1], ",")}
			}
		}

		// check field type
		fieldType := getFieldType(field.Type, customKindTypes)
		if fieldType.Name == "" {
			continue
		}
		fieldParamsItem.Type = fieldType
		fieldParams = append(fieldParams, fieldParamsItem)
	}

	return
}

func getFieldType(fieldTypeExpr ast.Expr, customKindTypes map[string]string) (fieldType FieldType) {
	var fieldTypeName string
	if fieldTypeIdent, ok := fieldTypeExpr.(*ast.Ident); ok {
		fieldTypeName = getIdentTypeName(fieldTypeIdent, customKindTypes)
		fieldType.IsSlice = false
	} else if fieldTypeArr, ok := fieldTypeExpr.(*ast.ArrayType); ok {
		if fieldTypeIdent, ok := fieldTypeArr.Elt.(*ast.Ident); ok {
			fieldTypeName = getIdentTypeName(fieldTypeIdent, customKindTypes)
			fieldType.IsSlice = true
		}
	}
	if _, ok := supportedTypes[fieldTypeName]; ok {
		fieldType.Name = fieldTypeName
	}
	return
}

func getIdentTypeName(fieldTypeIdent *ast.Ident, customKindTypes map[string]string) string {
	fieldName := fieldTypeIdent.Name
	if kind, ok := customKindTypes[fieldName]; ok {
		fieldName = kind
	}
	return fieldName
}

func parseCustomKindType(decls []ast.Decl) map[string]string {
	result := make(map[string]string)
	for _, f := range decls {
		g, ok := f.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range g.Specs {
			currType, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			typeIdent, ok := currType.Type.(*ast.Ident)
			if !ok {
				continue
			}
			result[currType.Name.Name] = typeIdent.Name
		}
	}
	return result
}
