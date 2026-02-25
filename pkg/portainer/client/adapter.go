package client

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	sdkclient "github.com/portainer/client-api-go/v2/client"
	swaggerclient "github.com/portainer/client-api-go/v2/pkg/client"
	"github.com/portainer/client-api-go/v2/pkg/client/backup"
	"github.com/portainer/client-api-go/v2/pkg/client/custom_templates"
	"github.com/portainer/client-api-go/v2/pkg/client/endpoints"
	"github.com/portainer/client-api-go/v2/pkg/client/motd"
	"github.com/portainer/client-api-go/v2/pkg/client/registries"
	"github.com/portainer/client-api-go/v2/pkg/client/roles"
	"github.com/portainer/client-api-go/v2/pkg/client/tags"
	"github.com/portainer/client-api-go/v2/pkg/client/teams"
	"github.com/portainer/client-api-go/v2/pkg/client/users"
	"github.com/portainer/client-api-go/v2/pkg/client/webhooks"
	apimodels "github.com/portainer/client-api-go/v2/pkg/models"
)

// portainerAPIAdapter wraps the SDK PortainerClient and adds methods
// that are available in the Swagger-generated client but not exposed
// by the SDK's high-level client (e.g., delete operations).
type portainerAPIAdapter struct {
	*sdkclient.PortainerClient
	swagger *swaggerclient.PortainerClientAPI
}

// newPortainerAPIAdapter creates a new adapter that embeds the SDK high-level
// client and also holds a reference to the low-level Swagger client for
// operations not exposed by the SDK.
func newPortainerAPIAdapter(host, apiKey string, skipTLSVerify bool) *portainerAPIAdapter {
	sdkCli := sdkclient.NewPortainerClient(host, apiKey, sdkclient.WithSkipTLSVerify(skipTLSVerify))

	transport := httptransport.New(host, "/api", []string{"https"})
	if skipTLSVerify {
		transport.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	apiKeyAuth := runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		return r.SetHeaderParam("x-api-key", apiKey)
	})
	transport.DefaultAuthentication = apiKeyAuth

	return &portainerAPIAdapter{
		PortainerClient: sdkCli,
		swagger:         swaggerclient.New(transport, nil),
	}
}

// DeleteTag deletes a tag by ID using the low-level Swagger client.
func (a *portainerAPIAdapter) DeleteTag(id int64) error {
	params := tags.NewTagDeleteParams().WithID(id)
	_, err := a.swagger.Tags.TagDelete(params, nil)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return nil
}

// DeleteTeam deletes a team by ID using the low-level Swagger client.
func (a *portainerAPIAdapter) DeleteTeam(id int64) error {
	params := teams.NewTeamDeleteParams().WithID(id)
	_, err := a.swagger.Teams.TeamDelete(params, nil)
	if err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}
	return nil
}

