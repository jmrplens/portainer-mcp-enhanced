package mcp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/portainer/portainer-mcp/pkg/portainer/models"
	"github.com/stretchr/testify/assert"
)

func TestHandleGetSSLSettings(t *testing.T) {
	tests := []struct {
		name          string
		sslSettings   models.SSLSettings
		mockError     error
		expectError   bool
		errorContains string
	}{
		{
			name: "successful SSL settings retrieval",
			sslSettings: models.SSLSettings{
				CertPath:    "/certs/cert.pem",
				KeyPath:     "/certs/key.pem",
				HTTPEnabled: true,
				SelfSigned:  false,
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:          "client error",
			sslSettings:   models.SSLSettings{},
			mockError:     assert.AnError,
			expectError:   true,
			errorContains: "failed to get SSL settings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockPortainerClient)
			mockClient.On("GetSSLSettings").Return(tt.sslSettings, tt.mockError)

			srv := &PortainerMCPServer{
				srv:   server.NewMCPServer("Test Server", "1.0.0"),
				cli:   mockClient,
				tools: make(map[string]mcp.Tool),
			}

			handler := srv.HandleGetSSLSettings()
			result, err := handler(context.Background(), mcp.CallToolRequest{})

			if tt.expectError {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.True(t, result.IsError)
				textContent, ok := result.Content[0].(mcp.TextContent)
				assert.True(t, ok)
				assert.Contains(t, textContent.Text, tt.errorContains)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				textContent, ok := result.Content[0].(mcp.TextContent)
				assert.True(t, ok)

				var settings models.SSLSettings
				err = json.Unmarshal([]byte(textContent.Text), &settings)
				assert.NoError(t, err)
				assert.Equal(t, tt.sslSettings, settings)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestHandleUpdateSSLSettings(t *testing.T) {
	httpEnabled := true

	tests := []struct {
		name          string
		request       mcp.CallToolRequest
		setupMock     func(*MockPortainerClient)
		expectError   bool
		errorContains string
	}{
		{
			name: "successful SSL settings update with all params",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"cert":        "-----BEGIN CERTIFICATE-----\nMIIB...",
						"key":         "-----BEGIN PRIVATE KEY-----\nMIIE...",
						"httpEnabled": true,
					},
				},
			},
			setupMock: func(m *MockPortainerClient) {
				m.On("UpdateSSLSettings", "-----BEGIN CERTIFICATE-----\nMIIB...", "-----BEGIN PRIVATE KEY-----\nMIIE...", &httpEnabled).Return(nil)
			},
			expectError: false,
		},
		{
			name: "successful SSL settings update with cert and key only",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"cert": "cert-data",
						"key":  "key-data",
					},
				},
			},
			setupMock: func(m *MockPortainerClient) {
				m.On("UpdateSSLSettings", "cert-data", "key-data", (*bool)(nil)).Return(nil)
			},
			expectError: false,
		},
		{
			name: "client error",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: map[string]any{
						"cert": "cert-data",
						"key":  "key-data",
					},
				},
			},
			setupMock: func(m *MockPortainerClient) {
				m.On("UpdateSSLSettings", "cert-data", "key-data", (*bool)(nil)).Return(assert.AnError)
			},
			expectError:   true,
			errorContains: "failed to update SSL settings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockPortainerClient)
			tt.setupMock(mockClient)

			srv := &PortainerMCPServer{
				srv:   server.NewMCPServer("Test Server", "1.0.0"),
				cli:   mockClient,
				tools: make(map[string]mcp.Tool),
			}

			handler := srv.HandleUpdateSSLSettings()
			result, err := handler(context.Background(), tt.request)

			if tt.expectError {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.True(t, result.IsError)
				textContent, ok := result.Content[0].(mcp.TextContent)
				assert.True(t, ok)
				assert.Contains(t, textContent.Text, tt.errorContains)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				textContent, ok := result.Content[0].(mcp.TextContent)
				assert.True(t, ok)
				assert.Contains(t, textContent.Text, "SSL settings updated successfully")
			}

			mockClient.AssertExpectations(t)
		})
	}
}
