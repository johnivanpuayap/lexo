package interpreter

import (
	"fmt"
	"strings"

	"github.com/johnivanpuayap/lexo/pkg/checker"
)

type Value interface {
	Type() checker.LexoType
	String() string
}

type IntVal int64
type FloatVal float64
type StringVal string
type BoolVal bool
type VoidVal struct{}

type ArrayVal struct {
	Elements []Value
	ElemType checker.LexoType
}

func (v IntVal) Type() checker.LexoType  { return checker.TypeInt }
func (v IntVal) String() string           { return fmt.Sprintf("%d", int64(v)) }

func (v FloatVal) Type() checker.LexoType { return checker.TypeFloat }
func (v FloatVal) String() string          { return fmt.Sprintf("%g", float64(v)) }

func (v StringVal) Type() checker.LexoType { return checker.TypeString }
func (v StringVal) String() string          { return string(v) }

func (v BoolVal) Type() checker.LexoType { return checker.TypeBool }
func (v BoolVal) String() string {
	if bool(v) {
		return "true"
	}
	return "false"
}

func (v VoidVal) Type() checker.LexoType { return checker.TypeVoid }
func (v VoidVal) String() string          { return "void" }

func (v *ArrayVal) Type() checker.LexoType {
	at, _ := checker.ArrayTypeOf(v.ElemType)
	return at
}

func (v *ArrayVal) String() string {
	parts := make([]string, len(v.Elements))
	for i, e := range v.Elements {
		if _, ok := e.(StringVal); ok {
			parts[i] = fmt.Sprintf("%q", e.String())
		} else {
			parts[i] = e.String()
		}
	}
	return "[" + strings.Join(parts, ", ") + "]"
}
