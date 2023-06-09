package service

import (
	"context"
	"devopscenter/helper"
	"fmt"
	v1 "k8s.io/api/apps/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

func DeploymentList(configFile, namespace string) (*v1.DeploymentList, error) {
	kubeEngine := helper.KubernetesConnect(configFile)
	kubeEngine.CoreV1().Pods(namespace).List(context.TODO(), metaV1.ListOptions{})
	deployment, err := kubeEngine.AppsV1().Deployments(namespace).List(context.TODO(), metaV1.ListOptions{})
	return deployment, err
}
func DeploymentGet(configFile, namespace, deploymentName string) (*v1.Deployment, error) {
	kubeEngine := helper.KubernetesConnect(configFile)
	deployment, err := kubeEngine.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metaV1.GetOptions{})
	return deployment, err
}

func DeploymentDelete(configFile, deploymentName, namespace string) error {
	kubeEngine := helper.KubernetesConnect(configFile)
	err := kubeEngine.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metaV1.DeleteOptions{})
	return err
}

func DeploymentImagePatch(configFile, namespace, containerName, deploymentName, image string) error {
	kubeEngine := helper.KubernetesConnect(configFile)
	data := []byte(fmt.Sprintf(`{"spec": {"template": {"spec": {"containers": [{"name": "%s", "image": "%s"}]}}}}`, containerName, image))
	_, err := kubeEngine.AppsV1().Deployments(namespace).Patch(context.TODO(), deploymentName, types.StrategicMergePatchType, data, metaV1.PatchOptions{})
	return err
}

func DeploymentAdd(configFile, namespace string, deployment *v1.Deployment) (*v1.Deployment, error) {
	kubeEngine := helper.KubernetesConnect(configFile)
	data, err := kubeEngine.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metaV1.CreateOptions{})
	return data, err
}

func PodList(configFile, namespace string) (watch.Interface, error) {
	kubeEngine := helper.KubernetesConnect(configFile)
	podList, err := kubeEngine.CoreV1().Pods(namespace).Watch(context.TODO(), metaV1.ListOptions{})
	return podList, err
}
