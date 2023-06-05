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

func CronJobDelete(configFile, jobName, namespace string) error {
	kubeEngine := helper.KubernetesConnect(configFile)
	err := kubeEngine.BatchV1beta1().CronJobs(namespace).Delete(context.TODO(), jobName, metaV1.DeleteOptions{})
	return err
}
