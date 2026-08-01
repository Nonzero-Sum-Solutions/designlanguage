package parser

import (
	"fmt"

	"github.com/Nonzero-Sum-Solutions/designlanguage/model"
	p "github.com/Nonzero-Sum-Solutions/designlanguage/model/designlanguage/parser"
	"github.com/antlr4-go/antlr/v4"
)

type antlrParser struct {
}

func NewANTLRParser() Parser {
	return &antlrParser{}
}

func (*antlrParser) Parse(path, namespace string) (model.Design, error) {
	fileBytes, err := readFile(path)
	if err != nil {
		return nil, err
	}
	lexer := p.NewDesignLanguageLexer(antlr.NewInputStream(string(fileBytes)))
	parser := p.NewDesignLanguageParser(antlr.NewCommonTokenStream(lexer, 0))

	components := []model.Component{}
	entities := []model.Entity{}
	objects := []model.Object{}
	design := parser.Design()
	author, comment := "", ""
	preamble := design.Preamble()
	if preamble != nil {
		authorElem := preamble.Author()
		if authorElem != nil {
			author = authorElem.GetText()
		}

		commentElem := preamble.COMMENT()
		if commentElem != nil {
			comment = commentElem.GetText()
		}
	}

	for _, cc := range design.AllComponent() {
		// TODO This is now wrong.
		objects, components, entities, err = addComponentANTLR(cc, objects, components, entities)
		if err != nil {
			return nil, err
		}
	}

	return model.NewDesign(author, comment, namespace, components), nil
}

func addComponentANTLR(cc p.IComponentContext, objects []model.Object, components []model.Component, entities []model.Entity) ([]model.Object, []model.Component, []model.Entity, error) {
	name := cc.SimpleComponent().NAME().GetText()
	objComment := ""
	comment := cc.SimpleComponent().COMMENT()
	if comment != nil {
		objComment = comment.GetText()
	}

	fcs := cc.AllField()
	attributes := []model.Attribute{}
	methods := []model.Method{}
	var err error
	for _, fc := range fcs {
		attr := fc.Attribute()
		if attr != nil {
			attributes, err = addAttribute(attr, attributes)
			if err != nil {
				return nil, nil, nil, err
			}

			continue
		}

		method := fc.Method()
		if method == nil {
			continue
		}

		methods, err = addMethod(method, methods)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	if len(methods) > 0 {
		obj, err := model.NewObject(name, objComment, nil, attributes, methods)
		if err != nil {
			return nil, nil, nil, err
		}
		objects = append(objects, obj)
		components = append(components, obj)
	} else {
		ent, err := model.NewEntity(name, objComment, nil, attributes)
		if err != nil {
			return nil, nil, nil, err
		}
		entities = append(entities, ent)
		components = append(components, ent)
	}
	return objects, components, entities, nil
}

func addAttribute(attr p.IAttributeContext, attributes []model.Attribute) ([]model.Attribute, error) {
	// param, err := GetParam(attr.Param())
	// if err != nil {
	// 	return nil, err
	// }

	//attrComment := ""
	panic("NYI attrComment")

	//return append(attributes, model.NewAttribute(param.Name(), attrComment, param.Type().Name(), param.Type().IsArray())), nil
}

func addMethod(method p.IMethodContext, methods []model.Method) ([]model.Method, error) {
	methodName := method.NAME().GetText()
	pcs := method.AllParams()
	if len(pcs) != 2 {
		return methods, fmt.Errorf("Expected 2 params")
	}

	param1, err := GetParams(pcs[0])
	if err != nil {
		return nil, err
	}

	param2, err := GetParams(pcs[1])
	if err != nil {
		return nil, err
	}

	methodComment := ""
	panic("NYI methodComment")

	newMethod, err := model.NewMethod(methodName, methodComment, param1, param2)
	if err != nil {
		return nil, err
	}
	return append(methods, newMethod), nil
}

func GetParams(params p.IParamsContext) ([]model.Param, error) {
	paramCount := len(params.AllParam())

	if paramCount < 1 {
		return []model.Param{}, nil
	}

	paramsList := make([]model.Param, paramCount)
	for i, param := range params.AllParam() {
		newParam, err := GetParam(param)
		if err != nil {
			return nil, err
		}

		paramsList[i] = newParam
	}

	return paramsList, nil
}

func GetParam(param p.IParamContext) (model.Param, error) {
	if param.NAME() == nil {
		return nil, fmt.Errorf("Expected param name")
	}

	attrName := param.NAME().GetText()

	attrType := param.Type_()
	if attrType == nil {
		return nil, fmt.Errorf("Expected param type")
	}

	typeName := attrType.NAME().GetText()

	typeArray := attrType.ARRAY()
	t, err := model.NewType(typeName, false)
	if err != nil {
		return nil, err
	}
	if typeArray == nil {
		return model.NewParam(attrName, t), nil
	}

	isArray := len(typeArray.GetText()) > 0
	t, err = model.NewType(typeName, isArray)
	return model.NewParam(attrName, t), err
}
