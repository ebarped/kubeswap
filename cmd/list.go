package cmd

import (
	"cmp"
	"fmt"
	"os"
	"slices"

	"github.com/ebarped/kubeswap/pkg/kubeconfig"
	"github.com/ebarped/kubeswap/pkg/kv"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var listCMD = &cobra.Command{
	Use:   "list",
	Short: "Lists all the kubeconfigs in the db",
	Run:   listFunc,
}

// init adds this command and his flags
func init() {
	rootCMD.AddCommand(listCMD)
}

func listFunc(cmd *cobra.Command, args []string) {
	retcode := 0
	defer func() { os.Exit(retcode) }()

	log.Debug().Str("command", "list").Str("database", dbPath).Send()

	db, err := kv.Open(dbPath)
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("error opening kv database")
		retcode = 1
		return
	}
	defer db.Close()

	items, err := db.GetAll()
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("error listing items from database")
		retcode = 1
		return
	}

	if len(items) == 0 {
		fmt.Println("Empty db...")
		retcode = 0
		return
	}

	slices.SortFunc(items, func(a, b kubeconfig.Kubeconfig) int {
		return cmp.Compare(a.Name, b.Name)
	})

	list := make([]pterm.BulletListItem, 0, len(items))
	for _, kc := range items {
		list = append(list, pterm.BulletListItem{
			Level:       0,
			Text:        kc.Name,
			Bullet:      "⎈",
			BulletStyle: pterm.NewStyle(pterm.FgBlue),
		})
	}
	pterm.DefaultBulletList.WithItems(list).Render()
	log.Debug().Str("command", "list").Str("result", "successful").Send()
}
