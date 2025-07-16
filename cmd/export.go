package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ebarped/kubeswap/pkg/kv"
	"github.com/spf13/cobra"
)

var exportCMD = &cobra.Command{
	Use:   "export",
	Short: "Exports all kubeconfigs from the db into the desired folder",
	Run:   exportFunc,
}

// init adds this command and his flags
func init() {
	rootCMD.AddCommand(exportCMD)
	exportCMD.Flags().StringP("path", "p", defaultKubeconfigPath(), "path to export the kubeconfigs")
}

func exportFunc(cmd *cobra.Command, args []string) {
	path, _ := cmd.Flags().GetString("path")

	retcode := 0
	defer func() { os.Exit(retcode) }()

	log.Debug().Str("command", "export").Str("database", dbPath).Msg("")

	db, err := kv.Open(dbPath)
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("error opening kv database")
		retcode = 1
		return
	}
	defer db.Close()

	kubeconfigs, err := db.GetAll()
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("error getting kubeconfigs from database")
		retcode = 1
		return
	}

	for _, kc := range kubeconfigs {
		// obtain path to export this kubeconfig
		kcDstPath := path
		// if user provides a path, maybe we have to append the '/'
		if !strings.Contains(kcDstPath, "/") {
			kcDstPath = kcDstPath + "/"
		}
		kcDstPath = kcDstPath + kc.Name

		log.Debug().Str("command", "export").Str("database", dbPath).Msgf("Creating kubeconfig '%s' on '%s'", kc.Name, kcDstPath)

		// create file to write
		f, err := os.Create(kcDstPath)
		if err != nil {
			log.Error().Str("command", "export").Str("error", err.Error()).Msgf("error creating kubeconfig file %s", kcDstPath)
			retcode = 1
			return
		}
		defer f.Close()

		// write the info
		_, err = f.WriteString(kc.Content)
		if err != nil {
			log.Error().Str("command", "export").Str("error", err.Error()).Msgf("error writing to kubeconfig file %s", kcDstPath)
			retcode = 1
			return
		}

		f.Sync()
	}

	fmt.Println("Export successful! 🚀")
}
