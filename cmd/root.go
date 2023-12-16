package cmd

import (
	"log"

	"github.com/fibonachyy/sternx/cmd/sternx"
	"github.com/spf13/cobra"
)

func Execute() {
	root := &cobra.Command{
		Use:     "Stern-X",
		Short:   "Stern-X hiering task",
		Version: "0.1",
	}

	sternx.Register(root)

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
