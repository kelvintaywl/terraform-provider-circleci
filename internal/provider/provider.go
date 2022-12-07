package circleci

import (
	"context"

	"github.com/go-openapi/runtime"
	rtc "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	api "github.com/kelvintaywl/circleci-schedule-go-sdk/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type apiClient struct {
	// Add whatever fields, client or connection info, etc. here
	// you would need to setup to communicate with the upstream
	// API.
	Client *api.Circleci
	Auth   runtime.ClientAuthInfoWriter
}

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"hostname": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://circleci.com",
				DefaultFunc: schema.EnvDefaultFunc("CIRCLECI_HOSTNAME", nil),
			},
			"api_token": &schema.Schema{
				Type:        schema.TypeString,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CIRCLECI_API_TOKEN", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"circleci_schedule": resourceSchedule(),
		},
		// DataSourcesMap: map[string]*schema.Resource{
		// 	"hashicups_coffees":     dataSourceCoffees(),
		// 	"hashicups_order":       dataSourceOrder(),
		// 	"hashicups_ingredients": dataSourceIngredients(),
		// },
	}
}

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := Provider()
		p.ConfigureContextFunc = configure(version, p)
		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (any, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		hostname := d.Get("hostname").(string)
		apiToken := d.Get("api_token").(string)

		// Warning or errors can be collected in a slice type
		var diags diag.Diagnostics

		if (hostname != "") && (apiToken != "") {
			cfg := api.DefaultTransportConfig().WithHost(hostname)
			client := api.NewHTTPClientWithConfig(strfmt.Default, cfg)
			auth := rtc.APIKeyAuth("Circle-Token", "header", apiToken)

			return &apiClient{Client: client, Auth: auth}, diags
		}

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create CircleCI API client",
			Detail:   "Missing required hostname and API token",
		})

		return nil, diags
	}
}
