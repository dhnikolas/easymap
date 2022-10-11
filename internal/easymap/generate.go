package easymap

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"strings"
	"text/template"
	
	"golang.org/x/tools/go/ast/inspector"
)

func GenerateMapping(inFile, outFile *StructField) ([]byte, error) {
	commonStruct := GetCommonStruct(outFile, inFile)
	result := GenerateMainTemplate(commonStruct, inFile.StructType)
	r, err := GoFmt(result)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func Scan(source ProcessFile) (*StructField, string, error) {
	files, packageName, err := GetPackageFiles(source)
	if err != nil {
		return nil, "", err
	}
	
	structField, err := ScanStruct(files, source.StructName, source.FieldName)
	if err != nil {
		return nil, "", err
	}
	return structField, packageName, nil
}

func ScanStruct(files []*ast.File, structName, fieldName string) (*StructField, error) {
	s, _ := getStructFromFiles(files, structName)
	if s == nil {
		return nil, errors.New("Struct not exist " + structName)
	}
	currentStruct := &StructField{}
	currentStruct.Name = structName
	currentStruct.NameIn = fieldName
	currentStruct.StructType = structName
	
	for _, f := range s.Fields.List {
		if !f.Names[0].IsExported() {
			
			continue
		}
		field := detectFieldCategory(f)
		if field == nil {
			continue
		}
		switch field.PrefixType {
		case PrefixTypeSlice:
			fallthrough
		case PrefixTypeStruct:
			fallthrough
		case PrefixTypePointer:
			newStruct, err := ScanStruct(files, field.FieldType, field.Name)
			if err != nil {
				continue
			}
			newStruct.PrefixType = field.PrefixType
			newStruct.ParentStruct = currentStruct
			currentStruct.ListStructFields = append(currentStruct.ListStructFields, newStruct)
		case PrefixTypeSimple:
			currentStruct.ListSimpleFields = append(currentStruct.ListSimpleFields, &SimpleField{
				Name:      field.Name,
				FieldType: field.FieldType,
			})
		}
	}
	
	return currentStruct, nil
}

func getStructFromFiles(files []*ast.File, structName string) (*ast.StructType, []ast.Decl) {
	var s *ast.StructType
	var decls []ast.Decl
	i := inspector.New(files)
	iFilter := []ast.Node{
		&ast.GenDecl{},
	}
	i.Nodes(iFilter, func(node ast.Node, push bool) (proceed bool) {
		genDecl := node.(*ast.GenDecl)
		if genDecl == nil {
			return false
		}
		typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec)
		if !ok {
			return false
		}
		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return false
		}
		if typeSpec.Name.Name != structName {
			return false
		}
		s = structType
		decls = append(decls, genDecl)
		return false
	})
	
	return s, decls
}

func detectFieldCategory(f *ast.Field) *CommonField {
	var resultField *CommonField
	switch ident := f.Type.(type) {
	case *ast.Ident:
		switch {
		case ident.Obj != nil && ident.Obj.Kind == ast.Typ:
			resultField = &CommonField{
				Name:       f.Names[0].Name,
				FieldType:  fmt.Sprint(ident.Obj.Name),
				PrefixType: PrefixTypeStruct,
			}
		default:
			resultField = &CommonField{
				Name:       f.Names[0].Name,
				FieldType:  ident.Name,
				PrefixType: PrefixTypeSimple,
			}
		}
	
	case *ast.MapType:
		resultField = &CommonField{
			Name:       f.Names[0].Name,
			FieldType:  fmt.Sprintf("map[%s]%s", ident.Key, ident.Value),
			PrefixType: PrefixTypeSimple,
		}
	case *ast.StarExpr:
		switch expr := ident.X.(type) {
		case *ast.SelectorExpr:
			resultField = &CommonField{
				Name:       f.Names[0].Name,
				FieldType:  fmt.Sprintf("*%s.%s", fmt.Sprint(expr.X), expr.Sel.Name),
				PrefixType: PrefixTypeSimple,
			}
		default:
			resultField = &CommonField{
				Name:       f.Names[0].Name,
				FieldType:  fmt.Sprint(ident.X),
				PrefixType: PrefixTypePointer,
			}
		}
	
	case *ast.ArrayType:
		switch arrayType := ident.Elt.(type) {
		case *ast.Ident:
			resultField = &CommonField{
				Name:       f.Names[0].Name,
				FieldType:  arrayType.Name,
				PrefixType: PrefixTypeSimple,
			}
		
		case *ast.StarExpr:
			
			resultField = &CommonField{
				Name:       f.Names[0].Name,
				FieldType:  fmt.Sprint(arrayType.X),
				PrefixType: PrefixTypeSlice,
			}
		}
	}
	
	return resultField
}

