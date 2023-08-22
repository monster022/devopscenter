package service

import (
	"context"
	"devopscenter/helper"
	"k8s.io/api/batch/v1beta1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CronJobList(configFile, namespace string) (*v1beta1.CronJobList, error) {
	kubeEngine := helper.KubernetesConnect(configFile)
	cronjob, err := kubeEngine.BatchV1beta1().CronJobs(namespace).List(context.TODO(), metaV1.ListOptions{})
	return cronjob, err
}

func CronJobCreate(configFile, namespace string, cronjob *v1beta1.CronJob) (*v1beta1.CronJob, error) {
	kubeEngine := helper.KubernetesConnect(configFile)
	result, err := kubeEngine.BatchV1beta1().CronJobs(namespace).Create(context.TODO(), cronjob, metaV1.CreateOptions{})
	return result, err
}
