package command

import (
	"fmt"
	"strings"
	"testing"

	"github.com/YardRat0117/foxbox/internal/types"
)

type fakePanel struct {
	installedByName  map[string]string // key = tool.Name
	installedByImage map[string]string // key = tool.Image
}

func (f *fakePanel) InstallTool(toolName, version string) error {
	name := strings.SplitN(toolName, ":", 2)[0]
	f.installedByName[name] = version
	f.installedByImage[toolName] = version
	return nil
}

func (f *fakePanel) RemoveTool(toolName, imgName, version string) error {
	delete(f.installedByName, toolName)
	delete(f.installedByImage, imgName)
	return nil
}

func (f *fakePanel) RunTool(tool types.Tool, version string, args []string) error {
	fmt.Printf("running %s %v\n", tool.Image, args)
	return nil
}

func (f *fakePanel) CheckTools(tools map[string]types.Tool) (map[string]types.ToolStatus, error) {
	status := make(map[string]types.ToolStatus)
	for name := range tools {
		_, ok := f.installedByName[name]
		status[name] = types.ToolStatus{Installed: ok}
	}
	return status, nil
}

func (f *fakePanel) CleanTools(tools map[string]types.Tool) error {
	f.installedByName = map[string]string{}
	f.installedByImage = map[string]string{}
	return nil
}

func TestCLICommands(t *testing.T) {
	fake := &fakePanel{
		installedByName:  map[string]string{},
		installedByImage: map[string]string{},
	}
	panel = fake
	cfg = &types.Config{
		Tools: map[string]types.Tool{
			"mytool": {Image: "mytool:latest"},
		},
	}

	// install
	installCmd.Run(installCmd, []string{"mytool@1.0"})
	if fake.installedByName["mytool"] != "1.0" {
		t.Fatalf("install failed")
	}

	// run
	runCmd.Run(runCmd, []string{"mytool@1.0", "--help"})

	// list
	listCmd.Run(listCmd, []string{})

	// remove
	removeCmd.Run(removeCmd, []string{"mytool@1.0"})
	if _, ok := fake.installedByName["mytool"]; ok {
		t.Fatalf("remove failed")
	}

	// clean
	fake.installedByName["a"] = "1.0"
	fake.installedByName["b"] = "2.0"
	fake.installedByImage["a"] = "1.0"
	fake.installedByImage["b"] = "2.0"
	cleanCmd.Run(cleanCmd, []string{})
	if len(fake.installedByName) != 0 || len(fake.installedByImage) != 0 {
		t.Fatalf("clean failed")
	}
}
