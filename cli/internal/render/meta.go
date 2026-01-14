package render

import (
	"fmt"
	"time"

	"github.com/IllumiKnowLabs/labstore/cli/internal/format"
)

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

func NewMetaString(str string) Meta {
	return Meta{Type: MetaTypeString, Value: str}
}

func NewMetaDate(date time.Time) Meta {
	return Meta{Type: MetaTypeDate, Value: date}
}

func NewMetaSize(size int64) Meta {
	return Meta{Type: MetaTypeSize, Value: size}
}

func (m Meta) Render() string {
	switch m.Type {
	case MetaTypeString:
		return m.Value.(string)
	case MetaTypeDate:
		return format.Date(m.Value.(time.Time))
	case MetaTypeSize:
		return format.Size(m.Value.(int64))
	default:
		return fmt.Sprint(m.Value)
	}
}
