package service

import (
	"context"
	"devopscenter/helper"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NamespaceList(configFile string) (*v1.NamespaceList, error) {
	kubeEngine := helper.KubernetesConnect(configFile)
	namespaceList, err := kubeEngine.CoreV1().Namespaces().List(context.TODO(), metaV1.ListOptions{})
	return namespaceList, err
}
