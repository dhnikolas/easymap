package easymap

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path"
)

const (
	PREFIX_TYPE_POINTER PrefixType = 0
	PREFIX_TYPE_SLICE   PrefixType = 1
)

type PrefixType uint

type ProcessFile struct {
	FullPath   string
	StructName string
	FieldName  string
}

type Content struct {
	*StructField
	Content       string
	ParentInName  string
	ParentOutName string
	InStructType  string
}

type Repository struct {
	typeSpec   *ast.TypeSpec
	structType *ast.StructType
}

type StructField struct {
	Name       string
	NameIn     string
	StructType string

	PrefixType   PrefixType
	ParentStruct *StructField

	ListScalarFields []*ScalarField
	ListStructFields []*StructField
}

type ScalarField struct {
	Name      string
	FieldType string
}

func GetPackageFiles(source ProcessFile) ([]*ast.File, error) {
	dirPath := path.Dir(source.FullPath)
	fmt.Println(dirPath)

	packageName, err := GetPackageName(source.FullPath)
	if err != nil {
		return nil, err
	}

	pkgs, err := parser.ParseDir(token.NewFileSet(), dirPath, func(info fs.FileInfo) bool {
		return true
	}, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("Error parse dir %s ", err)
	}

	currentPackage, ok := pkgs[packageName]
	if !ok {
		return nil, fmt.Errorf("Package not found %s ", packageName)
	}
	var files []*ast.File
	for _, f := range currentPackage.Files {
		files = append(files, f)
	}

	return files, nil
}
