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

var statusCMD = &cobra.Command{
	Use:   "status",
	Short: "Checks the reachability of each k8s cluster (from KUBECONFIG path or db)",
	Run:   statusFunc,
}

// init adds this command and his flags
func init() {
	rootCMD.AddCommand(statusCMD)
	statusCMD.PersistentFlags().Bool("db", false, "check the reachability of the kubeconfigs from the database")
}

func statusFunc(cmd *cobra.Command, args []string) {
	dbFlag, _ := cmd.Flags().GetBool("db")

	retcode := 0
	defer func() { os.Exit(retcode) }()

	if dbFlag { // check reachability of kubeconfigs stored in the database
		log.Debug().Str("command", "status").Str("database", dbPath).Send()

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

			itemColor := pterm.FgRed
			if kc.Reachable() {
				itemColor = pterm.FgGreen
			}

			list = append(list, pterm.BulletListItem{
				Level:       0,
				Text:        kc.Name,
				TextStyle:   pterm.NewStyle(itemColor),
				Bullet:      "⎈",
				BulletStyle: pterm.NewStyle(pterm.FgBlue),
			})
		}

		pterm.DefaultBulletList.WithItems(list).Render()

	} else { // check reachability of KUBECONFIG path files
		userHome, err := os.UserHomeDir()
		if err != nil {
			log.Fatal().Str("error", err.Error()).Msg("error getting the homeDir of the user")
		}
		kcRootDir = userHome + "/.kube/"
		// we dont have the filename arg, so we scan the $HOME/.kube/ directory
		files, err := os.ReadDir(kcRootDir)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}

		// check if we have any kubeconfig to list
		if len(files) == 0 {
			log.Info().Str("path", kcRootDir).Msg("Seems that you dont have any kubeconfig files ...")
			os.Exit(0)
		}

		// holds the list of kubeconfigs found
		var listItems []kubeconfig.Kubeconfig

		for _, f := range files {
			log.Debug().Str("file", f.Name()).Str("path", kcRootDir+f.Name()).Msg("loading kubeconfig...")
			// skip the default kubeconfig
			if f.Name() == "config" {
				log.Debug().Str("file", f.Name()).Msg("skipping default kubeconfig")
				continue
			}
			kc, err := kubeconfig.New(f.Name(), kcRootDir+f.Name())
			if err != nil {
				log.Debug().Str("file", f.Name()).Msg("not a valid kubeconfig")
				continue
			}
			listItems = append(listItems, *kc)
		}

		list := make([]pterm.BulletListItem, 0, len(listItems))
		for _, kc := range listItems {
			itemColor := pterm.FgRed
			if kc.Reachable() {
				itemColor = pterm.FgGreen
			}

			list = append(list, pterm.BulletListItem{
				Level:       0,
				Text:        kc.Name,
				TextStyle:   pterm.NewStyle(itemColor),
				Bullet:      "⎈",
				BulletStyle: pterm.NewStyle(pterm.FgBlue),
			})
		}

		pterm.DefaultBulletList.WithItems(list).Render()

	}

	log.Debug().Str("command", "status").Str("result", "successful").Send()
}
