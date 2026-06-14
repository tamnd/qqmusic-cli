package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tamnd/qqmusic-cli/qqmusic"
)

// chartsCmd returns a command that lists all available chart names.
func (a *App) chartsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "charts",
		Short: "List available QQ Music chart names",
		RunE: func(cmd *cobra.Command, _ []string) error {
			type chartRow struct {
				Name  string `json:"name"`
				Label string `json:"label"`
			}
			// Fixed order so output is deterministic.
			rows := []chartRow{
				{Name: "hot", Label: qqmusic.ChartLabels["hot"]},
				{Name: "new", Label: qqmusic.ChartLabels["new"]},
				{Name: "rising", Label: qqmusic.ChartLabels["rising"]},
			}
			if a.output == string(FormatTable) || a.output == "auto" {
				for _, r := range rows {
					_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%-10s  %s\n", r.Name, r.Label)
				}
				return nil
			}
			return a.render(rows)
		},
	}
}
