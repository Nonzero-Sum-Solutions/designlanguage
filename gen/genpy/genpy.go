package genpy

import (
	klog "github.com/go-kit/kit/log"

	"github.com/Nonzero-Sum-Solutions/designlanguage/gen"
	"github.com/Nonzero-Sum-Solutions/designlanguage/model"
	"github.com/Nonzero-Sum-Solutions/designlanguage/model/gen/pyast"
	pyastm "github.com/Nonzero-Sum-Solutions/designlanguage/model/pyast"
	"github.com/Nonzero-Sum-Solutions/designlanguage/translator/pytranslator"
	"github.com/Nonzero-Sum-Solutions/designlanguage/writer/pywriter"
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
