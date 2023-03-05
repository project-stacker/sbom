package sbom

import (
	"os"

	"github.com/spf13/cobra"
	"stackerbuild.io/sbom/pkg/distro"
)

func GenerateCmd() *cobra.Command {
	input := ""
	output := ""
	format := ""

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate",
		Long:  "generate",
		Run: func(cmd *cobra.Command, args []string) {
			distro.ParseFile(input)
		},
	}

	cmd.Flags().StringVarP(&input, "input", "i", "", "input file")
	cmd.MarkFlagRequired("input")
	cmd.Flags().StringVarP(&output, "output", "o", "", "output file")
	cmd.Flags().StringVarP(&format, "format", "f", "spdx", "output format (spdx, default:spdx)")

	return cmd
}

func BuildCmd() *cobra.Command {
	input := ""
	cmd := &cobra.Command{
		Use:   "build",
		Short: "build",
		Long:  "build",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	cmd.Flags().StringVarP(&input, "input", "i", "", "input file")
	cmd.MarkFlagRequired("input")

	return cmd
}

func NewCli() *cobra.Command {
	showVersion := false

	cmd := &cobra.Command{
		Use:   "sbom",
		Short: "sbom",
		Long:  `A SBOM generator tool`,
		Run: func(cmd *cobra.Command, args []string) {
			if showVersion {
				os.Exit(0)
			} else {
				_ = cmd.Usage()
			}
		},
	}

	cmd.AddCommand(BuildCmd())
	cmd.AddCommand(GenerateCmd())
	cmd.Flags().BoolVarP(&showVersion, "version", "v", false, "show the version and exit")

	return cmd
}
