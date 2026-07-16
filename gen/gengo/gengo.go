package gengo

import (
	"github.com/dave/jennifer/jen"
	klog "github.com/go-kit/kit/log"

	"github.com/mattmunz/designlanguage/gen"
	"github.com/mattmunz/designlanguage/model"
	"github.com/mattmunz/designlanguage/translator/gotranslator"
	"github.com/mattmunz/designlanguage/writer/gowriter"
)

func RenderDesignGoSource(design model.Design) (string, error) {
	outFile, s, err := gotranslator.TranslateDesign(design)
	if err != nil {
		return s, err
	}

	return gowriter.WriteGoFile(outFile)
}

func RenderEnumSource(enumExpr model.Enum) (string, error) {
	stmt := jen.Empty()
	gotranslator.AddEnum(stmt, enumExpr)
	return gowriter.WriteGo(stmt)
}

func RenderComponentSource(component model.Component) (string, error) {
	stmt := jen.Empty()
	gotranslator.AddBaseComponent(stmt, component)
	return gowriter.WriteGo(stmt)
}

func RenderEntitySource(entity model.Entity) (string, error) {
	stmt := jen.Empty()
	gotranslator.AddEntity(stmt, entity)
	return gowriter.WriteGo(stmt)
}

func RenderObjectSource(obj model.Object) (string, error) {
	stmt := jen.Empty()
	gotranslator.AddObject(stmt, obj)
	return gowriter.WriteGo(stmt)
}

func GenerateGoSourceForDL(projectDir string, logger klog.Logger, dryRun bool) error {
	projectDirPath, parsedDesigns, err := gen.LoadDesigns(logger, projectDir, dryRun)
	if err != nil {
		return err
	}

	if err = gen.WriteCode(logger, parsedDesigns, NewGoCodeWriter(), dryRun, projectDirPath, "go"); err != nil {
		return err
	}
	return nil
}

type goCodeWriter struct {
}

func (_ *goCodeWriter) GetSource(design model.Design) (string, error) {
	return RenderDesignGoSource(design)
}

func NewGoCodeWriter() gen.CodeWriter {
	return &goCodeWriter{}
}
