package checker

import "fmt"

type LexoType int

const (
	TypeInt LexoType = iota
	TypeFloat
	TypeString
	TypeBool
	TypeIntArray
	TypeFloatArray
	TypeStringArray
	TypeBoolArray
	TypeVoid
)

var typeNames = map[LexoType]string{
	TypeInt: "int", TypeFloat: "float", TypeString: "string", TypeBool: "bool",
	TypeIntArray: "int[]", TypeFloatArray: "float[]", TypeStringArray: "string[]", TypeBoolArray: "bool[]",
	TypeVoid: "void",
}

func (t LexoType) String() string {
	if name, ok := typeNames[t]; ok {
		return name
	}
	return fmt.Sprintf("unknown(%d)", int(t))
}

func TypeFromAnnotation(ann string) (LexoType, bool) {
	switch ann {
	case "int":
		return TypeInt, true
	case "float":
		return TypeFloat, true
	case "string":
		return TypeString, true
	case "bool":
		return TypeBool, true
	case "int[]":
		return TypeIntArray, true
	case "float[]":
		return TypeFloatArray, true
	case "string[]":
		return TypeStringArray, true
	case "bool[]":
		return TypeBoolArray, true
	case "void":
		return TypeVoid, true
	default:
		return 0, false
	}
}

func ElementType(arrayType LexoType) (LexoType, bool) {
	switch arrayType {
	case TypeIntArray:
		return TypeInt, true
	case TypeFloatArray:
		return TypeFloat, true
	case TypeStringArray:
		return TypeString, true
	case TypeBoolArray:
		return TypeBool, true
	default:
		return 0, false
	}
}

func ArrayTypeOf(elemType LexoType) (LexoType, bool) {
	switch elemType {
	case TypeInt:
		return TypeIntArray, true
	case TypeFloat:
		return TypeFloatArray, true
	case TypeString:
		return TypeStringArray, true
	case TypeBool:
		return TypeBoolArray, true
	default:
		return 0, false
	}
}

func IsArrayType(t LexoType) bool {
	return t == TypeIntArray || t == TypeFloatArray || t == TypeStringArray || t == TypeBoolArray
}
