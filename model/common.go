package model

import (
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/rancher/go-rancher/client"
)

type Id string

type Common struct {
	client.Resource

	Name        string
	Description string
	State       string
	Uuid        string
	Kind        string
	Type        string
	Data        types.JSONText
	Created     time.Time
	Removed     time.Time

	Transitioning        string
	TransitioningMessage string
}
