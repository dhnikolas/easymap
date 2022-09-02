package easymap

import (
	"fmt"
	"go/ast"
	"go/token"
)

func CopyGen(processFile ProcessFile, newStructName string) (*ast.File, error) {

	copyFile, err := CopyStruct(processFile, newStructName)
	if err != nil {
		return nil, err
	}
	var currentStructName string
	if len(newStructName) > 0 {
		currentStructName = newStructName
	} else {
		currentStructName = processFile.StructName
	}

	inFileStruct, err := Scan(processFile)
	if err != nil {
		return nil, fmt.Errorf("In File error %s ", err)
	}

	copyStruct, err := ScanStruct([]*ast.File{copyFile}, currentStructName, "")
	if err != nil {
		return nil, err
	}

	mapping, err := Generate(inFileStruct, copyStruct)
	if err != nil {
		return nil, err
	}

	for _, decl := range mapping.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if ok {
			if genDecl.Tok == token.IMPORT {
				continue
			}
		}

		copyFile.Decls = append(copyFile.Decls, decl)
	}

	return copyFile, nil
}
