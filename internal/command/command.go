package command

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dhnikolas/easymap/internal/easymap"
	"github.com/urfave/cli/v2"
)

type Application struct {
}

func (a *Application) Gen(cCtx *cli.Context) error {
	sourceArgument := cCtx.Args().Get(0)
	sourceArgumentSplit := strings.Split(sourceArgument, ":")
	if len(sourceArgumentSplit) < 2 {
		return errors.New("Wrong source argument " + sourceArgument)
	}

	outArgument := cCtx.Args().Get(1)
	outArgumentSplit := strings.Split(outArgument, ":")
	if len(outArgumentSplit) < 2 {
		return errors.New("Wrong out argument " + outArgument)
	}
	sourceFilePath := sourceArgumentSplit[0]
	sourceStructName := sourceArgumentSplit[1]

	outFilePath := outArgumentSplit[0]
	outStructName := outArgumentSplit[1]

	inFile := easymap.ProcessFile{
		FullPath:   sourceFilePath,
		StructName: sourceStructName,
	}

	outFile := easymap.ProcessFile{
		FullPath:   outFilePath,
		StructName: outStructName,
	}

	inFileStruct, err := easymap.Scan(inFile)
	if err != nil {
		return fmt.Errorf("In File error %s ", err)
	}

	outFileStruct, err := easymap.Scan(outFile)
	if err != nil {
		return fmt.Errorf("Out File error %s ", err)
	}

	result, err := easymap.GenerateMapping(inFileStruct, outFileStruct)
	if err != nil {
		return err
	}

	fmt.Println(string(result))

	return nil
}

func (a *Application) Copy(cCtx *cli.Context) error {
	sourceArgument := cCtx.Args().Get(0)
	sourceArgumentSplit := strings.Split(sourceArgument, ":")
	if len(sourceArgumentSplit) < 2 {
		return errors.New("Wrong source argument " + sourceArgument)
	}

	sourceFilePath := sourceArgumentSplit[0]
	sourceStructName := sourceArgumentSplit[1]

	newStructName := cCtx.Args().Get(1)

	processFile := easymap.ProcessFile{
		FullPath:   sourceFilePath,
		StructName: sourceStructName,
	}
	result, err := easymap.CopyStruct(processFile, newStructName)
	if err != nil {
		return err
	}

	stringResult, err := easymap.FileToString(result)
	if err != nil {
		return err
	}
	fmt.Println(stringResult)
	return nil
}

func (a *Application) CopyGen(cCtx *cli.Context) error {
	sourceArgument := cCtx.Args().Get(0)
	sourceArgumentSplit := strings.Split(sourceArgument, ":")
	if len(sourceArgumentSplit) < 2 {
		return errors.New("Wrong source argument " + sourceArgument)
	}

	sourceFilePath := sourceArgumentSplit[0]
	sourceStructName := sourceArgumentSplit[1]
	newStructName := cCtx.Args().Get(1)
	processFile := easymap.ProcessFile{
		FullPath:   sourceFilePath,
		StructName: sourceStructName,
	}

	copyFile, mappingFile, err := easymap.CopyGen(processFile, newStructName)
	if err != nil {
		return err
	}

	fmt.Print(strings.ReplaceAll(copyFile, "package main", ""))
	fmt.Print(strings.ReplaceAll(mappingFile, "package main", ""))

	return nil
}
