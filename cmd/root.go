package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	batchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
)

// Namespace to use
var Namespace string

// Label selector to use
var Label string

// Field selector to use
var Field string

// Verbose flag
var Verbose bool

// Api Object
var Api batchv1.BatchV1Interface

// RootCmd main cobra entry point
var RootCmd = &cobra.Command{
	Use: "kwatch",
	Run: func(cmd *cobra.Command, args []string) {
		err := Init(opt)
		if err != nil {
			log.Fatal(err)
		}
		if Verbose {
			fmt.Println("Runnig kwatch..")
		}
	},
}

var opt = &RuntimeOptions{}

// Default kube configuration location
var defaultKubeConfig = filepath.Join(
	os.Getenv("HOME"), ".kube", "config",
)

func init() {
	RootCmd.AddCommand(JobCmd)
	fmt.Printf("Configuration location: %s\n", defaultKubeConfig)
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().StringVarP(&opt.KubeconfigPath, "kubeconfig", "k", defaultKubeConfig, "The path to a kube config file on the local filesystem")
	RootCmd.PersistentFlags().StringVarP(&Namespace, "namespace", "n", "default", "Namespace to use")
	RootCmd.PersistentFlags().StringVarP(&Label, "label", "l", "", "Select Kubernetes resources based labels key/value pairs")
	RootCmd.PersistentFlags().StringVarP(&Field, "field", "f", "", "Select Kubernetes resources based on the value of on or more fields")
}
