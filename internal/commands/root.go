package commands

import (
	"github.com/spf13/cobra"
	"realess-server/internal/commands/create"
)

func CreateRootCmd(sersion string) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:     "rlss",
		Version: sersion,
		Short:   "Realess Server is a server for layer3 IP packet forwarding.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(create.CreateCreateCmd())
	// rootCmd.AddCommand(delete.CreateDeleteCmd())
	// rootCmd.AddCommand(uninstall.CreateUninstallCmd())

	rootCmd.Flags().BoolP("version", "v", false, "version of rlss")

	return rootCmd
}
