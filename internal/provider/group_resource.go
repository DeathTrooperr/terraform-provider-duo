package provider

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/srmullaney/terraform-provider-duo/internal/duo"
)

// GroupResource defines the resource implementation.
type GroupResource struct {
	client *duo.Client
}

// GroupResourceModel describes the resource data model.
type GroupResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Status      types.String `tfsdk:"status"`
	PushEnabled types.Bool   `tfsdk:"push_enabled"`
	SmsEnabled  types.Bool   `tfsdk:"sms_enabled"`
}

func NewGroupResource() resource.Resource {
	return &GroupResource{}
}

func (r *GroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *GroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Duo Group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The Duo Group ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the group.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description of the group.",
			},
			"status": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The status of the group (Active, Bypass, Disabled). Defaults to Active.",
			},
			"push_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether Push is enabled for this group.",
			},
			"sms_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether SMS is enabled for this group.",
			},
		},
	}
}

func (r *GroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *GroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data GroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := url.Values{}
	params.Set("name", data.Name.ValueString())
	if !data.Description.IsNull() {
		params.Set("desc", data.Description.ValueString())
	}
	if !data.Status.IsNull() {
		params.Set("status", data.Status.ValueString())
	}
	if !data.PushEnabled.IsNull() {
		params.Set("push_enabled", fmt.Sprintf("%t", data.PushEnabled.ValueBool()))
	}
	if !data.SmsEnabled.IsNull() {
		params.Set("sms_enabled", fmt.Sprintf("%t", data.SmsEnabled.ValueBool()))
	}

	result, err := r.client.CreateGroup(params)
	if err != nil {
		resp.Diagnostics.AddError("Error creating group", err.Error())
		return
	}

	data.ID = types.StringValue(result.GroupId)
	data.Name = types.StringValue(result.Name)
	data.Description = types.StringValue(result.Desc)
	data.Status = types.StringValue(result.Status)
	data.PushEnabled = types.BoolValue(result.PushEnabled)
	data.SmsEnabled = types.BoolValue(result.SmsEnabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data GroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	groups, err := r.client.GetGroups()
	if err != nil {
		resp.Diagnostics.AddError("Error reading groups", err.Error())
		return
	}

	var found *duo.Group
	for _, g := range groups {
		if g.GroupId == data.ID.ValueString() {
			found = &g
			break
		}
	}

	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data.Name = types.StringValue(found.Name)
	data.Description = types.StringValue(found.Desc)
	data.Status = types.StringValue(found.Status)
	data.PushEnabled = types.BoolValue(found.PushEnabled)
	data.SmsEnabled = types.BoolValue(found.SmsEnabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data GroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := url.Values{}
	params.Set("name", data.Name.ValueString())
	if !data.Description.IsNull() {
		params.Set("desc", data.Description.ValueString())
	} else {
		params.Set("desc", "")
	}
	if !data.Status.IsNull() {
		params.Set("status", data.Status.ValueString())
	}
	if !data.PushEnabled.IsNull() {
		params.Set("push_enabled", fmt.Sprintf("%t", data.PushEnabled.ValueBool()))
	}
	if !data.SmsEnabled.IsNull() {
		params.Set("sms_enabled", fmt.Sprintf("%t", data.SmsEnabled.ValueBool()))
	}

	result, err := r.client.UpdateGroup(data.ID.ValueString(), params)
	if err != nil {
		resp.Diagnostics.AddError("Error updating group", err.Error())
		return
	}

	data.Name = types.StringValue(result.Name)
	data.Description = types.StringValue(result.Desc)
	data.Status = types.StringValue(result.Status)
	data.PushEnabled = types.BoolValue(result.PushEnabled)
	data.SmsEnabled = types.BoolValue(result.SmsEnabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data GroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteGroup(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting group", err.Error())
		return
	}
}
