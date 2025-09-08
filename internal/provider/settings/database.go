package settings

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/supabase/cli/pkg/api"
)

// DatabaseConfig represents PostgreSQL database configuration
type DatabaseConfig struct {
	EffectiveCacheSize            types.String `tfsdk:"effective_cache_size"`
	LogicalDecodingWorkMem        types.String `tfsdk:"logical_decoding_work_mem"`
	MaintenanceWorkMem            types.String `tfsdk:"maintenance_work_mem"`
	MaxConnections                types.Int64  `tfsdk:"max_connections"`
	MaxLocksPerTransaction        types.Int64  `tfsdk:"max_locks_per_transaction"`
	MaxParallelMaintenanceWorkers types.Int64  `tfsdk:"max_parallel_maintenance_workers"`
	MaxParallelWorkers            types.Int64  `tfsdk:"max_parallel_workers"`
	MaxParallelWorkersPerGather   types.Int64  `tfsdk:"max_parallel_workers_per_gather"`
	MaxReplicationSlots           types.Int64  `tfsdk:"max_replication_slots"`
	MaxSlotWalKeepSize            types.String `tfsdk:"max_slot_wal_keep_size"`
	MaxStandbyArchiveDelay        types.String `tfsdk:"max_standby_archive_delay"`
	MaxStandbyStreamingDelay      types.String `tfsdk:"max_standby_streaming_delay"`
	MaxWalSenders                 types.Int64  `tfsdk:"max_wal_senders"`
	MaxWalSize                    types.String `tfsdk:"max_wal_size"`
	MaxWorkerProcesses            types.Int64  `tfsdk:"max_worker_processes"`
	RestartDatabase               types.Bool   `tfsdk:"restart_database"`
	SessionReplicationRole        types.String `tfsdk:"session_replication_role"`
	SharedBuffers                 types.String `tfsdk:"shared_buffers"`
	StatementTimeout              types.String `tfsdk:"statement_timeout"`
	TrackCommitTimestamp          types.Bool   `tfsdk:"track_commit_timestamp"`
	WalKeepSize                   types.String `tfsdk:"wal_keep_size"`
	WalSenderTimeout              types.String `tfsdk:"wal_sender_timeout"`
	WorkMem                       types.String `tfsdk:"work_mem"`
}

func GetDatabaseSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"effective_cache_size": schema.StringAttribute{
			MarkdownDescription: "Amount of memory available for disk caching by the OS and within the database itself",
			Optional:            true,
		},
		"logical_decoding_work_mem": schema.StringAttribute{
			MarkdownDescription: "Memory used for logical decoding",
			Optional:            true,
		},
		"maintenance_work_mem": schema.StringAttribute{
			MarkdownDescription: "Maximum amount of memory to be used by maintenance operations",
			Optional:            true,
		},
		"max_connections": schema.Int64Attribute{
			MarkdownDescription: "Maximum number of concurrent connections to the database server",
			Optional:            true,
		},
		"max_locks_per_transaction": schema.Int64Attribute{
			MarkdownDescription: "Maximum number of locks per transaction",
			Optional:            true,
		},
		"max_parallel_maintenance_workers": schema.Int64Attribute{
			MarkdownDescription: "Maximum number of parallel maintenance workers",
			Optional:            true,
		},
		"max_parallel_workers": schema.Int64Attribute{
			MarkdownDescription: "Maximum number of parallel worker processes",
			Optional:            true,
		},
		"max_parallel_workers_per_gather": schema.Int64Attribute{
			MarkdownDescription: "Maximum number of parallel workers per Gather node",
			Optional:            true,
		},
		"max_replication_slots": schema.Int64Attribute{
			MarkdownDescription: "Maximum number of replication slots",
			Optional:            true,
		},
		"max_slot_wal_keep_size": schema.StringAttribute{
			MarkdownDescription: "Maximum size of WAL files that replication slots are allowed to retain",
			Optional:            true,
		},
		"max_standby_archive_delay": schema.StringAttribute{
			MarkdownDescription: "Maximum delay before canceling queries when a hot standby server is processing archived WAL data",
			Optional:            true,
		},
		"max_standby_streaming_delay": schema.StringAttribute{
			MarkdownDescription: "Maximum delay before canceling queries when a hot standby server is processing streamed WAL data",
			Optional:            true,
		},
		"max_wal_senders": schema.Int64Attribute{
			MarkdownDescription: "Maximum number of WAL sender processes",
			Optional:            true,
		},
		"max_wal_size": schema.StringAttribute{
			MarkdownDescription: "Maximum size to let the WAL grow during automatic checkpoints",
			Optional:            true,
		},
		"max_worker_processes": schema.Int64Attribute{
			MarkdownDescription: "Maximum number of background worker processes",
			Optional:            true,
		},
		"restart_database": schema.BoolAttribute{
			MarkdownDescription: "Whether to restart the database to apply configuration changes",
			Optional:            true,
		},
		"session_replication_role": schema.StringAttribute{
			MarkdownDescription: "Controls firing of replication-related triggers and rules (origin, replica, local)",
			Optional:            true,
		},
		"shared_buffers": schema.StringAttribute{
			MarkdownDescription: "Amount of memory the database server uses for shared memory buffers",
			Optional:            true,
		},
		"statement_timeout": schema.StringAttribute{
			MarkdownDescription: "Maximum allowed duration of any statement",
			Optional:            true,
		},
		"track_commit_timestamp": schema.BoolAttribute{
			MarkdownDescription: "Whether to track commit time stamps of transactions",
			Optional:            true,
		},
		"wal_keep_size": schema.StringAttribute{
			MarkdownDescription: "Minimum size to retain in the pg_wal directory",
			Optional:            true,
		},
		"wal_sender_timeout": schema.StringAttribute{
			MarkdownDescription: "Maximum time to wait for WAL replication",
			Optional:            true,
		},
		"work_mem": schema.StringAttribute{
			MarkdownDescription: "Amount of memory to be used by internal sort operations and hash tables",
			Optional:            true,
		},
	}
}

