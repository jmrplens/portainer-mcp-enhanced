package client

import (
	"github.com/portainer/portainer-mcp/pkg/portainer/models"
)

// GetRoles retrieves all roles from the Portainer server.
func (c *PortainerClient) GetRoles() ([]models.Role, error) {
	rawRoles, err := c.cli.ListRoles()
	if err != nil {
		return nil, err
	}

	roles := make([]models.Role, len(rawRoles))
	for i, raw := range rawRoles {
		roles[i] = models.ConvertToRole(raw)
	}

	return roles, nil
}
