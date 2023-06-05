package service

import (
	"context"
	"devopscenter/helper"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func SvcList(configFile, namespace string) (*v1.ServiceList, error) {
	kubeEngine := helper.KubernetesConnect(configFile)
	serviceList, err := kubeEngine.CoreV1().Services(namespace).List(context.TODO(), metaV1.ListOptions{})
	return serviceList, err
}

func SvcDelete(configFile, namespace, serviceName string) error {
	kubeEngine := helper.KubernetesConnect(configFile)
	err := kubeEngine.CoreV1().Services(namespace).Delete(context.TODO(), serviceName, metaV1.DeleteOptions{})
	return err
}

func SvcCreate(configFile, namespace string, service *v1.Service) (*v1.Service, error) {
	kubeEngine := helper.KubernetesConnect(configFile)
	serviceList, err := kubeEngine.CoreV1().Services(namespace).Create(context.TODO(), service, metaV1.CreateOptions{})
	return serviceList, err
}
