package http

import "gopkg.in/yaml.v3"

type LineField[T any] struct {
	File  string
	Line  int
	Value T
}

var currentParseFile string

// 解析yaml
func (f *LineField[T]) UnmarshalYAML(node *yaml.Node) (err error) {
	var v T
	if err := node.Decode(&v); err != nil {
		return err
	}

	f.File = currentParseFile
	f.Line = node.Line
	f.Value = v

	return
}
