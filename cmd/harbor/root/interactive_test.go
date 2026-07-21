package root

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestCommandArgsFromPath(t *testing.T) {
	tests := []struct {
		path string
		want []string
	}{
		{path: "harbor", want: nil},
		{path: "harbor info", want: []string{"info"}},
		{path: "harbor artifact list", want: []string{"artifact", "list"}},
		{path: " harbor   project   view ", want: []string{"project", "view"}},
	}

	for _, tt := range tests {
		got := commandArgsFromPath(tt.path)
		if !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("commandArgsFromPath(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestChangedRootFlags(t *testing.T) {
	root := &cobra.Command{Use: "harbor"}
	root.PersistentFlags().String("config", "", "config file")
	root.PersistentFlags().Bool("verbose", false, "verbose output")

	if err := root.PersistentFlags().Set("config", "/tmp/config.yaml"); err != nil {
		t.Fatalf("failed to set config flag: %v", err)
	}
	if err := root.PersistentFlags().Set("verbose", "true"); err != nil {
		t.Fatalf("failed to set verbose flag: %v", err)
	}

	got := changedRootFlags(root)
	want := []string{"--config", "/tmp/config.yaml", "--verbose"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("changedRootFlags() = %v, want %v", got, want)
	}
}
