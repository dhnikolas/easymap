package easymap

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"path"

	"golang.org/x/tools/go/ast/inspector"
)

func CopyStruct(source ProcessFile, newName string) (string, error) {

	dirPath := path.Dir(source.FullPath)
	fmt.Println(dirPath)

	packageName, err := GetPackageName(source.FullPath)
	if err != nil {
		return "", err
	}

	pkgs, err := parser.ParseDir(token.NewFileSet(), dirPath, func(info fs.FileInfo) bool {
		return true
	}, parser.ParseComments)
	if err != nil {
		return "", fmt.Errorf("Error parse dir %s ", err)
	}

	currentPackage, ok := pkgs[packageName]
	if !ok {
		return "", fmt.Errorf("Package not found %s ", packageName)
	}
	var files []*ast.File
	for _, f := range currentPackage.Files {
		files = append(files, f)
	}
	decls := GetStruct(files, source.StructName)
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

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		var fields []*ast.Field
		for _, field := range structType.Fields.List {
			if !field.Names[0].IsExported() {
				continue
			}
			fields = append(fields, field)
		}

		structType.Fields.List = fields

		resultFile.Decls = append(resultFile.Decls, decl)
	}

	var bResult bytes.Buffer
	err = printer.Fprint(&bResult, token.NewFileSet(), resultFile)
	if err != nil {
		return "", err
	}

	return bResult.String(), nil
}

func GetStruct(files []*ast.File, structName string) []ast.Decl {

	i := inspector.New(files)
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
		if typeSpec.Name.Name != structName {
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
			decls = append(decls, GetStruct(files, fmt.Sprint(ident.X))...)
		case *ast.ArrayType:
			switch arrayType := ident.Elt.(type) {
			case *ast.StarExpr:
				decls = append(decls, GetStruct(files, fmt.Sprint(arrayType.X))...)
			}
		}
	}

	return decls
}

func GetPackageName(filePath string) (string, error) {
	astInFile, err := parser.ParseFile(token.NewFileSet(), filePath, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}
	if astInFile.Name == nil {
		return "", errors.New("No package name in file " + filePath)
	}

	return astInFile.Name.Name, nil
}
