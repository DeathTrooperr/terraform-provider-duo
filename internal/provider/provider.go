package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/srmullaney/terraform-provider-duo/internal/duo"
)

// DuoProvider defines the provider implementation.
type DuoProvider struct {
	version string
}

// DuoProviderModel describes the provider data model.
type DuoProviderModel struct {
	IntegrationKey types.String `tfsdk:"integration_key"`
	SecretKey      types.String `tfsdk:"secret_key"`
	ApiHost        types.String `tfsdk:"api_host"`
}

func (p *DuoProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "duo"
	resp.Version = p.version
}

func (p *DuoProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"integration_key": schema.StringAttribute{
				MarkdownDescription: "Duo Admin API integration key.",
				Optional:            true,
			},
			"secret_key": schema.StringAttribute{
				MarkdownDescription: "Duo Admin API secret key.",
				Optional:            true,
				Sensitive:           true,
			},
			"api_host": schema.StringAttribute{
				MarkdownDescription: "Duo API host (e.g. api-XXXXXXXX.duosecurity.com).",
				Optional:            true,
			},
		},
	}
}

func (p *DuoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data DuoProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Default to environment variables if not provided in configuration
	ikey := os.Getenv("DUO_INTEGRATION_KEY")
	skey := os.Getenv("DUO_SECRET_KEY")
	host := os.Getenv("DUO_API_HOST")

	if !data.IntegrationKey.IsNull() {
		ikey = data.IntegrationKey.ValueString()
	}
	if !data.SecretKey.IsNull() {
		skey = data.SecretKey.ValueString()
	}
	if !data.ApiHost.IsNull() {
		host = data.ApiHost.ValueString()
	}

	if ikey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("integration_key"),
			"Missing Duo Integration Key",
			"The provider cannot create the Duo API client as there is no integration key configuration.",
		)
	}

	if skey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("secret_key"),
			"Missing Duo Secret Key",
			"The provider cannot create the Duo API client as there is no secret key configuration.",
		)
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_host"),
			"Missing Duo API Host",
			"The provider cannot create the Duo API client as there is no API host configuration.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create Duo Admin API Client
	duoClient := duo.NewClient(ikey, skey, host)

	resp.DataSourceData = duoClient
	resp.ResourceData = duoClient
}

func (p *DuoProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewGroupResource,
		NewApplicationResource,
		NewPolicyResource,
		NewSettingsResource,
		NewUserResource,
	}
}

func (p *DuoProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) provider.Provider {
	return &DuoProvider{
		version: version,
	}
}
