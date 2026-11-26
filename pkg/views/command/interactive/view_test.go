package interactive

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestVisibleSubcommandsSkipsHiddenAndHelperCommands(t *testing.T) {
	root := &cobra.Command{Use: "harbor"}
	root.AddCommand(&cobra.Command{Use: "artifact", Short: "Manage artifacts"})
	root.AddCommand(&cobra.Command{Use: "interactive", Short: "Browse interactively"})
	root.AddCommand(&cobra.Command{Use: "completion", Short: "Generate completion"})
	root.AddCommand(&cobra.Command{Use: "secret", Hidden: true})
	root.SetHelpCommand(&cobra.Command{Use: "help"})

	commands := visibleSubcommands(root)
	if len(commands) != 1 {
		t.Fatalf("expected 1 visible subcommand, got %d", len(commands))
	}

	if commands[0].Name() != "artifact" {
		t.Fatalf("expected artifact command, got %q", commands[0].Name())
	}
}

func TestRenderSelectionSummaryIncludesCommandPath(t *testing.T) {
	root := &cobra.Command{Use: "harbor"}
	artifact := &cobra.Command{
		Use:     "artifact",
		Short:   "Manage artifacts",
		Long:    "Manage artifacts in Harbor Repository",
		Example: "harbor artifact list",
	}
	root.AddCommand(artifact)

	summary := renderSelectionSummary(artifact)
	if !strings.Contains(summary, "harbor artifact") {
		t.Fatalf("expected command path in summary, got %q", summary)
	}

	if !strings.Contains(summary, "Manage artifacts") {
		t.Fatalf("expected summary text in output, got %q", summary)
	}
}
