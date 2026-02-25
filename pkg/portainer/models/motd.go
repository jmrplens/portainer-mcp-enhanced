package models

import (
	apimodels "github.com/portainer/client-api-go/v2/pkg/models"
)

// MOTD represents the message of the day
type MOTD struct {
	Title         string            `json:"title"`
	Message       string            `json:"message"`
	Style         string            `json:"style"`
	Hash          []int64           `json:"hash"`
	ContentLayout map[string]string `json:"contentLayout,omitempty"`
}

// ConvertToMOTD converts a raw MotdMotdResponse to a local MOTD
func ConvertToMOTD(raw *apimodels.MotdMotdResponse) MOTD {
	return MOTD{
		Title:         raw.Title,
		Message:       raw.Message,
		Style:         raw.Style,
		Hash:          raw.Hash,
		ContentLayout: raw.ContentLayout,
	}
}
