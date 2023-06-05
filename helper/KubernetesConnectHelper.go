package helper

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func KubernetesConnect(configFile string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", configFile)
	if err != nil {
		panic(err.Error())
	}
	client, err1 := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err1.Error())
	}
	return client
}