// DeleteUser deletes a user by ID using the low-level Swagger client.
func (a *portainerAPIAdapter) DeleteUser(id int64) error {
	params := users.NewUserDeleteParams().WithID(id)
	_, err := a.swagger.Users.UserDelete(params, nil)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// DeleteEndpoint deletes an endpoint by ID using the low-level Swagger client.
func (a *portainerAPIAdapter) DeleteEndpoint(id int64) error {
	params := endpoints.NewEndpointDeleteParams().WithID(id)
	_, err := a.swagger.Endpoints.EndpointDelete(params, nil)
	if err != nil {
		return fmt.Errorf("failed to delete endpoint: %w", err)
	}
	return nil
}

// SnapshotEndpoint triggers a snapshot for a single endpoint.
func (a *portainerAPIAdapter) SnapshotEndpoint(id int64) error {
	params := endpoints.NewEndpointSnapshotParams().WithID(id)
	_, err := a.swagger.Endpoints.EndpointSnapshot(params, nil)
	if err != nil {
		return fmt.Errorf("failed to snapshot endpoint: %w", err)
	}
	return nil
}

// SnapshotAllEndpoints triggers a snapshot for all endpoints.
func (a *portainerAPIAdapter) SnapshotAllEndpoints() error {
	params := endpoints.NewEndpointSnapshotsParams()
	_, err := a.swagger.Endpoints.EndpointSnapshots(params, nil)
	if err != nil {
		return fmt.Errorf("failed to snapshot all endpoints: %w", err)
	}
	return nil
}

// ListWebhooks retrieves all webhooks using the low-level Swagger client.
func (a *portainerAPIAdapter) ListWebhooks() ([]*apimodels.PortainerWebhook, error) {
	params := webhooks.NewGetWebhooksParams()
	resp, err := a.swagger.Webhooks.GetWebhooks(params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list webhooks: %w", err)
	}
	return resp.Payload, nil
}

// CreateWebhook creates a new webhook using the low-level Swagger client.
func (a *portainerAPIAdapter) CreateWebhook(resourceId string, endpointId int64, webhookType int64) (int64, error) {
	payload := &apimodels.WebhooksWebhookCreatePayload{
		ResourceID:  resourceId,
		EndpointID:  endpointId,
		WebhookType: webhookType,
	}
	params := webhooks.NewPostWebhooksParams().WithBody(payload)
	resp, err := a.swagger.Webhooks.PostWebhooks(params, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create webhook: %w", err)
	}
	return resp.Payload.ID, nil
}

// DeleteWebhook deletes a webhook by ID using the low-level Swagger client.
func (a *portainerAPIAdapter) DeleteWebhook(id int64) error {
	params := webhooks.NewDeleteWebhooksIDParams().WithID(id)
	_, err := a.swagger.Webhooks.DeleteWebhooksID(params, nil)
	if err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}
	return nil
}

// ListCustomTemplates lists all custom templates.
func (a *portainerAPIAdapter) ListCustomTemplates() ([]*apimodels.PortainereeCustomTemplate, error) {
	params := custom_templates.NewCustomTemplateListParams()
	resp, err := a.swagger.CustomTemplates.CustomTemplateList(params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list custom templates: %w", err)
	}
	return resp.Payload, nil
}

// GetCustomTemplate retrieves a custom template by ID.
func (a *portainerAPIAdapter) GetCustomTemplate(id int64) (*apimodels.PortainereeCustomTemplate, error) {
	params := custom_templates.NewCustomTemplateInspectParams().WithID(id)
	resp, err := a.swagger.CustomTemplates.CustomTemplateInspect(params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get custom template: %w", err)
	}
	return resp.Payload, nil
}

// GetCustomTemplateFile retrieves the file content of a custom template.
func (a *portainerAPIAdapter) GetCustomTemplateFile(id int64) (string, error) {
	params := custom_templates.NewCustomTemplateFileParams().WithID(id)
	resp, err := a.swagger.CustomTemplates.CustomTemplateFile(params, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get custom template file: %w", err)
	}
	return resp.Payload.FileContent, nil
}

// CreateCustomTemplate creates a new custom template from file content.
func (a *portainerAPIAdapter) CreateCustomTemplate(payload *apimodels.CustomtemplatesCustomTemplateFromFileContentPayload) (*apimodels.PortainereeCustomTemplate, error) {
	params := custom_templates.NewCustomTemplateCreateStringParams().WithBody(payload)
	resp, err := a.swagger.CustomTemplates.CustomTemplateCreateString(params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create custom template: %w", err)
	}
	return resp.Payload, nil
}

// DeleteCustomTemplate deletes a custom template by ID.
func (a *portainerAPIAdapter) DeleteCustomTemplate(id int64) error {
	params := custom_templates.NewCustomTemplateDeleteParams().WithID(id)
	_, err := a.swagger.CustomTemplates.CustomTemplateDelete(params, nil)
	if err != nil {
		return fmt.Errorf("failed to delete custom template: %w", err)
	}
	return nil
}

// ListRegistries lists all registries.
func (a *portainerAPIAdapter) ListRegistries() ([]*apimodels.PortainereeRegistry, error) {
	params := registries.NewRegistryListParams()
	resp, err := a.swagger.Registries.RegistryList(params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list registries: %w", err)
	}
	return resp.Payload, nil
}

