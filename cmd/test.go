package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ebarped/kubeswap/pkg/kubeconfig"
	"github.com/spf13/cobra"
)

var testCMD = &cobra.Command{
	Use:   "test",
	Short: "test",
	Run:   testFunc,
}

// init adds this command and his flags
func init() {
	rootCMD.AddCommand(testCMD)
}

func testFunc(cmd *cobra.Command, args []string) {
	retcode := 0
	defer func() { os.Exit(retcode) }()

	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("error getting the homeDir of the user")
	}
	kcRootDir := userHome + "/.kube/"

	files, err := ioutil.ReadDir(kcRootDir)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	for _, f := range files {
		log.Debug().Msgf("loading kubeconfig %s from %s", f.Name(), kcRootDir+f.Name())
		kc, err := kubeconfig.New(f.Name(), kcRootDir+f.Name())
		if err != nil {
			log.Error().Msg(err.Error())
			continue
		}
		fmt.Println(kc.Name)
		//fmt.Println(kc.Path)
	}
}