// ReadDatabaseConfig reads database configuration from the API
func ReadDatabaseConfig(ctx context.Context, client *api.ClientWithResponses, state *SettingsResourceModel) diag.Diagnostics {
	httpResp, err := client.V1GetPostgresConfigWithResponse(ctx, state.Id.ValueString())
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read database settings: %s", err))}
	}

	switch httpResp.StatusCode() {
	case http.StatusNotFound, http.StatusNotAcceptable:
		return nil
	}

	if httpResp.JSON200 == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read database settings, got status %d: %s", httpResp.StatusCode(), httpResp.Body))}
	}

	if state.Database == nil {
		state.Database = &DatabaseConfig{}
	}

	resp := httpResp.JSON200

	if resp.EffectiveCacheSize != nil {
		state.Database.EffectiveCacheSize = types.StringValue(*resp.EffectiveCacheSize)
	}
	if resp.LogicalDecodingWorkMem != nil {
		state.Database.LogicalDecodingWorkMem = types.StringValue(*resp.LogicalDecodingWorkMem)
	}
	if resp.MaintenanceWorkMem != nil {
		state.Database.MaintenanceWorkMem = types.StringValue(*resp.MaintenanceWorkMem)
	}
	if resp.MaxConnections != nil {
		state.Database.MaxConnections = types.Int64Value(int64(*resp.MaxConnections))
	}
	if resp.StatementTimeout != nil {
		state.Database.StatementTimeout = types.StringValue(*resp.StatementTimeout)
	}
	if resp.SharedBuffers != nil {
		state.Database.SharedBuffers = types.StringValue(*resp.SharedBuffers)
	}
	if resp.WorkMem != nil {
		state.Database.WorkMem = types.StringValue(*resp.WorkMem)
	}

	return nil
}

// UpdateDatabaseConfig updates database configuration via the API
func UpdateDatabaseConfig(ctx context.Context, client *api.ClientWithResponses, plan *SettingsResourceModel) diag.Diagnostics {
	body := api.UpdatePostgresConfigBody{}

	if !plan.Database.EffectiveCacheSize.IsNull() {
		body.EffectiveCacheSize = plan.Database.EffectiveCacheSize.ValueStringPointer()
	}
	if !plan.Database.LogicalDecodingWorkMem.IsNull() {
		body.LogicalDecodingWorkMem = plan.Database.LogicalDecodingWorkMem.ValueStringPointer()
	}
	if !plan.Database.MaintenanceWorkMem.IsNull() {
		body.MaintenanceWorkMem = plan.Database.MaintenanceWorkMem.ValueStringPointer()
	}
	if !plan.Database.MaxConnections.IsNull() {
		val := int(plan.Database.MaxConnections.ValueInt64())
		body.MaxConnections = &val
	}
	if !plan.Database.StatementTimeout.IsNull() {
		body.StatementTimeout = plan.Database.StatementTimeout.ValueStringPointer()
	}
	if !plan.Database.SharedBuffers.IsNull() {
		body.SharedBuffers = plan.Database.SharedBuffers.ValueStringPointer()
	}
	if !plan.Database.WorkMem.IsNull() {
		body.WorkMem = plan.Database.WorkMem.ValueStringPointer()
	}
	if !plan.Database.RestartDatabase.IsNull() {
		body.RestartDatabase = plan.Database.RestartDatabase.ValueBoolPointer()
	}

	httpResp, err := client.V1UpdatePostgresConfigWithResponse(ctx, plan.ProjectRef.ValueString(), body)
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to update database settings: %s", err))}
	}

	if httpResp.JSON200 == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to update database settings, got status %d: %s", httpResp.StatusCode(), httpResp.Body))}
	}

	return nil
}
