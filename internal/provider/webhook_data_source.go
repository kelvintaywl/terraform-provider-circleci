package provider

import (
	"context"
	"fmt"

	"github.com/go-openapi/strfmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/kelvintaywl/circleci-webhook-go-sdk/client/webhook"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &WebhookDataSource{}

func NewWebhookDataSource() datasource.DataSource {
	return &WebhookDataSource{}
}

type WebhookDataSource struct {
	client *CircleciAPIClient
}

// WebhookDataSourceModel describes the data source data model.
type WebhookDataSourceModel struct {
	ProjectId types.String `tfsdk:"project_id"`
}

func (d *WebhookDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

func (d *WebhookDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "CircleCI project webhook data source.",

		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				MarkdownDescription: "CircleCI project ID.",
				Required:            true,
			},
		},
	}
}

func (d *WebhookDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*CircleciAPIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *CircleciAPIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *WebhookDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WebhookDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	param := webhook.NewListWebhooksParams().WithDefaults()
	param = param.WithScopeID(strfmt.UUID(data.ProjectId.ValueString()))

	res, err := d.client.Client.Webhook.ListWebhooks(param, d.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError("Encountered error fetching API", fmt.Sprintf("%s", err))
		return
	}

	if ! res.IsSuccess() {
		resp.Diagnostics.AddError("Failed to fetch API with HTTP 4xx or 5xx errors", res.Error())
		return
	}

	webhooksInfo := res.GetPayload()
	// TODO: paginate
	webhooks := webhooksInfo.Items
	for i, w := range webhooks {
		fmt.Printf("%dth schedule: %s, %s, %s, %v", (i + 1), w.ID, w.Name, w.URL, w.Events)
	}
	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	// data.Id = types.StringValue("example-id")

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
