package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jharrington22/k8s-watch/cmd"
	"k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("Using kubeconfig: ", kubeconfig)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	api := clientset.BatchV1()

	listOptions := metav1.ListOptions{
		LabelSelector: cmd.Label,
		FieldSelector: cmd.Field,
	}

	jobs, err := api.Jobs(cmd.Namespace).List(listOptions)
	if err != nil {
		log.Fatal(err)
	}

	printJobs(jobs)
	fmt.Println()

	watcher, err := api.Jobs(cmd.Namespace).Watch(listOptions)
	if err != nil {
		log.Fatal(err)
	}

	ch := watcher.ResultChan()

	for event := range ch {
		job, ok := event.Object.(*v1.Job)
		if !ok {
			log.Fatal("unexpected type")
		}
		switch event.Type {
		case watch.Added:
			fmt.Printf("Job added\n")
			fmt.Printf(job.ObjectMeta.Name + "\n")
			fmt.Printf(job.Spec.Template.ObjectMeta.Labels["build"] + "\n")
			printJobs(jobs)
			fmt.Printf("added")
		case watch.Deleted:
			fmt.Printf("Job deleted\n")
		}
	}
}

func printJobs(jobs *v1.JobList) {
	if len(jobs.Items) == 0 {
		log.Println("No jobs found")
		return
	}
	template := "%-32s%-8s\n"
	fmt.Printf(template, "NAME", "STATUS")
	for _, job := range jobs.Items {
		fmt.Printf(template, job.Name, string(job.Status.Active))
	}
}
