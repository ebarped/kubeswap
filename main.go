package main

import (
	"github.com/ebarped/kubeswap/cmd"
)

func main() {
	cmd.Execute()
	cmd.CloseDB()
}
