package command

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

)

func newListCommand(rootCtx *rootContext) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all configured tools status",
		RunE: func(_ *cobra.Command, _ []string) error {
			panel, cfg := rootCtx.panel, rootCtx.cfg

			installed, err := panel.CheckTools(cfg.Tools)
			if err != nil {
				return fmt.Errorf("Error listing installed tools: %e", err)
			}

			const nameWidth, parenWidth = 15, 25
			fmt.Println("Configured tools:")
			for name, tool := range cfg.Tools {
				s, exists := installed[name]
				status := "[not installed]"
				if exists && s.Installed {
					status = "[Installed]"
				}
				status = fmt.Sprintf("%-15s", status)

				tags := ""
				if exists && len(s.LocalTags) > 0 {
					tags = "tags " + strings.Join(s.LocalTags, ", ")
				}

				fmt.Printf("- %-*s %-*s %s %s\n",
					nameWidth, name,
					parenWidth, fmt.Sprintf("(%s)", tool.Image),
					status, tags)
			}

			return nil
		},
	}
}
