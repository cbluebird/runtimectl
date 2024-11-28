package sdk

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func dynamicConfigInit(kubeConfig string) (*dynamic.DynamicClient, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(config)
}

func configInit(kubeConfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		panic(err.Error())
	}

	return kubernetes.NewForConfig(config)
}
