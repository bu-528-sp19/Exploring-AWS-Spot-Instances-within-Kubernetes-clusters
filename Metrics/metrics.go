package main

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"

	"os"
)

type k8s struct {
	clientset kubernetes.Interface
}

func newK8s() (*kubernetes.Clientset, error) {
	//path := os.Getenv()
	//absPath, _ := filepath.Abs("Aditya\.kube")
	config, err := clientcmd.BuildConfigFromFlags("", "C:/Users/Aditya/.kube/config")
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func main() {
	k8s, err := newK8s()
	if err != nil {
		fmt.Println(err)
		return
	}
	node, err := k8s.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
		}
	fmt.Printf("There are %d nodes in the cluster\n", len(node.Items))

	fmt.Println("-----------------------------------------------------")
	config, err := clientcmd.BuildConfigFromFlags("", "C:/Users/Aditya/.kube/config")

	mc, err := metrics.NewForConfig(config)


	list, err := k8s.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listing nodes: %v", err)
	}

	for  _ , nodes  := range list.Items  {
		fmt.Printf("Node: %s\n", nodes.Name)
		nodes, err := k8s.CoreV1().Nodes().Get(nodes.Name, metav1.GetOptions{})
		nodemetric, err := mc.MetricsV1beta1().NodeMetricses().Get(nodes.Name, metav1.GetOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting node: %v", err)
		}
		fmt.Println("Allocatable memory:",nodes.Status.Allocatable.Memory().Value())
		fmt.Println("Memory:",nodemetric.Usage.Memory().Value())
		fmt.Println("Percent memory used:",(float64(nodemetric.Usage.Memory().Value())/float64(nodes.Status.Allocatable.Memory().Value())*100))
		fmt.Println("CPU Usage:", nodemetric.Usage.Cpu().String())
		fmt.Println("-----------------------------------------------------")

	}

	if err != nil {
		panic(err)
	}


}
