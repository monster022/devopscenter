package router

import (
	"devopscenter/controller/kubernetes"
	"github.com/gin-gonic/gin"
)

func KubernetesRegister(c *gin.Engine) {
	api := c.Group("/devops")
	{
		api.GET("/namespaces", kubernetes.NsList)

		api.GET("/deployment", kubernetes.DeployList)
		api.GET("/deployment/:name", kubernetes.DeployGet)
		api.GET("/deploymentV2", kubernetes.DeployListV2)
		api.PATCH("/deployment", kubernetes.DeployPatch)

		api.GET("/ingress", kubernetes.IngressList)
		api.GET("/ingressV2", kubernetes.IngressListV2)
		api.DELETE("/ingress", kubernetes.IngressDelete)

		api.GET("/services", kubernetes.ServiceList)
		api.GET("/servicesV2", kubernetes.ServiceListV2)
		api.DELETE("/services", kubernetes.ServiceDelete)
		api.POST("/services", kubernetes.ServiceCreate)

		api.GET("/cronjob", kubernetes.CronJobList)
		api.GET("/cronjobV2", kubernetes.CronJobListV2)

	}
}
