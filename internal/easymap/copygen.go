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
	
	inFileStruct, packageName, err := Scan(processFile)
	if err != nil {
		return "", "", fmt.Errorf("In File error %s ", err)
	}
	
	outFile, err := ScanStruct([]*ast.File{copyFile}, currentStructName, "")
	if err != nil {
		return "", "", err
	}
	
	commonStruct := GetCommonStruct(outFile, inFileStruct)
	result := GenerateMainTemplate(commonStruct, packageName+"."+inFileStruct.StructType)
	mapping, err := GoFmt(result)
	if err != nil {
		return "", "", err
	}
	
	copyFileString, err := FileToString(copyFile)
	if err != nil {
		return "", "", err
	}
	
	return copyFileString, string(mapping), nil
}
