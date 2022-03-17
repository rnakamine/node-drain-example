package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/kubectl/pkg/drain"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// get node
	nodeName := "node-xxx"
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	drainer := &drain.Helper{
		Ctx:                 context.TODO(),
		Client:              clientset,
		Force:               true,
		IgnoreAllDaemonSets: true,

		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	err = drain.RunCordonOrUncordon(drainer, node, true)
	if err != nil {
		panic(err.Error())
	}
	err = drain.RunNodeDrain(drainer, nodeName)
	if err != nil {
		panic(err.Error())
	}
}
