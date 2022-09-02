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

	outArgument := cCtx.Args().Get(0)
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

	result, err := easymap.Generate(inFile, outFile)
	if err != nil {
		return err
	}

	fmt.Println(result)

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

	result, err := easymap.CopyStruct(easymap.ProcessFile{
		FullPath:   sourceFilePath,
		StructName: sourceStructName,
	}, newStructName)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
