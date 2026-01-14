package render

type MetaType int

const (
	MetaTypeSize MetaType = iota
	MetaTypeDate
	MetaTypeString
)

type Meta struct {
	Type  MetaType
	Value any
}
