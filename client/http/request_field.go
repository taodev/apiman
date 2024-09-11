package http

import "gopkg.in/yaml.v3"

type LineField[T any] struct {
	Line  int
	Value T
}

// 解析yaml
func (f *LineField[T]) UnmarshalYAML(node *yaml.Node) (err error) {
	var v T
	if err := node.Decode(&v); err != nil {
		return err
	}

	f.Line = node.Line
	f.Value = v

	return
}
