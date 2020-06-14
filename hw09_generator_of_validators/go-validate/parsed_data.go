package main

import (
	"fmt"
	"reflect"
)

type Node struct {
	IsSet bool
	Value interface{}
}

type FieldType struct {
	IsSlice bool
	Name    string
}

type FieldParams struct {
	Regexp Node
	Len    Node
	Enum   Node
	Min    Node
	Max    Node
	Type   FieldType
	Name   string
}

type ParsedData struct {
	PackageName string
	Structures  map[string][]FieldParams
}

func (p *FieldParams) GetEnumList() []string {
	if !p.Enum.IsSet {
		return nil
	}
	v := p.Enum.Value
	if reflect.TypeOf(v).Kind() == reflect.Slice {
		s := reflect.ValueOf(v)
		result := make([]string, 0, s.Len())
		for i := 0; i < s.Len(); i++ {
			str, ok := s.Index(i).Interface().(string)
			if !ok {
				return nil
			}
			if p.Type.Name == "string" {
				str = fmt.Sprintf("\"%s\"", str)
			}
			result = append(result, str)
		}
		return result
	}
	return nil
}
