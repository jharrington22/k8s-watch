package cmd

import (
	"fmt"
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	k8sclient *kubernetes.Clientset
	options   *RuntimeOptions
)

// RuntimeOptions for kubernetes configuration
type RuntimeOptions struct {
	KubeconfigPath string
}

// Init Kubernetes config
func Init(opt *RuntimeOptions) error {
	fmt.Println()
	fmt.Println("Using kubeconfig: ", opt.KubeconfigPath)

	fmt.Printf("Namespace: %s\n", Namespace)
	fmt.Printf("Label: %s\n", Label)
	fmt.Printf("Field: %s\n", Field)

	config, err := clientcmd.BuildConfigFromFlags("", opt.KubeconfigPath)
	if err != nil {
		panic(err.Error())

	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	k8sclient = clientset
	options = opt
	return nil
}
