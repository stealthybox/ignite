package vmcmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdKill kills running VMs
func NewCmdKill(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kill <vm>...",
		Short: "Kill running VMs",
		Long: dedent.Dedent(`
			Kill (force stop) one or multiple VMs. The VMs are matched by prefix based
			on their ID and name. To kill multiple VMs, chain the matches separated
			by spaces.
		`),
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				so, err := (&run.StopFlags{Kill: true}).NewStopOptions(args)
				if err != nil {
					return err
				}

				return run.Stop(so)
			}())
		},
	}

	return cmd
}
