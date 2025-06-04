package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/go-openapi/runtime"
	rtc "github.com/go-openapi/runtime/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	api "github.com/kelvintaywl/circleci-go-sdk/client"
	rapi "github.com/kelvintaywl/circleci-runner-go-sdk/client"
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
	Client       *api.Circleci
	RunnerClient *rapi.Circleci
	V1Client     *http.Client
	Hostname     string
	Auth         runtime.ClientAuthInfoWriter
}

type httpClientTransport struct {
	APIToken string
}

func (t *httpClientTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Circle-Token", t.APIToken)
	return http.DefaultTransport.RoundTrip(req)
}

// retryOn429Transport wraps an http.RoundTripper to retry on HTTP 429 responses.
type retryOn429Transport struct {
	Base       http.RoundTripper
	MaxRetries int
	Backoff    func(attempt int) int // returns milliseconds
}

func (t *retryOn429Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	attempts := 0
	for {
		resp, err := t.Base.RoundTrip(req)
		if err != nil {
			return resp, err
		}
		if resp.StatusCode != 429 || attempts >= t.MaxRetries {
			return resp, err
		}
		resp.Body.Close()
		backoffMs := t.Backoff(attempts)
		tflog.Debug(req.Context(), fmt.Sprintf("Received HTTP 429, retrying in %d ms (attempt %d/%d)", backoffMs, attempts+1, t.MaxRetries))
		time.Sleep(time.Duration(backoffMs) * time.Millisecond)
		attempts++
	}
}

// CircleciProviderModel describes the provider data model.
type CircleciProviderModel struct {
	ApiToken   types.String `tfsdk:"api_token"`
	Hostname   types.String `tfsdk:"hostname"`
	Retry      types.Bool   `tfsdk:"retry"`
	MaxRetries types.Int64  `tfsdk:"max_retries"`
}

func (p *CircleciProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "circleci"
	resp.Version = p.version
}

func (p *CircleciProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				MarkdownDescription: "A CircleCI user API token. This can also be set via the `CIRCLE_TOKEN` environment variable.",
				Optional:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("CircleCI hostname (default: %s). This can also be set via the `CIRCLE_HOSTNAME` environment variable.", defaultHostName),
				Optional:            true,
			},
			"retry": schema.BoolAttribute{
				MarkdownDescription: "Whether to retry API calls when provider receives an HTTP 429 status code (default: false).",
				Optional:            true,
			},
			"max_retries": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of retries for API calls when retry is enabled (default: 3).",
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

	// Retry settings
	retry := false
	maxRetries := int64(3)

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

	if data.Retry.ValueBool() {
		retry = data.Retry.ValueBool()
	}

	if data.MaxRetries.ValueInt64() > 0 {
		maxRetries = data.MaxRetries.ValueInt64()
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
		tflog.Info(ctx, fmt.Sprintf("Using default value for hostname: %s", hostname))

	}

	cfg := api.DefaultTransportConfig().WithHost(hostname)

	client := api.NewHTTPClientWithConfig(strfmt.Default, cfg)
	auth := rtc.APIKeyAuth("Circle-Token", "header", apiToken)

	rhostname := hostname
	if hostname == defaultHostName {
		rhostname = fmt.Sprintf("runner.%s", defaultHostName)
	}
	rcfg := rapi.DefaultTransportConfig().WithHost(rhostname)

	defaultTransport := &httpClientTransport{
		APIToken: apiToken,
	}

	var rclient *rapi.Circleci
	if retry {
		// Add retry transport for 429s
		retryTransport := &retryOn429Transport{
			Base:       defaultTransport,
			MaxRetries: int(maxRetries),
			Backoff: func(attempt int) int {
				// Exponential backoff: 500ms, 1000ms, 2000ms
				return 500 << attempt
			},
		}
		transport := rtc.New(rcfg.Host, rcfg.BasePath, rcfg.Schemes)
		transport.Transport = retryTransport
		rclient = rapi.New(transport, strfmt.Default)
	} else {
		rclient = rapi.NewHTTPClientWithConfig(strfmt.Default, rcfg)
	}

	httpClient := &http.Client{Transport: defaultTransport}

	apiClient := &CircleciAPIClient{
		Client:       client,
		RunnerClient: rclient,
		V1Client:     httpClient,
		Hostname:     hostname,
		Auth:         auth,
	}

	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

func (p *CircleciProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewWebhookResource,
		NewScheduleResource,
		NewEnvVarResource,
		NewCheckoutKeyResource,
		NewContextResource,
		NewContextEnvVarResource,
		NewRunnerResourceClassResource,
		NewRunnerTokenResource,
		NewProjectResource,
	}
}

func (p *CircleciProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectDataSource,
		NewWebhooksDataSource,
		NewCheckoutKeysDataSource,
		NewContextDataSource,
		NewRunnerResourceClassesDataSource,
		NewRunnerTokensDataSource,
	}
}
