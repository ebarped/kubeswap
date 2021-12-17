package cmd

import (
	"github.com/ebarped/kubeswap/pkg/kv"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var listCMD = &cobra.Command{
	Use:   "list",
	Short: "lists all the kubeconfigs in the db",
	Run:   listFunc,
}

// init adds this command and his flags
func init() {
	rootCMD.AddCommand(listCMD)
}

func listFunc(cmd *cobra.Command, args []string) {
	log.Debug().Str("command", "list").Str("database", dbPath).Send()

	db, err := kv.Open(dbPath)
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("error opening kv database")
	}
	defer db.CloseDB()

	var list []pterm.BulletListItem

	items, err := db.GetAll()
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("error listing items from database")
	}
	for _, kc := range items {
		list = append(list, pterm.BulletListItem{
			Level:  0,
			Text:   kc.Name,
			Bullet: "⎈",
		})
	}
	pterm.DefaultBulletList.WithItems(list).Render()
}
