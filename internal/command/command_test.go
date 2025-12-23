//go:build legacy

package command

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/YardRat0117/foxbox/internal/types"
)

// ------ Fake Panel ------

// type runCall logs each fall
type runCall struct {
	tool    types.Tool
	version string
	args    []string
}

// type fakePanel is simulating type Panel to run mock test.
type fakePanel struct {
	installed   map[string]string // key = tool.Name
	lastRunCall *runCall
}

// func InstallTool simulates the installation procedure.
func (f *fakePanel) InstallTool(toolName, version string) error {
	name := strings.SplitN(toolName, ":", 2)[0]
	f.installed[name] = version
	return nil
}

// func RemoveTool simulates the removing procedure.
func (f *fakePanel) RemoveTool(toolName, imgName, version string) error {
	delete(f.installed, toolName)
	return nil
}

// func RunTool simulates the running procedure.
func (f *fakePanel) RunTool(tool types.Tool, version string, args []string) error {
	f.lastRunCall = &runCall{tool, version, args}
	fmt.Printf("running %s %v\n", tool.Image, args)
	return nil
}

// func CheckTools simulates the tool installation checking procedure.
func (f *fakePanel) CheckTools(tools map[string]types.Tool) (map[string]types.ToolStatus, error) {
	status := make(map[string]types.ToolStatus)
	for name := range tools {
		_, ok := f.installed[name]
		status[name] = types.ToolStatus{Installed: ok}
	}
	return status, nil
}

// func CleanTools simulates the cleaning procedure.
func (f *fakePanel) CleanTools(tools map[string]types.Tool) error {
	for name := range tools {
		delete(f.installed, name)
	}
	return nil
}

// For Error Injection
type fakePanelErr struct{}

func (f *fakePanelErr) InstallTool(toolName, version string) error {
	return fmt.Errorf("mock install error")
}
func (f *fakePanelErr) RemoveTool(toolName, imgName, version string) error {
	return fmt.Errorf("mock remove error")
}
func (f *fakePanelErr) RunTool(tool types.Tool, version string, args []string) error {
	return fmt.Errorf("mock run error")
}
func (f *fakePanelErr) CheckTools(tools map[string]types.Tool) (map[string]types.ToolStatus, error) {
	return nil, fmt.Errorf("mock check error")
}
func (f *fakePanelErr) CleanTools(tools map[string]types.Tool) error {
	return fmt.Errorf("mock clean error")
}

// ------ Test ------
func TestPackageCommand(t *testing.T) {
	// Initialization
	// Initialize fakeRootCtx
	fakeP := &fakePanel{
		installed: make(map[string]string),
	}
	fakeC := &types.Config{
		Tools: map[string]types.Tool{
			"mytool": {Image: "mytool:latest"},
			"rat":    {Image: "rat:1.17"},
		},
	}
	fakeRootCtx := &rootContext{
		panel: fakeP,
		cfg:   fakeC,
	}

	// Init commands
	installCmd := newInstallCommand(fakeRootCtx)
	runCmd := newRunCommand(fakeRootCtx)
	removeCmd := newRemoveCommand(fakeRootCtx)
	cleanCmd := newCleanCommand(fakeRootCtx)
	listCmd := newListCommand(fakeRootCtx)
	versionCmd := newVersionCommand()

	// Test install
	t.Run("install success", func(t *testing.T) {
		err := installCmd.RunE(installCmd, []string{"mytool@latest"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fakeP.installed["mytool"] != "latest" {
			t.Fatalf("expected mytool@latest installed")
		}
	})

	// Test run
	t.Run("run success", func(t *testing.T) {
		err := runCmd.RunE(runCmd, []string{"mytool@latest", "--help"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := fakeP.lastRunCall
		if got == nil || got.tool.Image != "mytool:latest" {
			t.Fatalf("expected RunTool called with mytool:latest")
		}
		if !reflect.DeepEqual(got.args, []string{"--help"}) {
			t.Fatalf("args mismatch: %v", got.args)
		}
	})

	// Test List
	t.Run("test list for coverage", func(t *testing.T) {
		listCmd.RunE(listCmd, []string{})
	})

	// Test remove
	t.Run("remove success", func(t *testing.T) {
		err := removeCmd.RunE(removeCmd, []string{"mytool@latest"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := fakeP.installed["mytool"]; ok {
			t.Fatalf("expected mytool removed")
		}
	})

	// Test clean
	t.Run("clean only own tools", func(t *testing.T) {
		fakeP.installed["mytool"] = "latest"
		fakeP.installed["othertool"] = "latest"
		err := cleanCmd.RunE(cleanCmd, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := fakeP.installed["mytool"]; ok {
			t.Fatalf("failed to clean mytool")
		}
		if _, ok := fakeP.installed["othertool"]; !ok {
			t.Fatalf("should not clean othertool")
		}
	})

	// Test Version
	t.Run("test version for coverage", func(t *testing.T) {
		versionCmd.Run(versionCmd, []string{})
	})

	// Error injection
	badPanel := &fakePanelErr{}
	badRootCtx := &rootContext{
		panel: badPanel,
		cfg:   fakeC,
	}

	// Init commands
	badInstallCmd := newInstallCommand(badRootCtx)
	badRunCmd := newRunCommand(badRootCtx)
	badRemoveCmd := newRemoveCommand(badRootCtx)
	badCleanCmd := newCleanCommand(badRootCtx)

	t.Run("install error injected", func(t *testing.T) {
		err := badInstallCmd.RunE(badInstallCmd, []string{"mytool@latest"})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "mock install error") {
			t.Fatalf("unexpected error message: %v", err)
		}
	})

	t.Run("run error injected", func(t *testing.T) {
		err := badRunCmd.RunE(badRunCmd, []string{"mytool@latest", "--test"})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "mock run error") {
			t.Fatalf("unexpected error message: %v", err)
		}
	})

	t.Run("remove error injected", func(t *testing.T) {
		err := badRemoveCmd.RunE(badRemoveCmd, []string{"mytool@latest"})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "mock remove error") {
			t.Fatalf("unexpected error message: %v", err)
		}
	})

	t.Run("clean error injected", func(t *testing.T) {
		err := badCleanCmd.RunE(installCmd, []string{"mytool@latest"})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "mock clean error") {
			t.Fatalf("unexpected error message: %v", err)
		}
	})
}
