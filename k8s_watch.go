package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var ns, label, field string
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)

	flag.StringVar(&ns, "namespace", "", "namespace")
	flag.StringVar(&label, "l", "", "Label Selector")
	flag.StringVar(&field, "f", "", "Field Selector")
	flag.Parse()

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
		LabelSelector: label,
		FieldSelector: field,
	}

	jobs, err := api.Jobs(ns).List(listOptions)
	if err != nil {
		log.Fatal(err)
	}

	printJobs(jobs)
	fmt.Println()

	watcher, err := api.Jobs(ns).Watch(listOptions)
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
