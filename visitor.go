package pkgst

import (
	"go/ast"
	"go/token"
	"log"
	"reflect"
)

type WalkVisitor struct {
	Packages map[string]*Package
}

func (v *WalkVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node == nil {
		return nil
	}

	switch node := node.(type) {
	case *ast.File:
		if v.Packages == nil {
			v.Packages = map[string]*Package{}
		}

		pkgName := node.Name.Name
		pkg, ok := v.Packages[pkgName]
		if !ok {
			pkg = &Package{Name: pkgName}
			v.Packages[pkgName] = pkg
		}

		return &fileVisitor{
			Package:  pkg,
			Packages: v.Packages,
		}
	}
	return v
}

type fileVisitor struct {
	Package  *Package
	Packages map[string]*Package
}

func (v *fileVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node == nil {
		return nil
	}

	switch node := node.(type) {
	case *ast.GenDecl:
		switch node.Tok {
		case token.IMPORT:
			for _, sp := range node.Specs {
				imp := importSpec(sp.(*ast.ImportSpec))
				if imp != nil {
					v.Package.Impoprts = append(v.Package.Impoprts, imp)
				}
			}
		case token.TYPE:
			for _, sp := range node.Specs {
				decl := typeSpec(sp.(*ast.TypeSpec))
				if decl != nil {
					v.Package.Types = append(v.Package.Types, decl)
				}
			}
		case token.CONST:
			for _, sp := range node.Specs {
				decls := valueSpec(sp.(*ast.ValueSpec))
				if len(decls) > 0 {
					v.Package.Consts = append(v.Package.Consts, decls...)
				}
			}
		case token.VAR:
			for _, sp := range node.Specs {
				valueSpec(sp.(*ast.ValueSpec))
				decls := valueSpec(sp.(*ast.ValueSpec))
				if len(decls) > 0 {
					v.Package.Vars = append(v.Package.Vars, decls...)
				}
			}
		}
		return v
	case *ast.FuncDecl:
		f := fun(node.Type)
		if node.Recv != nil {
			f.Receiver = ty(node.Recv.List[0].Type)
		}
		v.Package.Funcs = append(v.Package.Funcs, &Decl{
			Name: node.Name.Name,
			Type: &Type{F: f},
		})
		return v
	}
	return nil
}

func importSpec(spec *ast.ImportSpec) *Import {
	name := ""
	if spec.Name != nil {
		name = spec.Name.Name
	}
	return &Import{Name: name, Path: spec.Path.Value}
}

func typeSpec(spec *ast.TypeSpec) *Decl {
	t := ty(spec.Type)
	if t == nil {
		return nil
	}
	return &Decl{Name: spec.Name.Name, Type: t}
}

func valueSpec(spec *ast.ValueSpec) []*Decl {
	var decls []*Decl
	for i, name := range spec.Names {
		if spec.Type != nil {
			decls = append(decls, &Decl{
				Name: name.Name,
				Type: ty(spec.Type),
			})
		} else {
			decls = append(decls, &Decl{
				Name: name.Name,
				Type: litTy(spec.Values[i]),
			})
		}
	}
	return decls
}

func ty(expr ast.Expr) *Type {
	switch expr := expr.(type) {
	case *ast.StarExpr:
		ty := ty(expr.X)
		ty.Star = ty.Star + 1
		return ty
	case *ast.Ident:
		name := expr.Name
		if isPrimitive(name) {
			return &Type{P: Primitive(expr.Name)}
		}
		return &Type{R: &Ref{Name: expr.Name}}
	case *ast.ArrayType:
		return &Type{A: &Array{Type: ty(expr.Elt)}}
	case *ast.MapType:
		return &Type{M: &Map{
			KeyType:   ty(expr.Key),
			ValueType: ty(expr.Value),
		}}
	case *ast.FuncType:
		return &Type{F: fun(expr)}
	case *ast.StructType:
		var f []*Decl
		if expr.Fields.NumFields() > 0 {
			f = fields(expr.Fields.List)
		}
		return &Type{S: &Struct{Fields: f}}
	case *ast.InterfaceType:
		var f []*Decl
		if expr.Methods.NumFields() > 0 {
			f = fields(expr.Methods.List)
		}
		return &Type{I: &Interface{Methods: f}}
	case *ast.Ellipsis:
		return &Type{E: &Ellipse{Type: ty(expr.Elt)}}
	case *ast.SelectorExpr:
		return &Type{R: &Ref{
			Package: pkg(expr.X),
			Name:    expr.Sel.Name,
		}}
	default:
		log.Println("**** unknown type:", reflect.TypeOf(expr).String(), expr)
		return nil
	}
}

func pkg(expr ast.Expr) string {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr.Name
	}
	return ""
}

func fun(fun *ast.FuncType) *Function {
	var args []*Decl
	var results []*Decl

	if fun.Params.NumFields() > 0 {
		args = fields(fun.Params.List)
	}
	if fun.Results.NumFields() > 0 {
		results = fields(fun.Results.List)
	}
	return &Function{
		Args:    args,
		Results: results,
	}
}

func fields(fields []*ast.Field) []*Decl {
	var decls []*Decl
	for _, field := range fields {
		if len(field.Names) == 0 {
			decls = append(decls, &Decl{Type: ty(field.Type)})
			continue
		}
		for _, name := range field.Names {
			decls = append(decls, &Decl{
				Name: name.Name,
				Type: ty(field.Type),
			})
		}
	}
	return decls
}

func litTy(expr ast.Expr) *Type {
	switch expr := expr.(type) {
	case *ast.CompositeLit:
		return ty(expr.Type)
	case *ast.BasicLit:
		return &Type{P: basicLit(expr.Kind)}
	case *ast.FuncLit:
		return ty(expr.Type)
	}
	return nil
}

func basicLit(kind token.Token) Primitive {
	switch kind {
	case token.INT: // 12345
		return Int
	case token.FLOAT: // 123.45
		return Float32
	case token.IMAG: // 123.45i
		return Complex
	case token.CHAR: // 'a'
		return Rune
	case token.STRING: // "abc"
		return String
	}
	return Undefined
}
