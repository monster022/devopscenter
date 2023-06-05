package service

import (
	"context"
	"devopscenter/helper"
	"k8s.io/api/extensions/v1beta1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func IngressList(configFile, namespace string) (*v1beta1.IngressList, error) {
	kubeEngine := helper.KubernetesConnect(configFile)
	// IngressList, err := kubeEngine.NetworkingV1beta1().Ingresses(namespace).List(context.TODO(), metaV1.ListOptions{})
	IngressList, err := kubeEngine.ExtensionsV1beta1().Ingresses(namespace).List(context.TODO(), metaV1.ListOptions{})
	return IngressList, err
}

func IngressDelete(configFile, namespace, ingressName string) error {
	kubeEngine := helper.KubernetesConnect(configFile)
	err := kubeEngine.ExtensionsV1beta1().Ingresses(namespace).Delete(context.TODO(), ingressName, metaV1.DeleteOptions{})
	return err
}