// GetRegistryByID retrieves a registry by ID.
func (a *portainerAPIAdapter) GetRegistryByID(id int64) (*apimodels.PortainereeRegistry, error) {
	params := registries.NewRegistryInspectParams().WithID(id)
	resp, err := a.swagger.Registries.RegistryInspect(params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get registry: %w", err)
	}
	return resp.Payload, nil
}

// CreateRegistry creates a new registry.
func (a *portainerAPIAdapter) CreateRegistry(body *apimodels.RegistriesRegistryCreatePayload) (int64, error) {
	params := registries.NewRegistryCreateParams().WithBody(body)
	resp, err := a.swagger.Registries.RegistryCreate(params, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create registry: %w", err)
	}
	return resp.Payload.ID, nil
}

// UpdateRegistry updates an existing registry.
func (a *portainerAPIAdapter) UpdateRegistry(id int64, body *apimodels.RegistriesRegistryUpdatePayload) error {
	params := registries.NewRegistryUpdateParams().WithID(id).WithBody(body)
	_, err := a.swagger.Registries.RegistryUpdate(params, nil)
	if err != nil {
		return fmt.Errorf("failed to update registry: %w", err)
	}
	return nil
}

// DeleteRegistry deletes a registry by ID.
func (a *portainerAPIAdapter) DeleteRegistry(id int64) error {
	params := registries.NewRegistryDeleteParams().WithID(id)
	_, err := a.swagger.Registries.RegistryDelete(params, nil)
	if err != nil {
		return fmt.Errorf("failed to delete registry: %w", err)
	}
	return nil
}

// GetBackupStatus retrieves the status of the last backup.
func (a *portainerAPIAdapter) GetBackupStatus() (*apimodels.BackupBackupStatus, error) {
	params := backup.NewBackupStatusFetchParams()
	resp, err := a.swagger.Backup.BackupStatusFetch(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get backup status: %w", err)
	}
	return resp.Payload, nil
}

// GetBackupSettings retrieves the S3 backup settings.
func (a *portainerAPIAdapter) GetBackupSettings() (*apimodels.PortainereeS3BackupSettings, error) {
	params := backup.NewBackupSettingsFetchParams()
	resp, err := a.swagger.Backup.BackupSettingsFetch(params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get backup settings: %w", err)
	}
	return resp.Payload, nil
}

// CreateBackup triggers a backup with an optional password.
func (a *portainerAPIAdapter) CreateBackup(password string) error {
	body := &apimodels.BackupBackupPayload{
		Password: password,
	}
	params := backup.NewBackupParams().WithBody(body)
	_, err := a.swagger.Backup.Backup(params, nil)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	return nil
}

// BackupToS3 triggers a backup to S3.
func (a *portainerAPIAdapter) BackupToS3(body *apimodels.BackupS3BackupPayload) error {
	params := backup.NewBackupToS3Params().WithBody(body)
	_, err := a.swagger.Backup.BackupToS3(params, nil)
	if err != nil {
		return fmt.Errorf("failed to backup to S3: %w", err)
	}
	return nil
}

// RestoreFromS3 triggers a restore from S3.
func (a *portainerAPIAdapter) RestoreFromS3(body *apimodels.BackupRestoreS3Settings) error {
	params := backup.NewRestoreFromS3Params().WithBody(body)
	_, err := a.swagger.Backup.RestoreFromS3(params)
	if err != nil {
		return fmt.Errorf("failed to restore from S3: %w", err)
	}
	return nil
}

// ListRoles lists all roles.
func (a *portainerAPIAdapter) ListRoles() ([]*apimodels.PortainereeRole, error) {
	params := roles.NewRoleListParams()
	resp, err := a.swagger.Roles.RoleList(params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	return resp.Payload, nil
}

// GetMOTD retrieves the message of the day.
func (a *portainerAPIAdapter) GetMOTD() (*apimodels.MotdMotdResponse, error) {
	params := motd.NewMOTDParams()
	resp, err := a.swagger.Motd.MOTD(params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get MOTD: %w", err)
	}
	return resp.Payload, nil
}
