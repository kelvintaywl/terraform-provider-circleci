package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/go-openapi/strfmt"

	"github.com/go-openapi/runtime"
	rtc "github.com/go-openapi/runtime/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	api "github.com/kelvintaywl/circleci-webhook-go-sdk/client"
)

// Ensure CircleciProvider satisfies various provider interfaces.
var _ provider.Provider = &CircleciProvider{}

const (
	defaultHostName string = "circleci.com"
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &CircleciProvider{}
}

// CircleciProvider defines the provider implementation.
type CircleciProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type CircleciAPIClient struct {
	Client *api.Circleci
	Auth   runtime.ClientAuthInfoWriter
}

// CircleciProviderModel describes the provider data model.
type CircleciProviderModel struct {
	ApiToken types.String `tfsdk:"api_token"`
	Hostname types.String `tfsdk:"hostname"`
	Https    types.Bool   `tfsdk:"https"`
}

func (p *CircleciProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "circleci"
	resp.Version = p.version
}

func (p *CircleciProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				MarkdownDescription: "CircleCI User API token",
				Optional:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("CircleCI hostname (default: %s)", defaultHostName),
				Optional:            true,
			},
			"https": schema.BoolAttribute{
				MarkdownDescription: "Use HTTPS (default: true)",
				Optional:            true,
			},
		},
	}
}

func (p *CircleciProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data CircleciProviderModel

	// Check environment variables
	apiToken := os.Getenv("CIRCLE_TOKEN")
	hostname := os.Getenv("CIRCLE_HOSTNAME")

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check configuration data, which should take precedence over
	// environment variable data, if found.
	if data.ApiToken.ValueString() != "" {
		apiToken = data.ApiToken.ValueString()
	}

	if data.Hostname.ValueString() != "" {
		hostname = data.Hostname.ValueString()
	}

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing CircleCI user API Token configuration",
			"While configuring the provider, the CircleCI user API token was not found in "+
				"the CIRCLE_TOKEN environment variable or provider "+
				"configuration block api_token attribute.",
		)
		// Not returning early allows the logic to collect all errors.
	}

	if hostname == "" {
		hostname = defaultHostName
		resp.Diagnostics.AddWarning(
			"Missing CircleCI API hostname configuration",
			"While configuring the provider, the CircleCI API hostname was not found in "+
				"the CIRCLE_HOSTNAME environment variable or provider "+
				fmt.Sprintf("configuration block hostname attribute.\nUsing default: %s", hostname),
		)
	}

	cfg := api.DefaultTransportConfig().WithHost(hostname)

	// only for acceptance testing
	if data.Https.ValueBool() == false {
		cfg = cfg.WithSchemes([]string{"http"})
	}

	client := api.NewHTTPClientWithConfig(strfmt.Default, cfg)
	auth := rtc.APIKeyAuth("Circle-Token", "header", apiToken)
	apiClient := &CircleciAPIClient{
		Client: client,
		Auth:   auth,
	}

	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

func (p *CircleciProvider) Resources(ctx context.Context) []func() resource.Resource {
	// TODO:
	return []func() resource.Resource{}
}

func (p *CircleciProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewWebhooksDataSource,
	}
}
