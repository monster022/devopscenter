package model

import "time"

type HarborTagResponse struct {
	Digest        string    `json:"digest"`
	Name          string    `json:"name"`
	Size          int       `json:"size"`
	Architecture  string    `json:"architecture"`
	Os            string    `json:"os"`
	OsVersion     string    `json:"os.version"`
	DockerVersion string    `json:"docker_version"`
	Author        string    `json:"author"`
	Created       time.Time `json:"created"`
	Config        struct {
		Labels interface{} `json:"labels"`
	} `json:"config"`
	Immutable bool          `json:"immutable"`
	Signature interface{}   `json:"signature"`
	Labels    []interface{} `json:"labels"`
	PushTime  time.Time     `json:"push_time"`
	PullTime  time.Time     `json:"pull_time"`
}
