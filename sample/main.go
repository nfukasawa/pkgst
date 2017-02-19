package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"

	"github.com/nfukasawa/pkgst"
)

func main() {

	fset := token.NewFileSet()

	dirname := "./fixture"
	pkgs, err := parser.ParseDir(fset, dirname, func(info os.FileInfo) bool {
		log.Println("file:", info.Name())
		return true
	}, 0)
	if err != nil {
		exit(err)
	}

	visitor := new(pkgst.WalkVisitor)
	for _, pkg := range pkgs {
		ast.Walk(visitor, pkg)
	}
	fmt.Println(marshalJSON(visitor.Packages))
}

func marshalJSON(obj interface{}) string {
	b, err := json.MarshalIndent(obj, "", " ")
	if err != nil {
		exit(err)
	}
	return string(b)
}

func exit(err error) {
	os.Stderr.WriteString(err.Error() + "\n")
	os.Exit(1)
}
