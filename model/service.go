package model

type Service struct {
	Common

	ContainerIds       []Id                `json:"containerIds"`
	ContainerSelector  string              `json:"containerSelector"`
	ContainerTemplates []ContainerTemplate `json:"containerTemplates"`
	CreateIndex        string              `json:"createIndex"`
	Fqdn               string              `json:"fqdn"`
	HealthState        string              `json:"healthState"`
	HostnameOverride   string              `json:"hostnameOverride"`
	LinkSelector       string              `json:"linkSelector"`
	LinkedServiceIds   []Id                `json:"linkedServiceIds"`
	Metadata           string              `json:"metadata"`
	PublicEndpoints    []PublicEndpoint    `json:"publicEndpoints"`
	PullImage          string              `json:"pullImage"`
	RequestedIpAddress string              `json:"requestedIpAddress"`
	RetainIpAddress    bool                `json:"retainIpAddress"`
	Scale              int                 `json:"scale"`
	ServiceIpAddress   string              `json:"serviceIpAddress"`
	StackId            Id                  `json:"stackId"`
	StartOnce          bool                `json:"startOnce"`
}
