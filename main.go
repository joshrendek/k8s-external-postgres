package main

import (
	"flag"
	"time"

	"github.com/golang/glog"
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"os"

	"path/filepath"

	"github.com/joshrendek/k8s-external-postgres/pkg/apis/postgresql/v1"
	clientset "github.com/joshrendek/k8s-external-postgres/pkg/client/clientset/versioned"
	informers "github.com/joshrendek/k8s-external-postgres/pkg/client/informers/externalversions"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"k8s.io/sample-controller/pkg/signals"
)

var (
	masterURL   string
	kubeconfig  string
	postgresURL string
	isConsole   bool
)

func main() {
	flag.Parse()

	if isConsole {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	if home := homeDir(); home != "" && kubeconfig == "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)

	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	exampleClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %s", err.Error())
	}

	crdConfig, _ := GetClientConfig(kubeconfig)
	crdClient, err := apiextcs.NewForConfig(crdConfig)

	v1.CreateCRD(crdClient)

	exampleInformerFactory := informers.NewSharedInformerFactory(exampleClient, time.Second*1)

	controller := NewController(kubeClient, exampleClient, exampleInformerFactory)

	go exampleInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}

func GetClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&postgresURL, "postgres-uri", "postgres://localhost/template1?sslmode=disable", "URI to connect to postgres")
	flag.BoolVar(&isConsole, "console", false, "whether to console log or json log")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
