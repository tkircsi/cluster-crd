package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/tkircsi/cluster-crd/api/types/v1alpha1"
	clientV1alpha1 "github.com/tkircsi/cluster-crd/clientset/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
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

	store := WatchResources(clientSet)

	for {
		clustersFromStore := store.List()
		for _, cl := range clustersFromStore {
			fmt.Printf("clusters in store: %s\n", cl.(*v1alpha1.Cluster).Spec.ClusterName)
		}
		time.Sleep(1 * time.Second)

	}

}

func WatchResources(clientSet clientV1alpha1.ExtensionV1Alpha1Interface) cache.Store {
	clusterStore, clusterController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return clientSet.Clusters("").List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return clientSet.Clusters("").Watch(options)
			},
		},
		&v1alpha1.Cluster{},
		30*time.Second,
		cache.ResourceEventHandlerFuncs{},
	)

	go clusterController.Run(wait.NeverStop)
	return clusterStore
}
