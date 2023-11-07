package cmd

import (
	"fmt"
	"strings"
)

// NukeParameters struct
type NukeParameters struct {
	ConfigPath string

	Targets      []string
	Excludes     []string
	CloudControl []string

	NoDryRun   bool
	Force      bool
	ForceSleep int
	Quiet      bool

	MaxWaitRetries int
}

// Validate nuke params
func (p *NukeParameters) Validate() error {
	if strings.TrimSpace(p.ConfigPath) == "" {
		return fmt.Errorf("you have to specify the --config flag")
	}

	return nil
}
