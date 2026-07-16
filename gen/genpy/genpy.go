package genpy

import (
	klog "github.com/go-kit/kit/log"

	"github.com/mattmunz/designlanguage/gen"
	"github.com/mattmunz/designlanguage/model"
	"github.com/mattmunz/designlanguage/model/gen/pyast"
	pyastm "github.com/mattmunz/designlanguage/model/pyast"
	"github.com/mattmunz/designlanguage/translator/pytranslator"
	"github.com/mattmunz/designlanguage/writer/pywriter"
)

func RenderDesignPySource(design model.Design) (string, error) {
	ast1, s, err := pytranslator.TranslateDesign(design)
	if err != nil {
		return s, err
	}

	text, err := pywriter.RenderPy(ast1, "", "")
	if err != nil {
		return "", err
	}
	return text, nil
}

func RenderComponentPySource(component model.Component) (string, error) {
	var ast1 pyast.AST = pyastm.NewAST()
	pytranslator.AddBaseComponentPy(ast1, component)
	return pywriter.RenderPy(ast1, "", "")
}

func RenderEntityPySource(entity model.Entity) (string, error) {
	var ast1 pyast.AST = pyastm.NewAST()
	pytranslator.AddEntityPy(ast1, entity)
	return pywriter.RenderPy(ast1, "", "")
}

func RenderObjectPySource(obj model.Object) (string, error) {
	var ast1 pyast.AST = pyastm.NewAST()
	if err := pytranslator.AddObjectPy(ast1, obj); err != nil {
		return "", err
	}
	return pywriter.RenderPy(ast1, "", "")
}

func GeneratePySourceForDL(projectDir string, logger klog.Logger, dryRun bool) error {
	projectDirPath, parsedDesigns, err := gen.LoadDesigns(logger, projectDir, dryRun)
	if err != nil {
		return err
	}

	if err = gen.WriteCode(logger, parsedDesigns, NewPyCodeWriter(), dryRun, projectDirPath, "py"); err != nil {
		return err
	}
	return nil
}

type pyCodeWriter struct {
}

func (_ *pyCodeWriter) GetSource(design model.Design) (string, error) {
	return RenderDesignPySource(design)
}

func NewPyCodeWriter() gen.CodeWriter {
	return &pyCodeWriter{}
}
