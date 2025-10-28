package command

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

)

// func newListCommand creates a `cobra.Command` object `list` with given context
func newListCommand(rootCtx *rootContext) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all configured tools status",
		RunE: func(_ *cobra.Command, _ []string) error {
			// Assign panel and config
			panel, cfg := rootCtx.panel, rootCtx.cfg

			// Call panel to list installed tools
			installed, err := panel.CheckTools(cfg.Tools)
			if err != nil {
				return fmt.Errorf("Error listing installed tools: %e", err)
			}

			// Const width
			const nameWidth, parenWidth = 15, 25

			// Print to console with formatting
			fmt.Println("Configured tools:")

			for name, tool := range cfg.Tools {
				// Parse
				s, exists := installed[name]

				// Determine status
				status := "[not installed]"
				if exists && s.Installed {
					status = "[Installed]"
				}

				// Format status
				status = fmt.Sprintf("%-15s", status)

				// Determine and format tags
				tags := ""
				if exists && len(s.LocalTags) > 0 {
					tags = "tags " + strings.Join(s.LocalTags, ", ")
				}

				// Print current tool
				fmt.Printf("- %-*s %-*s %s %s\n",
					nameWidth, name,
					parenWidth, fmt.Sprintf("(%s)", tool.Image),
					status, tags)
			}

			return nil
		},
	}
}
