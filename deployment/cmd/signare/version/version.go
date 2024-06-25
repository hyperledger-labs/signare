// Package version provides a command utility to retrieve the signare version through the binary.
package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

type BuildInfo struct {
	CommitHash string
	BuildTime  string
	Tag        string
	Branch     string
}

var buildInfo BuildInfo

func Command(currentBuildInfo BuildInfo) *cobra.Command {
	buildInfo = currentBuildInfo
	cmd := &cobra.Command{
		Use:  "version",
		Long: "prints version",
		RunE: executeVersionCmd,
	}
	return cmd
}

func executeVersionCmd(_ *cobra.Command, _ []string) error {
	if buildInfo.Tag != "" {
		fmt.Print("version:", buildInfo.Tag)
	} else {
		fmt.Print("branch:", buildInfo.Branch)
	}
	fmt.Println(" commit-hash:", buildInfo.CommitHash, " build-time:", buildInfo.BuildTime)
	return nil
}