func GenerateMainTemplate(s *StructField, inStructType string) []byte {
	if s == nil {
		return []byte{}
	}
	
	c := Content{StructField: s, InStructType: inStructType}
	content := ""
	for _, node := range c.ListStructFields {
		content += GenerateCheckTemplate(node, "out", "in")
	}
	
	c.Content = content
	
	templ := mappingTemplate
	
	var b bytes.Buffer
	t := template.Must(template.New("").Parse(templ))
	if err := t.Execute(&b, c); err != nil {
		panic(err)
	}
	
	return b.Bytes()
}

func GenerateCheckTemplate(s *StructField, parentOutName, parentInName string) string {
	if s == nil {
		return ""
	}
	
	if len(parentOutName) > 0 {
		parentOutName = parentOutName + "."
	}
	if len(parentInName) > 0 {
		parentInName = parentInName + "."
	}
	
	var newParentInName, newParentOutName string
	
	compareStruct := s
	if s.ParentStruct != nil {
		compareStruct = s.ParentStruct
	}
	
	//Parent params
	switch compareStruct.PrefixType {
	case PrefixTypePointer:
		fallthrough
	case PrefixTypeStruct:
		newParentInName = parentInName + s.NameIn
		newParentOutName = parentOutName + s.NameIn
	
	case PrefixTypeSlice:
		newParentInName = s.ParentStruct.NameIn + "Item." + s.NameIn
		newParentOutName = "new" + s.ParentStruct.NameIn + "." + s.NameIn
		
		parentInName = compareStruct.NameIn + "Item."
		parentOutName = "new" + compareStruct.NameIn + "."
	}
	
	c := Content{StructField: s, ParentOutName: parentOutName, ParentInName: parentInName}
	
	content := ""
	
	for _, node := range c.ListStructFields {
		content += GenerateCheckTemplate(node, newParentOutName, newParentInName)
	}
	
	c.Content = content
	
	templ := ""
	switch s.PrefixType {
	case PrefixTypePointer:
		templ = ifConditionPointerTemplate
	
	case PrefixTypeSlice:
		templ = slicePointerTemplate
	
	case PrefixTypeStruct:
		templ = structTemplate
	}
	
	var b bytes.Buffer
	t := template.Must(template.New("").Parse(templ))
	if err := t.Execute(&b, c); err != nil {
		panic(err)
	}
	
	return b.String()
}

func GetCommonStruct(outStruct *StructField, inStruct *StructField) *StructField {
	resultStruct := &StructField{}
	
	resultStruct.Name = outStruct.Name
	resultStruct.NameIn = inStruct.NameIn
	resultStruct.StructType = outStruct.StructType
	resultStruct.PrefixType = outStruct.PrefixType
	resultStruct.ParentStruct = outStruct.ParentStruct
	
	for _, scalarField := range outStruct.ListSimpleFields {
		for _, inField := range inStruct.ListSimpleFields {
			if strings.ToLower(scalarField.Name) == strings.ToLower(inField.Name) &&
				strings.ToLower(scalarField.FieldType) == strings.ToLower(inField.FieldType) {
				resultStruct.ListSimpleFields = append(resultStruct.ListSimpleFields, scalarField)
			}
		}
	}
	
	for _, outFieldStruct := range outStruct.ListStructFields {
		for _, inFieldStruct := range inStruct.ListStructFields {
			if strings.ToLower(outFieldStruct.NameIn) == strings.ToLower(inFieldStruct.NameIn) {
				resultStruct.ListStructFields = append(resultStruct.ListStructFields, GetCommonStruct(outFieldStruct, inFieldStruct))
			}
		}
	}
	
	return resultStruct
}
