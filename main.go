package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/tkircsi/cluster-crd/api/types/v1alpha1"
	clientV1alpha1 "github.com/tkircsi/cluster-crd/clientset/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
	flag.Parse()
}

func main() {
	var config *rest.Config
	var err error

	if kubeconfig == "" {
		log.Printf("using in-cluster configuration\n")
		config, err = rest.InClusterConfig()
	} else {
		log.Printf("using configuration from %q\n", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		log.Fatalln(err)
	}

	v1alpha1.AddToScheme(scheme.Scheme)

	clientSet, err := clientV1alpha1.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}

	clusters, err := clientSet.Clusters("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("clusters found: %+v\n", clusters)

}
