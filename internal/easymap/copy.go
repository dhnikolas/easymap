package easymap

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"

	"golang.org/x/tools/go/ast/inspector"
)

func CopyStruct(source ProcessFile, newName string) {
	decls := GetStruct(source)
	resultFile := &ast.File{Name: &ast.Ident{
		NamePos: 0,
		Name:    "main",
	}}

	for _, decl := range decls {
		genDecl := decl.(*ast.GenDecl)
		if genDecl == nil {
			continue
		}

		typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec)
		if !ok {
			continue
		}

		if len(newName) > 0 {
			if typeSpec.Name.Name == source.StructName {
				typeSpec.Name.Name = newName
			}
		}

		resultFile.Decls = append(resultFile.Decls, decl)
	}

	var bResult bytes.Buffer
	printer.Fprint(&bResult, token.NewFileSet(), resultFile)
	fmt.Println(bResult.String())
}

func GetStruct(source ProcessFile) []ast.Decl {
	astInFile, err := parser.ParseFile(token.NewFileSet(), source.FullPath, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("parse file: %v", err)
	}

	ast.FileExports(astInFile)

	i := inspector.New([]*ast.File{astInFile})
	iFilter := []ast.Node{
		&ast.GenDecl{},
	}

	var structType *ast.StructType

	var decls []ast.Decl
	i.Nodes(iFilter, func(node ast.Node, push bool) (proceed bool) {
		genDecl := node.(*ast.GenDecl)
		if genDecl == nil {
			return false
		}
		typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec)
		if !ok {
			return false
		}
		st, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return false
		}
		if typeSpec.Name.Name != source.StructName {
			return false
		}
		structType = st

		decls = append(decls, genDecl)

		return false
	})

	if structType == nil {
		return decls
	}

	for _, f := range structType.Fields.List {
		if !f.Names[0].IsExported() {
			continue
		}
		if f.Tag != nil {
			f.Tag.Value = ""
		}
		switch ident := f.Type.(type) {
		case *ast.StarExpr:
			processFile := ProcessFile{
				FullPath:   source.FullPath,
				StructName: fmt.Sprint(ident.X),
				FieldName:  f.Names[0].Name,
			}
			decls = append(decls, GetStruct(processFile)...)

		case *ast.ArrayType:
			switch arrayType := ident.Elt.(type) {
			case *ast.StarExpr:
				processFile := ProcessFile{
					FullPath:   source.FullPath,
					StructName: fmt.Sprint(arrayType.X),
					FieldName:  f.Names[0].Name,
				}
				decls = append(decls, GetStruct(processFile)...)
			}
		}
	}

	return decls
}
