package main

import "time"

type TableStruct struct {
	Header []string
	Data   [][]string
}

type StatusResponse struct {
	ID     string
	Status string
}

type SubmittedFiles struct {
	Options string
}

type CallItem struct {
	ExecutionStatus     string
	Stdout              string
	Stderr              string
	Attempt             int
	ShardIndex          int
	Start               time.Time
	End                 time.Time
	Labels              Label
	MonitoringLog       string
	CommandLine         string
	DockerImageUsed     string
	SubWorkflowID       string
	SubWorkflowMetadata MetadataResponse
	RuntimeAttributes   RuntimeAttributes
	CallCaching         CallCachingData
	ExecutionEvents     []ExecutionEvents
	BackendLogs         BackendLogs
	Failures            []string
}

type BackendLogs struct {
	Log string
}

type ExecutionEvents struct {
	StartTime   time.Time
	Description string
	EndTime     time.Time
}

type CallCachingData struct {
	Result string
	Hit    bool
}

type RuntimeAttributes struct {
	BootDiskSizeGb string
	CPU            string
	Disks          string
	Docker         string
	Memory         string
	Preemptible    string
}

type Label struct {
	CromwellWorkflowID string `json:"cromwell-workflow-id"`
	WdlTaskName        string `json:"wdl-task-name"`
}

type Failure struct {
	CausedBy []Failure
	Message  string
}

type MetadataResponse struct {
	WorkflowRoot   string
	WorkflowName   string
	SubmittedFiles SubmittedFiles
	RootWorkflowID string
	Calls          map[string][]CallItem
	Inputs         map[string]interface{}
	Outputs        map[string]interface{}
	Start          time.Time
	End            time.Time
	Status         string
	Failures       []Failure
}

type MetadataTableResponse struct {
	Metadata MetadataResponse
}

type ErrorResponse struct {
	HTTPStatus string
}
