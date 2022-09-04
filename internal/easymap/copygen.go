package easymap

import (
	"fmt"
	"go/ast"
)

func CopyGen(processFile ProcessFile, newStructName string) (string, string, error) {

	copyFile, err := CopyStruct(processFile, newStructName)
	if err != nil {
		return "", "", err
	}
	var currentStructName string
	if len(newStructName) > 0 {
		currentStructName = newStructName
	} else {
		currentStructName = processFile.StructName
	}

	inFileStruct, err := Scan(processFile)
	if err != nil {
		return "", "", fmt.Errorf("In File error %s ", err)
	}

	copyStruct, err := ScanStruct([]*ast.File{copyFile}, currentStructName, "")
	if err != nil {
		return "", "", err
	}

	mapping, err := GenerateMapping(inFileStruct, copyStruct)
	if err != nil {
		return "", "", err
	}

	copyFileString, err := FileToString(copyFile)
	if err != nil {
		return "", "", err
	}

	return copyFileString, string(mapping), nil
}
