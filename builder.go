package pkgst

import (
	"go/ast"
	"go/token"
	"log"
	"reflect"
)

func Build(pkgs map[string]*ast.Package) map[string]*Package {
	res := map[string]*Package{}
	for _, pkg := range pkgs {
		res[pkg.Name] = &Package{Name: pkg.Name}
	}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			visitor := &fileVisitor{
				Package: res[file.Name.Name],
				Scope:   new(FileScope),
			}
			ast.Walk(visitor, file)
		}
	}
	return res
}

type fileVisitor struct {
	Package *Package
	Scope   *FileScope
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
				imp := v.importSpec(sp.(*ast.ImportSpec))
				if imp != nil {
					v.Scope.Impoprts = append(v.Scope.Impoprts, imp)
				}
			}
		case token.TYPE:
			for _, sp := range node.Specs {
				decl := v.typeSpec(sp.(*ast.TypeSpec))
				if decl != nil {
					v.Package.Types = append(v.Package.Types, decl)
				}
			}
		case token.CONST:
			for _, sp := range node.Specs {
				decls := v.valueSpec(sp.(*ast.ValueSpec))
				if len(decls) > 0 {
					v.Package.Consts = append(v.Package.Consts, decls...)
				}
			}
		case token.VAR:
			for _, sp := range node.Specs {
				decls := v.valueSpec(sp.(*ast.ValueSpec))
				if len(decls) > 0 {
					v.Package.Vars = append(v.Package.Vars, decls...)
				}
			}
		}
	case *ast.FuncDecl:
		f := v.fun(node.Type)
		if node.Recv != nil {
			f.Receiver = v.ty(node.Recv.List[0].Type)
		}
		v.Package.Funcs = append(v.Package.Funcs, &Decl{
			Name: node.Name.Name,
			Type: v.newType(f),
		})
	}
	return v
}

func (v *fileVisitor) importSpec(spec *ast.ImportSpec) *Import {
	name := ""
	if spec.Name != nil {
		name = spec.Name.Name
	}
	return &Import{Name: name, Path: spec.Path.Value}
}

func (v *fileVisitor) typeSpec(spec *ast.TypeSpec) *Decl {
	t := v.ty(spec.Type)
	if t == nil {
		return nil
	}
	return &Decl{Name: spec.Name.Name, Type: t}
}

func (v *fileVisitor) valueSpec(spec *ast.ValueSpec) []*Decl {
	var decls []*Decl
	for i, name := range spec.Names {
		if spec.Type != nil {
			decls = append(decls, &Decl{
				Name: name.Name,
				Type: v.ty(spec.Type),
			})
		} else {
			decls = append(decls, &Decl{
				Name: name.Name,
				Type: v.litTy(spec.Values[i]),
			})
		}
	}
	return decls
}

func (v *fileVisitor) ty(expr ast.Expr) *Type {
	switch expr := expr.(type) {
	case *ast.StarExpr:
		ty := v.ty(expr.X)
		ty.Ptr = ty.Ptr + 1
		return ty
	case *ast.Ident:
		name := expr.Name
		pn := primName(name)
		if pn != Undefined {
			return v.newType(&Primitive{Base: pn})
		}
		return v.newType(&Ref{Name: expr.Name})
	case *ast.ArrayType:
		return v.newType(&Array{Type: v.ty(expr.Elt)})
	case *ast.MapType:
		return v.newType(&Map{
			KeyType:   v.ty(expr.Key),
			ValueType: v.ty(expr.Value),
		})
	case *ast.FuncType:
		return v.newType(v.fun(expr))
	case *ast.StructType:
		var f []*Decl
		if expr.Fields.NumFields() > 0 {
			f = v.fields(expr.Fields.List)
		}
		return v.newType(&Struct{Fields: f})
	case *ast.InterfaceType:
		var f []*Decl
		if expr.Methods.NumFields() > 0 {
			f = v.fields(expr.Methods.List)
		}
		return v.newType(&Interface{Methods: f})
	case *ast.Ellipsis:
		return v.newType(&Ellipse{Type: v.ty(expr.Elt)})
	case *ast.SelectorExpr:
		return v.newType(&Ref{
			Package: pkg(expr.X),
			Name:    expr.Sel.Name,
		})
	default:
		log.Println("**** unknown type:", reflect.TypeOf(expr).String(), expr)
		return nil
	}
}

func (v *fileVisitor) fun(fun *ast.FuncType) *Function {
	var args []*Decl
	var results []*Decl

	if fun.Params.NumFields() > 0 {
		args = v.fields(fun.Params.List)
	}
	if fun.Results.NumFields() > 0 {
		results = v.fields(fun.Results.List)
	}
	return &Function{
		Args:    args,
		Results: results,
	}
}

func (v *fileVisitor) fields(fields []*ast.Field) []*Decl {
	var decls []*Decl
	for _, field := range fields {
		if len(field.Names) == 0 {
			decls = append(decls, &Decl{Type: v.ty(field.Type)})
			continue
		}
		for _, name := range field.Names {
			decls = append(decls, &Decl{
				Name: name.Name,
				Type: v.ty(field.Type),
			})
		}
	}
	return decls
}

func (v *fileVisitor) litTy(expr ast.Expr) *Type {
	switch expr := expr.(type) {
	case *ast.CompositeLit:
		return v.ty(expr.Type)
	case *ast.BasicLit:
		return v.newType(&Primitive{Base: basicLit(expr.Kind)})
	case *ast.FuncLit:
		return v.ty(expr.Type)
	}
	return nil
}

func (v *fileVisitor) newType(t interface{}) *Type {
	ty := newType(t)
	if ty == nil {
		return nil
	}
	ty.pkg = v.Package
	ty.imports = v.Scope.Impoprts
	return ty
}

func basicLit(kind token.Token) PrimName {
	switch kind {
	case token.INT: // 12345
		return Int
	case token.FLOAT: // 123.45
		return Float64
	case token.IMAG: // 123.45i
		return Complex64
	case token.CHAR: // 'a'
		return Rune
	case token.STRING: // "abc"
		return String
	}
	return Undefined
}

func pkg(expr ast.Expr) string {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr.Name
	}
	return ""
}
