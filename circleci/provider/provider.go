package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure CircleciProvider satisfies various provider interfaces.
var _ provider.Provider = &CircleciProvider{}
var _ provider.ProviderWithMetadata = &CircleciProvider{}

const (
	defaultHostName string = "https://circleci.com"
)

// CircleciProvider defines the provider implementation.
type CircleciProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// CircleciProviderModel describes the provider data model.
type CircleciProviderModel struct {
	ApiToken types.String `tfsdk:"api_token"`
	Hostname types.String `tfsdk:"hostname"`
}

func (p *CircleciProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "circleci"
	resp.Version = p.version
}

func (p *CircleciProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"api_token": {
				MarkdownDescription: "CircleCI user API token",
				Optional:            true,
				Type:                types.StringType,
			},
			"hostname": {
				MarkdownDescription: "CircleCI API server hostname",
				Optional:            true,
				Type:                types.StringType,
			},
		},
	}, nil
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
		resp.Diagnostics.AddError(
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
}

func (p *CircleciProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewScheduleResource,
	}
}

func (p *CircleciProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CircleciProvider{
			version: version,
		}
	}
}
