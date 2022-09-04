package easymap

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/ast/inspector"
)

func CopyStruct(source ProcessFile, newName string) (*ast.File, error) {

	files, err := GetPackageFiles(source)
	if err != nil {
		return nil, err
	}
	decls := GetStruct(files, source.StructName)
	resultFile := &ast.File{
		Name: &ast.Ident{
			NamePos: 0,
			Name:    "main",
		},
	}

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
	return resultFile, nil
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
		case *ast.Ident:
			if ident.Obj != nil && ident.Obj.Kind == ast.Typ {
				decls = append(decls, GetStruct(files, fmt.Sprint(ident.Obj.Name))...)
			}
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
