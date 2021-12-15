package cmd

import (
	"fmt"

	"github.com/ebarped/kubeswap/pkg/kv"
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
		log.Fatal().Str("error", err.Error()).Msg("error opening kv database")
	}
	defer db.CloseDB()
	items, err := db.GetAll()
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("error listing items from database")
	}
	fmt.Printf("List of kubeconfigs:\n")
	for _, kc := range items {
		fmt.Printf("Name: %s\n", kc.Name)
		fmt.Printf("Content:\n%s\n", kc.Content)
	}
}
