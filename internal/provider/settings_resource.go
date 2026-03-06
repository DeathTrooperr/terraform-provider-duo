package provider

import (
	"context"
	"fmt"
	"net/url"

	"github.com/DeathTrooperr/terraform-provider-duo/internal/duo"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SettingsResource defines the resource implementation.
type SettingsResource struct {
	client *duo.Client
}

// SettingsResourceModel describes the resource data model.
type SettingsResourceModel struct {
	LockoutThreshold types.Int64 `tfsdk:"lockout_threshold"`
	LockoutDuration  types.Int64 `tfsdk:"lockout_duration"`
	InactiveTimeout  types.Int64 `tfsdk:"inactive_timeout"`
	UserApprove      types.Bool  `tfsdk:"user_approve"`
	UserTelephony    types.Bool  `tfsdk:"user_telephony"`
}

func NewSettingsResource() resource.Resource {
	return &SettingsResource{}
}

func (r *SettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_settings"
}

func (r *SettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Duo General Settings. Note: This resource manages global settings, only one should exist.",
		Attributes: map[string]schema.Attribute{
			"lockout_threshold": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Number of failed logins before lockout.",
			},
			"lockout_duration": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Lockout duration in minutes.",
			},
			"inactive_timeout": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Inactive timeout in minutes.",
			},
			"user_approve": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether users can approve their own devices.",
			},
			"user_telephony": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether users can use telephony for authentication.",
			},
		},
	}
}

func (r *SettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*duo.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *duo.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *SettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// In Duo Admin API, settings are updated via UpdateSettings.
	// There is no "Create" for global settings, so we just update them.
	r.updateSettings(ctx, data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SettingsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	settings, err := r.client.GetSettings()
	if err != nil {
		resp.Diagnostics.AddError("Error reading settings", err.Error())
		return
	}

	// Map API response to model.
	data.LockoutThreshold = types.Int64Value(int64(settings.LockoutThreshold))
	data.LockoutDuration = types.Int64Value(int64(settings.LockoutDuration))
	data.InactiveTimeout = types.Int64Value(int64(settings.InactiveExpiration))
	data.UserApprove = types.BoolValue(settings.UserApprove)
	data.UserTelephony = types.BoolValue(settings.UserTelephony)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.updateSettings(ctx, data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Global settings cannot be deleted, only reset or left as is.
	// We just remove it from Terraform state.
}

func (r *SettingsResource) updateSettings(ctx context.Context, data SettingsResourceModel, diags *diag.Diagnostics) {
	params := url.Values{}
	if !data.LockoutThreshold.IsUnknown() && !data.LockoutThreshold.IsNull() {
		params.Set("lockout_threshold", fmt.Sprintf("%d", data.LockoutThreshold.ValueInt64()))
	}
	if !data.LockoutDuration.IsUnknown() && !data.LockoutDuration.IsNull() {
		params.Set("lockout_duration", fmt.Sprintf("%d", data.LockoutDuration.ValueInt64()))
	}
	if !data.InactiveTimeout.IsUnknown() && !data.InactiveTimeout.IsNull() {
		params.Set("inactive_expiration", fmt.Sprintf("%d", data.InactiveTimeout.ValueInt64()))
	}
	if !data.UserApprove.IsUnknown() && !data.UserApprove.IsNull() {
		params.Set("user_approve", fmt.Sprintf("%t", data.UserApprove.ValueBool()))
	}
	if !data.UserTelephony.IsUnknown() && !data.UserTelephony.IsNull() {
		params.Set("user_telephony", fmt.Sprintf("%t", data.UserTelephony.ValueBool()))
	}

	_, err := r.client.UpdateSettings(params)
	if err != nil {
		diags.AddError("Error updating settings", err.Error())
	}
}
