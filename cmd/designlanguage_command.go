package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Nonzero-Sum-Solutions/appkit"
	model "github.com/Nonzero-Sum-Solutions/appkit/model"
	gmodel "github.com/Nonzero-Sum-Solutions/appkit/model/gen/appkit"
)

const version = "1.1.1"

var cli = newCLI()

func newCLI() gmodel.CLI {
	f := newCommandFactory()
	c := appkit.NewCLI(
		model.NewApp("designlanguage", version, ".sundries"),
		"designlanguage",
		"Tools for the Nonzero Sum Design Language.",
		f,
	)

	f.setCLI(c)

	return c
}

type cmdFactory struct {
	cli gmodel.CLI
}

func newCommandFactory() *cmdFactory {
	return &cmdFactory{}
}

func (f *cmdFactory) setCLI(i gmodel.CLI) {
	f.cli = i
}

func (f *cmdFactory) New() *cobra.Command {
	cmd := appkit.NewCommandBase(f.cli)
	cmd.AddCommand(newGenCmd(f.cli))
	return cmd
}

func init() {
	appkit.DoInit(cli)
}

func CLI() gmodel.CLI {
	return cli
}
