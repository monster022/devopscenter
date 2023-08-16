package model

type ConfigMap struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Env       string `json:"env"`
	//Data      struct {
	//	ASPNETCOREENVIRONMENT string `json:"ASPNETCORE_ENVIRONMENT"`
	//} `json:"data"`
}
type TTT struct {
	Env       string `json:"env"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Data      []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"data"`
}
