package service

import (
	"context"
	"devopscenter/helper"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ConfigMapList(configFile, namespace string) (*v1.ConfigMapList, error) {
	kubeEngine := helper.KubernetesConnect(configFile)
	data, err := kubeEngine.CoreV1().ConfigMaps(namespace).List(context.TODO(), metaV1.ListOptions{})
	return data, err
}

func ConfigMapAdd(configFile, namespace string, configMap *v1.ConfigMap) (*v1.ConfigMap, error) {
	kubeEngine := helper.KubernetesConnect(configFile)
	data, err := kubeEngine.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metaV1.CreateOptions{})
	return data, err
}
