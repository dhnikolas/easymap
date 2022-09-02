package easymap

import "go/ast"

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
