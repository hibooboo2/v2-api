package model

type Stack struct {
	Common

	ExternalId     string
	DockerCompose  string
	RancherCompose string
	StartOnCreate  bool
}
