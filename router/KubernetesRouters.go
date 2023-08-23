package router

import (
	"devopscenter/controller/kubernetes"
	"devopscenter/middleware"
	"github.com/gin-gonic/gin"
)

func KubernetesRegister(c *gin.Engine) {
	api := c.Group("/devops")
	{
		api.GET("/namespaces", middleware.JwtAuth(), kubernetes.NsList)

		api.GET("/deployment", middleware.JwtAuth(), kubernetes.DeployList)
		api.GET("/deployment/:name", middleware.JwtAuth(), kubernetes.DeployGet)
		api.GET("/deploymentV2", middleware.JwtAuth(), kubernetes.DeployListV3)
		api.PATCH("/deployment", middleware.JwtAuth(), kubernetes.DeployPatch)
		api.POST("/deployment", middleware.JwtAuth(), kubernetes.DeployAddV2)

		api.GET("/pod/:env/:namespace", kubernetes.PodList)

		api.GET("/ingress", middleware.JwtAuth(), kubernetes.IngressList)
		api.GET("/ingressV2", middleware.JwtAuth(), kubernetes.IngressListV2)
		//api.DELETE("/ingress", middleware.JwtAuth(), kubernetes.IngressDelete)
		api.POST("/ingress", middleware.JwtAuth(), kubernetes.IngressCreate)

		api.GET("/services", middleware.JwtAuth(), kubernetes.ServiceList)
		api.GET("/servicesV2", middleware.JwtAuth(), kubernetes.ServiceListV2)
		//api.DELETE("/services", middleware.JwtAuth(), kubernetes.ServiceDelete)
		api.POST("/services", middleware.JwtAuth(), kubernetes.ServiceCreateV2)

		api.GET("/cronjob", middleware.JwtAuth(), kubernetes.CronJobList)
		api.GET("/cronjobV2", middleware.JwtAuth(), kubernetes.CronJobListV2)
		api.POST("/cronjob", middleware.JwtAuth(), kubernetes.CronJobCreate)

		api.GET("/configmap", middleware.JwtAuth(), kubernetes.ConfigMapListV2)
		api.POST("/configmapV2", middleware.JwtAuth(), kubernetes.ConfigMapAddV2)
	}
}
