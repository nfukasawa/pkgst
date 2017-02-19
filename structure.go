package pkgst

import "unicode"

type Package struct {
	Name     string    `json:"name,omitempty"`
	Impoprts []*Import `json:"imports,omitempty"`
	Types    []*Decl   `json:"types,omitempty"`
	Funcs    []*Decl   `json:"functions,omitempty"`
	Consts   []*Decl   `json:"consts,omitempty"`
	Vars     []*Decl   `json:"vars,omitempty"`
}

type Import struct {
	Name string
	Path string
}

type Decl struct {
	Name string `json:"name,omitempty"`
	Type *Type  `json:"type,omitempty"`
}

func (d *Decl) IsPublic() bool {
	return isPublicName(d.Name)
}

type Type struct {
	Star int        `json:"star,omitempty"`
	S    *Struct    `json:"s,omitempty"`
	I    *Interface `json:"i,omitempty"`
	F    *Function  `json:"f,omitempty"`
	A    *Array     `json:"a,omitempty"`
	M    *Map       `json:"m,omitempty"`
	E    *Ellipse   `json:"e,omitempty"`
	R    *Ref       `json:"r,omitempty"`
	P    Primitive  `json:"p,omitempty"`
}

type Struct struct {
	Fields []*Decl `json:"fields,omitempty"`
}

type Interface struct {
	Methods []*Decl `json:"methods,omitempty"`
}

type Function struct {
	Receiver *Type   `json:"receiver,omitempty"`
	Args     []*Decl `json:"args,omitempty"`
	Results  []*Decl `json:"results,omitempty"`
}

type Array struct {
	Type *Type `json:"type,omitempty"`
}

type Map struct {
	KeyType   *Type `json:"keyType,omitempty"`
	ValueType *Type `json:"valueType,omitempty"`
}

type Ellipse struct {
	Type *Type `json:"type,omitempty"`
}

type Ref struct {
	Package string `json:"package,omitempty"`
	Name    string `json:"name,omitempty"`
}

type Primitive string

const (
	Undefined  Primitive = ""
	String               = "string"
	Bool                 = "bool"
	Int                  = "int"
	Int8                 = "int8"
	Int16                = "int16"
	Int32                = "int32"
	Int64                = "int64"
	Uint                 = "uint"
	Uint8                = "uint8"
	Uint16               = "uint16"
	Uint32               = "uint32"
	Uint64               = "uint64"
	Uintptr              = "uintptr"
	Byte                 = "byte"
	Rune                 = "rune"
	Float32              = "float32"
	Float64              = "float64"
	Complex64            = "complex64"
	Complex128           = "complex128"
)

func isPrimitive(s string) bool {
	p := Primitive(s)
	return (p == String ||
		p == Bool ||
		p == Int ||
		p == Int8 ||
		p == Int16 ||
		p == Int32 ||
		p == Int64 ||
		p == Uint ||
		p == Uint8 ||
		p == Uint16 ||
		p == Uint32 ||
		p == Uint64 ||
		p == Uintptr ||
		p == Byte ||
		p == Rune ||
		p == Float32 ||
		p == Float64 ||
		p == Complex64 ||
		p == Complex128)
}

func isPublicName(s string) bool {
	if s == "" {
		return false
	}
	r := rune(s[0])
	return unicode.IsUpper(r)
}
