package cmd

import (
	klog "github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/Nonzero-Sum-Solutions/appkit/cmd"
	gmodel "github.com/Nonzero-Sum-Solutions/appkit/model/gen/appkit"
	"github.com/Nonzero-Sum-Solutions/designlanguage/gen/gengo"
	"github.com/Nonzero-Sum-Solutions/designlanguage/gen/genpy"
)

func newGenCmd(cli gmodel.CLI) *cobra.Command {
	genCmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate model source code from Design Language Model files.",
		Args:  validateGenArgs,
	}

	dryRunFlag := genCmd.Flags().BoolP("dryRun", "d", true, "True to read but not make changes on disk, true by default.")
	projectDirFlag := genCmd.Flags().StringP("projectDir", "p", "",
		"Path to the project directory.")
	targetFlag := genCmd.Flags().StringP("target", "t", "",
		"Target language for the generated source code.")

	doDLM1 := func(cli gmodel.CLI, cmd *cobra.Command, args []string) error {
		return doGen(cli.Logger(), cmd, args, *dryRunFlag, *projectDirFlag, *targetFlag)
	}
	genCmd.Run = cmd.WrapRunner(doDLM1, cli)

	return genCmd
}

func validateGenArgs(_ *cobra.Command, args []string) error {
	if len(args) > 0 {
		return errors.New("No args expected")
	}
	return nil
}

func doGen(logger klog.Logger, cmd *cobra.Command, args []string, dryRun bool, projectDir, target string) error {
	err := validateGenArgs(cmd, args)
	if err != nil {
		return err
	}

	switch target {
	case "go", "golang":
		return gengo.GenerateGoSourceForDL(projectDir, logger, dryRun)
	case "python", "py":
		return genpy.GeneratePySourceForDL(projectDir, logger, dryRun)
	default:
		return errors.New("Invalid target language: " + target)
	}
}
