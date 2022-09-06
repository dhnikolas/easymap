package easymap

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
)

func CopyStruct(source ProcessFile, newName string) (*ast.File, error) {
	
	files, _, err := GetPackageFiles(source)
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
	var decls []ast.Decl
	structType, d := getStructFromFiles(files, structName)
	if structType == nil {
		return decls
	}
	decls = append(decls, d...)
	
	for _, f := range structType.Fields.List {
		if !f.Names[0].IsExported() {
			continue
		}
		if f.Tag != nil {
			f.Tag.Value = ""
		}
		field := detectFieldCategory(f)
		if field == nil {
			continue
		}
		if field.PrefixType != PrefixTypeSimple {
			decls = append(decls, GetStruct(files, field.FieldType)...)
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
