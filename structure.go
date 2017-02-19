package pkgst

import "unicode"
import "log"

type _map map[string]interface{}

type FileScope struct {
	Impoprts []*Import
}

type Package struct {
	Name   string  `json:"name,omitempty"`
	Types  []*Decl `json:"types,omitempty"`
	Funcs  []*Decl `json:"functions,omitempty"`
	Consts []*Decl `json:"consts,omitempty"`
	Vars   []*Decl `json:"vars,omitempty"`
}

func (p *Package) Type(name string) *Decl {
	for _, d := range p.Types {
		if d.Name == name {
			return d
		}
	}
	return nil
}

func (p *Package) Func(name string) *Decl {
	for _, d := range p.Funcs {
		if d.Name == name && d.Type.F.Receiver == nil {
			return d
		}
	}
	return nil
}

func (p *Package) Const(name string) *Decl {
	for _, d := range p.Consts {
		if d.Name == name {
			return d
		}
	}
	return nil
}

func (p *Package) Var(name string) *Decl {
	for _, d := range p.Vars {
		if d.Name == name {
			return d
		}
	}
	return nil
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
	if d.Name == "" {
		return false
	}
	r := rune(d.Name[0])
	return unicode.IsUpper(r)
}

type Type struct {
	S *Struct    `json:"s,omitempty"`
	I *Interface `json:"i,omitempty"`
	F *Function  `json:"f,omitempty"`
	A *Array     `json:"a,omitempty"`
	M *Map       `json:"m,omitempty"`
	E *Ellipse   `json:"e,omitempty"`
	P *Primitive `json:"p,omitempty"`
	R *Ref       `json:"r,omitempty"`

	Ptr int `json:"ptr,omitempty"`

	pkg     *Package
	imports []*Import
	methods []*Decl
}

func (t *Type) Methods() []*Decl {
	if t.methods != nil {
		return t.methods
	}

	t.methods = []*Decl{}
	for _, f := range t.pkg.Funcs {
		rcv := f.Type.F.Receiver
		if rcv == nil {
			continue
		}
		if t.Real() == rcv.Real() {
			t.methods = append(t.methods, f)
		}
	}
	return t.methods
}

func (t *Type) Real() *Type {
	if t.R != nil {
		if t.R.Package != "" {
			log.Println("not support external package", t.R.Package)
		}
		return t.pkg.Type(t.R.Name).Type
	}
	return t
}

func newType(t interface{}) *Type {
	switch t := t.(type) {
	case *Struct:
		return &Type{S: t}
	case *Interface:
		return &Type{I: t}
	case *Function:
		return &Type{F: t}
	case *Array:
		return &Type{A: t}
	case *Map:
		return &Type{M: t}
	case *Ellipse:
		return &Type{E: t}
	case *Ref:
		return &Type{R: t}
	case *Primitive:
		return &Type{P: t}
	default:
		return nil
	}
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

type Primitive struct {
	Base PrimName
}

type PrimName string

const (
	Undefined  PrimName = ""
	String              = "string"
	Bool                = "bool"
	Int                 = "int"
	Int8                = "int8"
	Int16               = "int16"
	Int32               = "int32"
	Int64               = "int64"
	Uint                = "uint"
	Uint8               = "uint8"
	Uint16              = "uint16"
	Uint32              = "uint32"
	Uint64              = "uint64"
	Uintptr             = "uintptr"
	Byte                = "byte"
	Rune                = "rune"
	Float32             = "float32"
	Float64             = "float64"
	Complex64           = "complex64"
	Complex128          = "complex128"
)

func primName(s string) PrimName {
	p := PrimName(s)
	if p == String ||
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
		p == Complex128 {
		return p
	}
	return Undefined
}
