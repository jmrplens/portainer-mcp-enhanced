package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// metaAction maps an action name to its handler and access metadata.
type metaAction struct {
	name     string
	handler  func(s *PortainerMCPServer) server.ToolHandlerFunc
	readOnly bool   // true = always available; false = hidden in read-only mode
	toolName string // the granular tool name (schema.go ToolXxx constant) whose
	// tools.yaml parameter schema gets merged into the parent meta-tool's
	// inputSchema — see registerOneMetaTool. Every action must set this so its
	// real parameter types (not just the bare "action" enum) reach MCP clients.
}

// metaToolDef describes a single grouped meta-tool.
type metaToolDef struct {
	name        string
	description string
	actions     []metaAction
	annotation  mcp.ToolAnnotation
}

// boolPtr is a convenience helper for creating *bool values.
func boolPtr(v bool) *bool { return &v }

// metaToolDefinitions returns the complete list of meta-tool groups.
// Each group aggregates several existing granular tools behind a single
// "action" enum parameter. Read-only mode filters write actions at
// registration time, so the enum only exposes permitted actions.
func metaToolDefinitions() []metaToolDef {
	return []metaToolDef{
		{
			name:        "manage_environments",
			description: "Manage Portainer environments, environment groups, and tags. Actions: list_environments, get_environment, delete_environment, snapshot_environment, snapshot_all_environments, update_environment_tags, update_environment_user_accesses, update_environment_team_accesses, list_environment_groups, create_environment_group, update_environment_group_name, update_environment_group_environments, update_environment_group_tags, list_environment_tags, create_environment_tag, delete_environment_tag. Set 'action' parameter to choose.",
			actions: []metaAction{
				{name: "list_environments", handler: (*PortainerMCPServer).HandleGetEnvironments, readOnly: true, toolName: ToolListEnvironments},
				{name: "get_environment", handler: (*PortainerMCPServer).HandleGetEnvironment, readOnly: true, toolName: ToolGetEnvironment},
				{name: "delete_environment", handler: (*PortainerMCPServer).HandleDeleteEnvironment, readOnly: false, toolName: ToolDeleteEnvironment},
				{name: "snapshot_environment", handler: (*PortainerMCPServer).HandleSnapshotEnvironment, readOnly: false, toolName: ToolSnapshotEnvironment},
				{name: "snapshot_all_environments", handler: (*PortainerMCPServer).HandleSnapshotAllEnvironments, readOnly: false, toolName: ToolSnapshotAllEnvironments},
				{name: "update_environment_tags", handler: (*PortainerMCPServer).HandleUpdateEnvironmentTags, readOnly: false, toolName: ToolUpdateEnvironmentTags},
				{name: "update_environment_user_accesses", handler: (*PortainerMCPServer).HandleUpdateEnvironmentUserAccesses, readOnly: false, toolName: ToolUpdateEnvironmentUserAccesses},
				{name: "update_environment_team_accesses", handler: (*PortainerMCPServer).HandleUpdateEnvironmentTeamAccesses, readOnly: false, toolName: ToolUpdateEnvironmentTeamAccesses},
				{name: "list_environment_groups", handler: (*PortainerMCPServer).HandleGetEnvironmentGroups, readOnly: true, toolName: ToolListEnvironmentGroups},
				{name: "create_environment_group", handler: (*PortainerMCPServer).HandleCreateEnvironmentGroup, readOnly: false, toolName: ToolCreateEnvironmentGroup},
				{name: "update_environment_group_name", handler: (*PortainerMCPServer).HandleUpdateEnvironmentGroupName, readOnly: false, toolName: ToolUpdateEnvironmentGroupName},
				{name: "update_environment_group_environments", handler: (*PortainerMCPServer).HandleUpdateEnvironmentGroupEnvironments, readOnly: false, toolName: ToolUpdateEnvironmentGroupEnvironments},
				{name: "update_environment_group_tags", handler: (*PortainerMCPServer).HandleUpdateEnvironmentGroupTags, readOnly: false, toolName: ToolUpdateEnvironmentGroupTags},
				{name: "list_environment_tags", handler: (*PortainerMCPServer).HandleGetEnvironmentTags, readOnly: true, toolName: ToolListEnvironmentTags},
				{name: "create_environment_tag", handler: (*PortainerMCPServer).HandleCreateEnvironmentTag, readOnly: false, toolName: ToolCreateEnvironmentTag},
				{name: "delete_environment_tag", handler: (*PortainerMCPServer).HandleDeleteEnvironmentTag, readOnly: false, toolName: ToolDeleteEnvironmentTag},
			},
			annotation: mcp.ToolAnnotation{
				Title:           "Manage Environments",
				ReadOnlyHint:    boolPtr(false),
				DestructiveHint: boolPtr(true),
				IdempotentHint:  boolPtr(false),
				OpenWorldHint:   boolPtr(false),
			},
		},
		{
			name:        "manage_stacks",
			description: "Manage Docker stacks (Compose and Edge deployments). Actions: list_stacks, list_regular_stacks, get_stack, get_stack_file, inspect_stack_file, create_stack, create_regular_stack, update_stack, delete_stack, update_stack_git, redeploy_stack_git, start_stack, stop_stack, migrate_stack. 'get_stack_file'/'create_stack' operate on edge stacks; 'inspect_stack_file'/'create_regular_stack' operate on regular (non-edge) stacks deployed directly to one environment. Set 'action' parameter to choose.",
			actions: []metaAction{
				{name: "list_stacks", handler: (*PortainerMCPServer).HandleGetStacks, readOnly: true, toolName: ToolListStacks},
				{name: "list_regular_stacks", handler: (*PortainerMCPServer).HandleListRegularStacks, readOnly: true, toolName: ToolListRegularStacks},
				{name: "get_stack", handler: (*PortainerMCPServer).HandleInspectStack, readOnly: true, toolName: ToolGetStack},
				{name: "get_stack_file", handler: (*PortainerMCPServer).HandleGetStackFile, readOnly: true, toolName: ToolGetStackFile},
				{name: "inspect_stack_file", handler: (*PortainerMCPServer).HandleInspectStackFile, readOnly: true, toolName: ToolInspectStackFile},
				{name: "create_stack", handler: (*PortainerMCPServer).HandleCreateStack, readOnly: false, toolName: ToolCreateStack},
				{name: "create_regular_stack", handler: (*PortainerMCPServer).HandleCreateRegularStack, readOnly: false, toolName: ToolCreateRegularStack},
				{name: "update_stack", handler: (*PortainerMCPServer).HandleUpdateStack, readOnly: false, toolName: ToolUpdateStack},
				{name: "delete_stack", handler: (*PortainerMCPServer).HandleDeleteStack, readOnly: false, toolName: ToolDeleteStack},
				{name: "update_stack_git", handler: (*PortainerMCPServer).HandleUpdateStackGit, readOnly: false, toolName: ToolUpdateStackGit},
				{name: "redeploy_stack_git", handler: (*PortainerMCPServer).HandleRedeployStackGit, readOnly: false, toolName: ToolRedeployStackGit},
				{name: "start_stack", handler: (*PortainerMCPServer).HandleStartStack, readOnly: false, toolName: ToolStartStack},
				{name: "stop_stack", handler: (*PortainerMCPServer).HandleStopStack, readOnly: false, toolName: ToolStopStack},
				{name: "migrate_stack", handler: (*PortainerMCPServer).HandleMigrateStack, readOnly: false, toolName: ToolMigrateStack},
			},
			annotation: mcp.ToolAnnotation{
				Title:           "Manage Stacks",
				ReadOnlyHint:    boolPtr(false),
				DestructiveHint: boolPtr(true),
				IdempotentHint:  boolPtr(false),
				OpenWorldHint:   boolPtr(false),
			},
		},
		{
			name:        "manage_access_groups",
			description: "Manage access groups for environment-level permissions. Actions: list_access_groups, create_access_group, update_access_group_name, update_access_group_user_accesses, update_access_group_team_accesses, add_environment_to_access_group, remove_environment_from_access_group. Set 'action' parameter to choose.",
			actions: []metaAction{
				{name: "list_access_groups", handler: (*PortainerMCPServer).HandleGetAccessGroups, readOnly: true, toolName: ToolListAccessGroups},
				{name: "create_access_group", handler: (*PortainerMCPServer).HandleCreateAccessGroup, readOnly: false, toolName: ToolCreateAccessGroup},
				{name: "update_access_group_name", handler: (*PortainerMCPServer).HandleUpdateAccessGroupName, readOnly: false, toolName: ToolUpdateAccessGroupName},
				{name: "update_access_group_user_accesses", handler: (*PortainerMCPServer).HandleUpdateAccessGroupUserAccesses, readOnly: false, toolName: ToolUpdateAccessGroupUserAccesses},
				{name: "update_access_group_team_accesses", handler: (*PortainerMCPServer).HandleUpdateAccessGroupTeamAccesses, readOnly: false, toolName: ToolUpdateAccessGroupTeamAccesses},
				{name: "add_environment_to_access_group", handler: (*PortainerMCPServer).HandleAddEnvironmentToAccessGroup, readOnly: false, toolName: ToolAddEnvironmentToAccessGroup},
				{name: "remove_environment_from_access_group", handler: (*PortainerMCPServer).HandleRemoveEnvironmentFromAccessGroup, readOnly: false, toolName: ToolRemoveEnvironmentFromAccessGroup},
			},
			annotation: mcp.ToolAnnotation{
				Title:           "Manage Access Groups",
				ReadOnlyHint:    boolPtr(false),
				DestructiveHint: boolPtr(true),
				IdempotentHint:  boolPtr(false),
				OpenWorldHint:   boolPtr(false),
			},
		},
		{
			name:        "manage_users",
			description: "Manage Portainer user accounts and roles. Actions: list_users, get_user, create_user, delete_user, update_user_role. Set 'action' parameter to choose.",
			actions: []metaAction{
				{name: "list_users", handler: (*PortainerMCPServer).HandleGetUsers, readOnly: true, toolName: ToolListUsers},
				{name: "get_user", handler: (*PortainerMCPServer).HandleGetUser, readOnly: true, toolName: ToolGetUser},
				{name: "create_user", handler: (*PortainerMCPServer).HandleCreateUser, readOnly: false, toolName: ToolCreateUser},
				{name: "delete_user", handler: (*PortainerMCPServer).HandleDeleteUser, readOnly: false, toolName: ToolDeleteUser},
				{name: "update_user_role", handler: (*PortainerMCPServer).HandleUpdateUserRole, readOnly: false, toolName: ToolUpdateUserRole},
			},
			annotation: mcp.ToolAnnotation{
				Title:           "Manage Users",
				ReadOnlyHint:    boolPtr(false),
				DestructiveHint: boolPtr(true),
				IdempotentHint:  boolPtr(false),
				OpenWorldHint:   boolPtr(false),
			},
		},
		{
			name:        "manage_teams",
			description: "Manage Portainer teams and membership. Actions: list_teams, get_team, create_team, delete_team, update_team_name, update_team_members. Set 'action' parameter to choose.",
			actions: []metaAction{
				{name: "list_teams", handler: (*PortainerMCPServer).HandleGetTeams, readOnly: true, toolName: ToolListTeams},
				{name: "get_team", handler: (*PortainerMCPServer).HandleGetTeam, readOnly: true, toolName: ToolGetTeam},
				{name: "create_team", handler: (*PortainerMCPServer).HandleCreateTeam, readOnly: false, toolName: ToolCreateTeam},
				{name: "delete_team", handler: (*PortainerMCPServer).HandleDeleteTeam, readOnly: false, toolName: ToolDeleteTeam},
				{name: "update_team_name", handler: (*PortainerMCPServer).HandleUpdateTeamName, readOnly: false, toolName: ToolUpdateTeamName},
				{name: "update_team_members", handler: (*PortainerMCPServer).HandleUpdateTeamMembers, readOnly: false, toolName: ToolUpdateTeamMembers},
			},
			annotation: mcp.ToolAnnotation{
				Title:           "Manage Teams",
				ReadOnlyHint:    boolPtr(false),
				DestructiveHint: boolPtr(true),
				IdempotentHint:  boolPtr(false),
				OpenWorldHint:   boolPtr(false),
			},
		},
		{
			name:        "manage_docker",
			description: "Interact with Docker environments via dashboards and proxy API calls. Actions: get_docker_dashboard, docker_proxy. Set 'action' parameter to choose.",
			actions: []metaAction{
				{name: "get_docker_dashboard", handler: (*PortainerMCPServer).HandleGetDockerDashboard, readOnly: true, toolName: ToolGetDockerDashboard},
				{name: "docker_proxy", handler: (*PortainerMCPServer).HandleDockerProxy, readOnly: false, toolName: ToolDockerProxy},
			},
			annotation: mcp.ToolAnnotation{
				Title:           "Manage Docker",
				ReadOnlyHint:    boolPtr(false),
				DestructiveHint: boolPtr(true),
				IdempotentHint:  boolPtr(false),
				OpenWorldHint:   boolPtr(true),
			},
		},
		{
			name:        "manage_kubernetes",
			description: "Interact with Kubernetes environments via dashboards, namespaces, kubeconfig, and proxy API calls. Actions: get_kubernetes_resource_stripped, get_kubernetes_dashboard, list_kubernetes_namespaces, get_kubernetes_config, kubernetes_proxy. Set 'action' parameter to choose.",
			actions: []metaAction{
				{name: "get_kubernetes_resource_stripped", handler: (*PortainerMCPServer).HandleKubernetesProxyStripped, readOnly: true, toolName: ToolKubernetesProxyStripped},
				{name: "get_kubernetes_dashboard", handler: (*PortainerMCPServer).HandleGetKubernetesDashboard, readOnly: true, toolName: ToolGetKubernetesDashboard},
				{name: "list_kubernetes_namespaces", handler: (*PortainerMCPServer).HandleListKubernetesNamespaces, readOnly: true, toolName: ToolListKubernetesNamespaces},
				{name: "get_kubernetes_config", handler: (*PortainerMCPServer).HandleGetKubernetesConfig, readOnly: true, toolName: ToolGetKubernetesConfig},
				{name: "kubernetes_proxy", handler: (*PortainerMCPServer).HandleKubernetesProxy, readOnly: false, toolName: ToolKubernetesProxy},
			},
			annotation: mcp.ToolAnnotation{
				Title:           "Manage Kubernetes",
				ReadOnlyHint:    boolPtr(false),
				DestructiveHint: boolPtr(true),
				IdempotentHint:  boolPtr(false),
				OpenWorldHint:   boolPtr(true),
			},
		},
		{
			name:        "manage_helm",
			description: "Manage Helm repositories, charts, and releases on Kubernetes environments. Actions: list_helm_repositories, search_helm_charts, list_helm_releases, get_helm_release_history, add_helm_repository, remove_helm_repository, install_helm_chart, delete_helm_release. Set 'action' parameter to choose.",
			actions: []metaAction{
				{name: "list_helm_repositories", handler: (*PortainerMCPServer).HandleListHelmRepositories, readOnly: true, toolName: ToolListHelmRepositories},
				{name: "search_helm_charts", handler: (*PortainerMCPServer).HandleSearchHelmCharts, readOnly: true, toolName: ToolSearchHelmCharts},
				{name: "list_helm_releases", handler: (*PortainerMCPServer).HandleListHelmReleases, readOnly: true, toolName: ToolListHelmReleases},
				{name: "get_helm_release_history", handler: (*PortainerMCPServer).HandleGetHelmReleaseHistory, readOnly: true, toolName: ToolGetHelmReleaseHistory},
				{name: "add_helm_repository", handler: (*PortainerMCPServer).HandleAddHelmRepository, readOnly: false, toolName: ToolAddHelmRepository},
				{name: "remove_helm_repository", handler: (*PortainerMCPServer).HandleRemoveHelmRepository, readOnly: false, toolName: ToolRemoveHelmRepository},
				{name: "install_helm_chart", handler: (*PortainerMCPServer).HandleInstallHelmChart, readOnly: false, toolName: ToolInstallHelmChart},
				{name: "delete_helm_release", handler: (*PortainerMCPServer).HandleDeleteHelmRelease, readOnly: false, toolName: ToolDeleteHelmRelease},
			},
			annotation: mcp.ToolAnnotation{
				Title:           "Manage Helm",
				ReadOnlyHint:    boolPtr(false),
				DestructiveHint: boolPtr(true),
				IdempotentHint:  boolPtr(false),
				OpenWorldHint:   boolPtr(false),
			},
		},
		{
			name:        "manage_registries",
			description: "Manage container registries (Quay, Azure, DockerHub, GitLab, ECR, custom). Actions: list_registries, get_registry, create_registry, update_registry, delete_registry. Set 'action' parameter to choose.",
			actions: []metaAction{
				{name: "list_registries", handler: (*PortainerMCPServer).HandleListRegistries, readOnly: true, toolName: ToolListRegistries},
				{name: "get_registry", handler: (*PortainerMCPServer).HandleGetRegistry, readOnly: true, toolName: ToolGetRegistry},
				{name: "create_registry", handler: (*PortainerMCPServer).HandleCreateRegistry, readOnly: false, toolName: ToolCreateRegistry},
				{name: "update_registry", handler: (*PortainerMCPServer).HandleUpdateRegistry, readOnly: false, toolName: ToolUpdateRegistry},
				{name: "delete_registry", handler: (*PortainerMCPServer).HandleDeleteRegistry, readOnly: false, toolName: ToolDeleteRegistry},
			},
			annotation: mcp.ToolAnnotation{
				Title:           "Manage Registries",
				ReadOnlyHint:    boolPtr(false),
				DestructiveHint: boolPtr(true),
				IdempotentHint:  boolPtr(false),
				OpenWorldHint:   boolPtr(false),
			},
		},
		{
			name:        "manage_templates",
			description: "Manage custom and application templates for stack deployment. Actions: list_custom_templates, get_custom_template, get_custom_template_file, create_custom_template, delete_custom_template, list_app_templates, get_app_template_file. Set 'action' parameter to choose.",
			actions: []metaAction{
				{name: "list_custom_templates", handler: (*PortainerMCPServer).HandleListCustomTemplates, readOnly: true, toolName: ToolListCustomTemplates},
				{name: "get_custom_template", handler: (*PortainerMCPServer).HandleGetCustomTemplate, readOnly: true, toolName: ToolGetCustomTemplate},
				{name: "get_custom_template_file", handler: (*PortainerMCPServer).HandleGetCustomTemplateFile, readOnly: true, toolName: ToolGetCustomTemplateFile},
				{name: "create_custom_template", handler: (*PortainerMCPServer).HandleCreateCustomTemplate, readOnly: false, toolName: ToolCreateCustomTemplate},
				{name: "delete_custom_template", handler: (*PortainerMCPServer).HandleDeleteCustomTemplate, readOnly: false, toolName: ToolDeleteCustomTemplate},
				{name: "list_app_templates", handler: (*PortainerMCPServer).HandleListAppTemplates, readOnly: true, toolName: ToolListAppTemplates},
				{name: "get_app_template_file", handler: (*PortainerMCPServer).HandleGetAppTemplateFile, readOnly: true, toolName: ToolGetAppTemplateFile},
			},
			annotation: mcp.ToolAnnotation{
				Title:           "Manage Templates",
				ReadOnlyHint:    boolPtr(false),
				DestructiveHint: boolPtr(true),
				IdempotentHint:  boolPtr(false),
				OpenWorldHint:   boolPtr(false),
			},
		},
		{
			name:        "manage_backups",
			description: "Manage Portainer server backups and restore (local and S3). Actions: get_backup_status, get_backup_s3_settings, create_backup, backup_to_s3, restore_from_s3. Set 'action' parameter to choose.",
			actions: []metaAction{
				{name: "get_backup_status", handler: (*PortainerMCPServer).HandleGetBackupStatus, readOnly: true, toolName: ToolGetBackupStatus},
				{name: "get_backup_s3_settings", handler: (*PortainerMCPServer).HandleGetBackupS3Settings, readOnly: true, toolName: ToolGetBackupS3Settings},
				{name: "create_backup", handler: (*PortainerMCPServer).HandleCreateBackup, readOnly: false, toolName: ToolCreateBackup},
				{name: "backup_to_s3", handler: (*PortainerMCPServer).HandleBackupToS3, readOnly: false, toolName: ToolBackupToS3},
				{name: "restore_from_s3", handler: (*PortainerMCPServer).HandleRestoreFromS3, readOnly: false, toolName: ToolRestoreFromS3},
			},
			annotation: mcp.ToolAnnotation{
				Title:           "Manage Backups",
				ReadOnlyHint:    boolPtr(false),
				DestructiveHint: boolPtr(true),
				IdempotentHint:  boolPtr(false),
				OpenWorldHint:   boolPtr(false),
			},
		},
		{
			name:        "manage_webhooks",
			description: "Manage webhooks for container services and automated deployments. Actions: list_webhooks, create_webhook, delete_webhook. Set 'action' parameter to choose.",
			actions: []metaAction{
				{name: "list_webhooks", handler: (*PortainerMCPServer).HandleListWebhooks, readOnly: true, toolName: ToolListWebhooks},
				{name: "create_webhook", handler: (*PortainerMCPServer).HandleCreateWebhook, readOnly: false, toolName: ToolCreateWebhook},
				{name: "delete_webhook", handler: (*PortainerMCPServer).HandleDeleteWebhook, readOnly: false, toolName: ToolDeleteWebhook},
			},
			annotation: mcp.ToolAnnotation{
				Title:           "Manage Webhooks",
				ReadOnlyHint:    boolPtr(false),
				DestructiveHint: boolPtr(true),
				IdempotentHint:  boolPtr(false),
				OpenWorldHint:   boolPtr(false),
			},
		},
		{
			name:        "manage_edge",
			description: "Manage Edge compute jobs and update schedules for remote environments. Actions: list_edge_jobs, get_edge_job, get_edge_job_file, create_edge_job, delete_edge_job, list_edge_update_schedules. Set 'action' parameter to choose.",
			actions: []metaAction{
				{name: "list_edge_jobs", handler: (*PortainerMCPServer).HandleListEdgeJobs, readOnly: true, toolName: ToolListEdgeJobs},
				{name: "get_edge_job", handler: (*PortainerMCPServer).HandleGetEdgeJob, readOnly: true, toolName: ToolGetEdgeJob},
				{name: "get_edge_job_file", handler: (*PortainerMCPServer).HandleGetEdgeJobFile, readOnly: true, toolName: ToolGetEdgeJobFile},
				{name: "create_edge_job", handler: (*PortainerMCPServer).HandleCreateEdgeJob, readOnly: false, toolName: ToolCreateEdgeJob},
				{name: "delete_edge_job", handler: (*PortainerMCPServer).HandleDeleteEdgeJob, readOnly: false, toolName: ToolDeleteEdgeJob},
				{name: "list_edge_update_schedules", handler: (*PortainerMCPServer).HandleListEdgeUpdateSchedules, readOnly: true, toolName: ToolListEdgeUpdateSchedules},
			},
			annotation: mcp.ToolAnnotation{
				Title:           "Manage Edge",
				ReadOnlyHint:    boolPtr(false),
				DestructiveHint: boolPtr(true),
				IdempotentHint:  boolPtr(false),
				OpenWorldHint:   boolPtr(false),
			},
		},
		{
			name:        "manage_settings",
			description: "Manage Portainer server settings, public settings, and SSL configuration. Actions: get_settings, get_public_settings, update_settings, get_ssl_settings, update_ssl_settings. Set 'action' parameter to choose.",
			actions: []metaAction{
				{name: "get_settings", handler: (*PortainerMCPServer).HandleGetSettings, readOnly: true, toolName: ToolGetSettings},
				{name: "get_public_settings", handler: (*PortainerMCPServer).HandleGetPublicSettings, readOnly: true, toolName: ToolGetPublicSettings},
				{name: "update_settings", handler: (*PortainerMCPServer).HandleUpdateSettings, readOnly: false, toolName: ToolUpdateSettings},
				{name: "get_ssl_settings", handler: (*PortainerMCPServer).HandleGetSSLSettings, readOnly: true, toolName: ToolGetSSLSettings},
				{name: "update_ssl_settings", handler: (*PortainerMCPServer).HandleUpdateSSLSettings, readOnly: false, toolName: ToolUpdateSSLSettings},
			},
			annotation: mcp.ToolAnnotation{
				Title:           "Manage Settings",
				ReadOnlyHint:    boolPtr(false),
				DestructiveHint: boolPtr(false),
				IdempotentHint:  boolPtr(true),
				OpenWorldHint:   boolPtr(false),
			},
		},
		{
			name:        "manage_system",
			description: "Portainer system info, roles, MOTD, and authentication. Actions: get_system_status, list_roles, get_motd, authenticate, logout. Set 'action' parameter to choose.",
			actions: []metaAction{
				{name: "get_system_status", handler: (*PortainerMCPServer).HandleGetSystemStatus, readOnly: true, toolName: ToolGetSystemStatus},
				{name: "list_roles", handler: (*PortainerMCPServer).HandleListRoles, readOnly: true, toolName: ToolListRoles},
				{name: "get_motd", handler: (*PortainerMCPServer).HandleGetMOTD, readOnly: true, toolName: ToolGetMOTD},
				{name: "authenticate", handler: (*PortainerMCPServer).HandleAuthenticateUser, readOnly: true, toolName: ToolAuthenticate},
				{name: "logout", handler: (*PortainerMCPServer).HandleLogout, readOnly: false, toolName: ToolLogout},
			},
			annotation: mcp.ToolAnnotation{
				Title:           "Manage System",
				ReadOnlyHint:    boolPtr(false),
				DestructiveHint: boolPtr(false),
				IdempotentHint:  boolPtr(false),
				OpenWorldHint:   boolPtr(false),
			},
		},
	}
}
