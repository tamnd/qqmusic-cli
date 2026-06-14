package cli

import (
	"github.com/spf13/cobra"
)

// topCmd returns the top command for browsing chart songs.
func (a *App) topCmd() *cobra.Command {
	var chart string
	cmd := &cobra.Command{
		Use:   "top",
		Short: "List QQ Music chart songs",
		Long: `List songs from a QQ Music chart.

Available charts: hot (热歌榜), new (新歌榜), rising (飙升榜).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			n := a.effectiveLimit(30)
			a.progressf("fetching %d songs from %q chart...", n, chart)
			songs, err := a.client.Top(cmd.Context(), chart, n)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.renderOrEmpty(songs, len(songs))
		},
	}
	cmd.Flags().StringVar(&chart, "chart", "hot", "chart to browse: hot|new|rising")
	return cmd
}
