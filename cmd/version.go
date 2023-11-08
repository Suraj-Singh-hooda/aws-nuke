package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// BuildVersion
// BuildDate
// BuildHash
// BuildEnvironment
var (
	//BuildVersion
	BuildVersion = "unknown"
	//BuildDate
	BuildDate = "unknown"
	//BuildHash
	BuildHash = "unknown"
	//BuildEnvironment
	BuildEnvironment = "unknown"
)

// NewVersionCommand function
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "shows version of this application",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("version:     %s\n", BuildVersion)
			fmt.Printf("build date:  %s\n", BuildDate)
			fmt.Printf("scm hash:    %s\n", BuildHash)
			fmt.Printf("environment: %s\n", BuildEnvironment)

			bi, ok := debug.ReadBuildInfo()
			if ok && bi != nil {
				fmt.Printf("go version:  %s\n", bi.GoVersion)
			}
		},
	}

	return cmd
}
