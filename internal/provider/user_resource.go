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

// UserResource defines the resource implementation.
type UserResource struct {
	client *duo.Client
}

// UserResourceModel describes the resource data model.
type UserResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Username  types.String `tfsdk:"username"`
	FirstName types.String `tfsdk:"firstname"`
	LastName  types.String `tfsdk:"lastname"`
	RealName  types.String `tfsdk:"realname"`
	Email     types.String `tfsdk:"email"`
	Status    types.String `tfsdk:"status"`
	Notes     types.String `tfsdk:"notes"`
	Alias1    types.String `tfsdk:"alias1"`
	Alias2    types.String `tfsdk:"alias2"`
	Alias3    types.String `tfsdk:"alias3"`
	Alias4    types.String `tfsdk:"alias4"`
}

func NewUserResource() resource.Resource {
	return &UserResource{}
}

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Duo User.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The Duo User ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"username": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The username of the user.",
			},
			"firstname": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The first name of the user.",
			},
			"lastname": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The last name of the user.",
			},
			"realname": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The real name of the user.",
			},
			"email": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The email of the user.",
			},
			"status": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The status of the user (Active, Bypass, Disabled). Defaults to Active.",
			},
			"notes": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The notes for the user.",
			},
			"alias1": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "An alias for the user.",
			},
			"alias2": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "An alias for the user.",
			},
			"alias3": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "An alias for the user.",
			},
			"alias4": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "An alias for the user.",
			},
		},
	}
}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := url.Values{}
	params.Set("username", data.Username.ValueString())
	if !data.FirstName.IsNull() {
		params.Set("firstname", data.FirstName.ValueString())
	}
	if !data.LastName.IsNull() {
		params.Set("lastname", data.LastName.ValueString())
	}
	if !data.RealName.IsNull() {
		params.Set("realname", data.RealName.ValueString())
	}
	if !data.Email.IsNull() {
		params.Set("email", data.Email.ValueString())
	}
	if !data.Status.IsNull() {
		params.Set("status", data.Status.ValueString())
	}
	if !data.Notes.IsNull() {
		params.Set("notes", data.Notes.ValueString())
	}
	if !data.Alias1.IsNull() {
		params.Set("alias1", data.Alias1.ValueString())
	}
	if !data.Alias2.IsNull() {
		params.Set("alias2", data.Alias2.ValueString())
	}
	if !data.Alias3.IsNull() {
		params.Set("alias3", data.Alias3.ValueString())
	}
	if !data.Alias4.IsNull() {
		params.Set("alias4", data.Alias4.ValueString())
	}

	result, err := r.client.CreateUser(params)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user", err.Error())
		return
	}

	data.ID = types.StringValue(result.UserId)
	data.Username = types.StringValue(result.Username)
	data.FirstName = types.StringValue(result.FirstName)
	data.LastName = types.StringValue(result.LastName)
	data.RealName = types.StringValue(result.RealName)
	data.Email = types.StringValue(result.Email)
	data.Status = types.StringValue(result.Status)
	data.Notes = types.StringValue(result.Notes)
	data.Alias1 = types.StringValue(result.Alias1)
	data.Alias2 = types.StringValue(result.Alias2)
	data.Alias3 = types.StringValue(result.Alias3)
	data.Alias4 = types.StringValue(result.Alias4)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.GetUser(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading user", err.Error())
		return
	}

	if user == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data.Username = types.StringValue(user.Username)
	data.FirstName = types.StringValue(user.FirstName)
	data.LastName = types.StringValue(user.LastName)
	data.RealName = types.StringValue(user.RealName)
	data.Email = types.StringValue(user.Email)
	data.Status = types.StringValue(user.Status)
	data.Notes = types.StringValue(user.Notes)
	data.Alias1 = types.StringValue(user.Alias1)
	data.Alias2 = types.StringValue(user.Alias2)
	data.Alias3 = types.StringValue(user.Alias3)
	data.Alias4 = types.StringValue(user.Alias4)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := url.Values{}
	params.Set("username", data.Username.ValueString())
	if !data.FirstName.IsNull() {
		params.Set("firstname", data.FirstName.ValueString())
	}
	if !data.LastName.IsNull() {
		params.Set("lastname", data.LastName.ValueString())
	}
	if !data.RealName.IsNull() {
		params.Set("realname", data.RealName.ValueString())
	}
	if !data.Email.IsNull() {
		params.Set("email", data.Email.ValueString())
	}
	if !data.Status.IsNull() {
		params.Set("status", data.Status.ValueString())
	}
	if !data.Notes.IsNull() {
		params.Set("notes", data.Notes.ValueString())
	}
	if !data.Alias1.IsNull() {
		params.Set("alias1", data.Alias1.ValueString())
	}
	if !data.Alias2.IsNull() {
		params.Set("alias2", data.Alias2.ValueString())
	}
	if !data.Alias3.IsNull() {
		params.Set("alias3", data.Alias3.ValueString())
	}
	if !data.Alias4.IsNull() {
		params.Set("alias4", data.Alias4.ValueString())
	}

	result, err := r.client.UpdateUser(data.ID.ValueString(), params)
	if err != nil {
		resp.Diagnostics.AddError("Error updating user", err.Error())
		return
	}

	data.Username = types.StringValue(result.Username)
	data.FirstName = types.StringValue(result.FirstName)
	data.LastName = types.StringValue(result.LastName)
	data.RealName = types.StringValue(result.RealName)
	data.Email = types.StringValue(result.Email)
	data.Status = types.StringValue(result.Status)
	data.Notes = types.StringValue(result.Notes)
	data.Alias1 = types.StringValue(result.Alias1)
	data.Alias2 = types.StringValue(result.Alias2)
	data.Alias3 = types.StringValue(result.Alias3)
	data.Alias4 = types.StringValue(result.Alias4)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUser(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting user", err.Error())
		return
	}
}
