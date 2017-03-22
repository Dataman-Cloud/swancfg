package types

import (
	"time"
)

type App struct {
	ID               string    `json:"id,omitempty"`
	Name             string    `json:"name,omitempty"`
	Instances        int       `json:"instances,omitempty"`
	UpdatedInstances int       `json:"updatedInstances,omitempty"`
	RunningInstances int       `json:"runningInstances"`
	RunAs            string    `json:"runAs,omitempty"`
	Priority         int       `json:"priority"`
	ClusterId        string    `json:"clusterId,omitempty"`
	Created          time.Time `json:"created,omitempty"`
	Updated          time.Time `json:"updated,omitempty"`
	Mode             string    `json:"mode,omitempty"`
	State            string    `json:"state"`

	// use task for compatability now, should be slot here
	Tasks          []*Task `json:"tasks,omitempty"`
	CurrentVersion *Spec   `json:"currentVersion"`
}

// use task for compatability now, should be slot here
// and together with task history
type Task struct {
	ID        string `json:"id,omitempty"`
	AppId     string `json:"appId,omitempty"`
	VersionId string `json:"versionId,omitempty"`

	Status string `json:"status,omitempty"`

	OfferId       string `json:"offerId,omitempty"`
	AgentId       string `json:"agentId,omitempty"`
	AgentHostname string `json:"agentHostname,omitempty"`

	Cpu  float64 `json:"cpu,omitempty"`
	Mem  float64 `json:"mem,omitempty"`
	Disk float64 `json:"disk,omitempty"`

	History []*TaskHistory `json:"history,omitempty"`

	IP string `json:"ip,omitempty"`

	Created time.Time `json:"created,omitempty"`

	Image   string `json:"image,omitempty"`
	Healthy bool   `json:"healthy,omitempty"`
}

type TaskHistory struct {
	ID        string `json:"id,omitempty"`
	AppId     string `json:"appId,omitempty"`
	VersionId string `json:"versionId,omitempty"`

	OfferId       string `json:"offerId,omitempty"`
	AgentId       string `json:"agentId,omitempty"`
	AgentHostname string `json:"agentHostname,omitempty"`

	Cpu  float64 `json:"cpu,omitempty"`
	Mem  float64 `json:"mem,omitempty"`
	Disk float64 `json:"disk,omitempty"`

	State  string `json:"state,omitempty"`
	Reason string `json:"Reason,omitempty"`
	Stdout string `json:"stdout,omitempty"`
	Stderr string `json:"stderr,omitempty"`
}
