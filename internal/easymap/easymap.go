package easymap

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"path"
	
	"github.com/dhnikolas/easymap/pkg/gofmt"
)

const (
	PrefixTypePointer PrefixType = 0
	PrefixTypeSlice   PrefixType = 1
	PrefixTypeStruct  PrefixType = 2
	PrefixTypeSimple  PrefixType = 3
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
	
	ListSimpleFields []*SimpleField
	ListStructFields []*StructField
}

type SimpleField struct {
	Name      string
	FieldType string
}

type CommonField struct {
	Name       string
	FieldType  string
	PrefixType PrefixType
}

func GetPackageFiles(source ProcessFile) ([]*ast.File, string, error) {
	dirPath := path.Dir(source.FullPath)
	
	packageName, err := GetPackageName(source.FullPath)
	if err != nil {
		return nil, "", err
	}
	
	pkgs, err := parser.ParseDir(token.NewFileSet(), dirPath, func(info fs.FileInfo) bool {
		return true
	}, parser.ParseComments)
	if err != nil {
		return nil, "", fmt.Errorf("Error parse dir %s ", err)
	}
	
	currentPackage, ok := pkgs[packageName]
	if !ok {
		return nil, "", fmt.Errorf("Package not found %s ", packageName)
	}
	var files []*ast.File
	for _, f := range currentPackage.Files {
		files = append(files, f)
	}
	
	return files, packageName, nil
}

func FileToString(file *ast.File) (string, error) {
	var bResult bytes.Buffer
	err := printer.Fprint(&bResult, token.NewFileSet(), file)
	if err != nil {
		return "", err
	}
	
	return bResult.String(), nil
}

func GoFmt(body []byte) ([]byte, error) {
	r := bytes.NewReader(body)
	var bResult bytes.Buffer
	err := gofmt.ProcessGofmt("new.go", r, &bResult, false)
	if err != nil {
		return nil, err
	}
	
	return bResult.Bytes(), nil
}
