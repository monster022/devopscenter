package model

type DockerRequestBody struct {
	Name   string   `json:"name"`
	Env    string   `json:"env"`
	Port   string   `json:"port"`
	Image  string   `json:"image"`
	Target []string `json:"target"`
}
