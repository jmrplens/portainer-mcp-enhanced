package client

import (
	"errors"
	"testing"

	apimodels "github.com/portainer/client-api-go/v2/pkg/models"
	"github.com/portainer/portainer-mcp/pkg/portainer/models"
	"github.com/stretchr/testify/assert"
)

func TestGetMOTD(t *testing.T) {
	tests := []struct {
		name          string
		mockMOTD      *apimodels.MotdMotdResponse
		mockError     error
		expected      models.MOTD
		expectedError bool
	}{
		{
			name: "successful retrieval",
			mockMOTD: &apimodels.MotdMotdResponse{
				Title:   "Welcome",
				Message: "Hello World",
				Style:   "info",
				Hash:    []int64{1, 2, 3},
				ContentLayout: map[string]string{
					"key": "value",
				},
			},
			expected: models.MOTD{
				Title:   "Welcome",
				Message: "Hello World",
				Style:   "info",
				Hash:    []int64{1, 2, 3},
				ContentLayout: map[string]string{
					"key": "value",
				},
			},
		},
		{
			name:          "api error",
			mockError:     errors.New("api error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := new(MockPortainerAPI)
			mockAPI.On("GetMOTD").Return(tt.mockMOTD, tt.mockError)

			client := &PortainerClient{cli: mockAPI}
			motd, err := client.GetMOTD()

			if tt.expectedError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, motd)
			mockAPI.AssertExpectations(t)
		})
	}
}
