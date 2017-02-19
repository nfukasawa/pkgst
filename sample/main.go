package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"

	"encoding/json"

	"github.com/nfukasawa/pkgst"
)

func main() {

	fset := token.NewFileSet()

	dirname := "./fixture"
	ps, err := parser.ParseDir(fset, dirname, func(info os.FileInfo) bool {
		log.Println("file:", info.Name())
		return true
	}, 0)
	if err != nil {
		panic(err)
	}

	pkgs := pkgst.Build(ps)
	for _, pkg := range pkgs {
		b, _ := json.MarshalIndent(pkg, "", " ")
		fmt.Println(string(b))
	}
}
