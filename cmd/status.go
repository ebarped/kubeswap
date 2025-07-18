package cmd

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"sync"

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

	// holds the list of kubeconfigs found (either in KUBECONFIG path or db)
	var items []kubeconfig.Kubeconfig

	if dbFlag { // check reachability of kubeconfigs stored in the database
		log.Debug().Str("command", "status").Str("database", dbPath).Send()

		db, err := kv.Open(dbPath)
		if err != nil {
			log.Error().Str("command", "status").Str("error", err.Error()).Msg("error opening kv database")
			retcode = 1
			return
		}
		defer db.Close()

		items, err = db.GetAll()
		if err != nil {
			log.Error().Str("command", "status").Str("error", err.Error()).Msg("error listing items from database")
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

	} else { // check reachability of KUBECONFIG path files
		// we dont have the filename arg, so we scan the $HOME/.kube/ directory
		files, err := os.ReadDir(defaultKubeconfigPath())
		if err != nil {
			log.Fatal().Str("command", "status").Msg(err.Error())
		}

		// check if we have any kubeconfig to list
		if len(files) == 0 {
			log.Info().Str("command", "status").Str("path", kcRootDir).Msg("Seems that you dont have any kubeconfig files ...")
			os.Exit(0)
		}

		for _, f := range files {
			log.Debug().Str("command", "status").Str("file", f.Name()).Str("path", kcRootDir+f.Name()).Msg("loading kubeconfig...")
			// skip the default kubeconfig
			if f.Name() == "config" {
				log.Debug().Str("command", "status").Str("file", f.Name()).Msg("skipping default kubeconfig")
				continue
			}
			kc, err := kubeconfig.New(f.Name(), kcRootDir+f.Name())
			if err != nil {
				log.Debug().Str("command", "status").Str("file", f.Name()).Msg("not a valid kubeconfig")
				continue
			}
			items = append(items, *kc)
		}

	}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{} // protects the list
	wg.Add(len(items))

	list := make([]pterm.BulletListItem, 0, len(items))
	for _, kc := range items {
		go func(kc kubeconfig.Kubeconfig) {
			defer wg.Done()

			itemColor := pterm.FgRed
			if kc.Reachable() {
				itemColor = pterm.FgGreen
			}

			item := pterm.BulletListItem{
				Level:       0,
				Text:        kc.Name,
				TextStyle:   pterm.NewStyle(itemColor),
				Bullet:      "⎈",
				BulletStyle: pterm.NewStyle(pterm.FgBlue),
			}

			mu.Lock()
			list = append(list, item)
			mu.Unlock()
		}(kc)
	}

	wg.Wait()

	pterm.DefaultBulletList.WithItems(list).Render()

	log.Debug().Str("command", "status").Str("result", "successful").Send()
}
