package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Namespace to use
var Namespace string

// Label selector to use
var Label string

// Field selector to use
var Field string

func init() {
	RootCmd.PersistentFlags().StringVarP(&Namespace, "namespace", "n", "default", "Namespace to use")
	RootCmd.PersistentFlags().StringVarP(&Label, "label", "l", "", "Label to use")
	RootCmd.PersistentFlags().StringVarP(&Field, "field", "f", "", "Field to use")
	RootCmd.AddCommand(GreenPlanetCmd)
}

var RootCmd = &cobra.Command{
	Use: "kwatch",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Namespace: %s\n", Namespace)
		fmt.Printf("Label %s:)\n", Label)
		fmt.Printf("Field %s:)\n", Field)
	},
}
