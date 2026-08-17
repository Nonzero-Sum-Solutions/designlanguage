package gen

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	klog "github.com/go-kit/kit/log"

	"github.com/Nonzero-Sum-Solutions/appkit/misc"
	"github.com/Nonzero-Sum-Solutions/designlanguage/model"
	"github.com/Nonzero-Sum-Solutions/designlanguage/parser"
)

type DesignList struct {
	designs []model.Design
}

func NewDesignList() *DesignList {
	return &DesignList{[]model.Design{}}
}

func (d *DesignList) Add(design model.Design) {
	d.designs = append(d.designs, design)
}

func (d *DesignList) All() []model.Design {
	return d.designs
}

// Collect all design objects from design files.
// designPath is the root dir of designs.
// path is the path to the specific design file being parsed.
func HandleDLMFile(logger klog.Logger, designParser parser.Parser, parsedDesigns *DesignList, designPath, path string, info fs.FileInfo, dryRun bool, outputPath string,
	err error) error {
	if err != nil {
		return err
	}

	namespace, _, err2 := ParseDLMFilePath(designPath, path)
	if err2 != nil {
		return err2
	}

	misc.LogMessage(logger, fmt.Sprintf("Parsing design (%s) / (%s)...", designPath, path))

	design, err3 := designParser.Parse(path, namespace)
	if err3 != nil {
		misc.LogMessage(logger, fmt.Sprintf("TODO Parse design. Parse error: [%s]", err3.Error()))
		return err3
	}
	parsedDesigns.Add(design)

	return nil
}

func ParseDLMFilePath(designPath, path string) (namespace string, fileName string, err error) {
	relativePath, err := filepath.Rel(designPath, path)
	if err != nil {
		return "", "", fmt.Errorf("no relative: %s, %s", designPath, path)
	}

	parts := strings.Split(relativePath, string(filepath.Separator))

	if len(parts) < 1 {
		return "", "", fmt.Errorf("invalid path: %s", path)
	}

	if len(parts) == 1 {
		return "", parts[0], nil
	}

	return strings.Join(parts[:len(parts)-1], "/"), parts[len(parts)-1], nil
}

func RenderDesignSummary(design model.Design) string {
	return fmt.Sprintf(
		"Namespace: %s, All Components: %d, Base Components: %d, Entities: %d, Objects: %d",
		design.Namespace(), len(design.AllComponents()), len(design.BaseComponents()),
		len(design.Entities()), len(design.Objects()),
	)
}

func LoadDesigns(logger klog.Logger, projectDir string, dryRun bool) (string, *DesignList, error) {
	misc.LogMessage(logger, "Loading designs...")

	projectDirPath, err := filepath.Abs(projectDir)
	if err != nil {
		return "", nil, err
	}
	designPath := filepath.Join(projectDirPath, "documentation", "design")

	designDirInfo, err := os.Stat(designPath)
	if err != nil {
		return "", nil, err
	}

	if !designDirInfo.IsDir() {
		return "", nil, fmt.Errorf("File is not dir: %s", designPath)
	}

	parsedDesigns := NewDesignList()

	misc.LogMessage(logger, fmt.Sprintf("Walking the path [%q]...", designPath))

	designParser := parser.NewParser()

	err2 := filepath.Walk(designPath, func(path string, info fs.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".nzsd.txt") {
			return nil
		}

		return HandleDLMFile(logger, designParser, parsedDesigns, designPath, path, info, dryRun, projectDirPath, err)
	})

	if err2 != nil {
		misc.LogMessage(logger, fmt.Sprintf("Walking the path. err2.error(): [%s]", err2.Error()))
		return "", nil, fmt.Errorf("Error walking the path [%s]: Cause: [%+v]", designPath, err2)
	}
	return projectDirPath, parsedDesigns, nil
}

func WriteCode(logger klog.Logger, parsedDesigns *DesignList, codeWriter CodeWriter, dryRun bool, projectDirPath string, extension string) error {
	misc.LogMessage(logger, "Generating code...")

	for _, design := range parsedDesigns.All() {

		misc.LogMessage(logger, fmt.Sprintf("Design summary:\n%s", RenderDesignSummary(design)))

		designSource, err := codeWriter.GetSource(design)

		if err != nil {
			misc.LogMessage(logger, fmt.Sprintf("Error rendering design source: %v", err))
			continue
		}

		misc.LogMessage(logger, fmt.Sprintf("Design source: %s", designSource))

		if dryRun {
			continue
		}

		dirPath := filepath.Join(projectDirPath, "model", "gen", design.Namespace())
		filePath := filepath.Join(dirPath, fmt.Sprintf("%s."+extension, design.Namespace()))

		if err = os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}

		if err = os.WriteFile(filePath, []byte(designSource), 0644); err != nil {
			return err
		}

		misc.LogMessage(logger, fmt.Sprintf("Wrote file: %s", filePath))
	}
	return nil
}

type CodeWriter interface {
	GetSource(design model.Design) (string, error)
}
