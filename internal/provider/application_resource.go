package provider

import (
	"context"
	"fmt"
	"net/url"

	"github.com/DeathTrooperr/terraform-provider-duo/internal/duo"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ApplicationResource defines the resource implementation.
type ApplicationResource struct {
	client *duo.Client
}

// ApplicationResourceModel describes the resource data model.
type ApplicationResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

func NewApplicationResource() resource.Resource {
	return &ApplicationResource{}
}

func (r *ApplicationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (r *ApplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Duo Application.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The Duo Application Integration Key (skey is NOT managed via this resource for security, but available in Duo GUI).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the application.",
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of the application (e.g. 'web_sdk', 'oidc', etc.).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *ApplicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ApplicationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := url.Values{}
	params.Set("name", data.Name.ValueString())
	params.Set("type", data.Type.ValueString())

	result, err := r.client.CreateApplication(params)
	if err != nil {
		resp.Diagnostics.AddError("Error creating application", err.Error())
		return
	}

	data.ID = types.StringValue(result.IntegrationKey)
	data.Name = types.StringValue(result.Name)
	data.Type = types.StringValue(result.Type)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ApplicationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apps, err := r.client.GetApplications()
	if err != nil {
		resp.Diagnostics.AddError("Error reading applications", err.Error())
		return
	}

	var found *duo.Application
	for _, a := range apps {
		if a.IntegrationKey == data.ID.ValueString() {
			found = &a
			break
		}
	}

	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data.Name = types.StringValue(found.Name)
	data.Type = types.StringValue(found.Type)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ApplicationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := url.Values{}
	params.Set("name", data.Name.ValueString())

	result, err := r.client.UpdateApplication(data.ID.ValueString(), params)
	if err != nil {
		resp.Diagnostics.AddError("Error updating application", err.Error())
		return
	}

	data.Name = types.StringValue(result.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ApplicationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteApplication(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting application", err.Error())
		return
	}
}
