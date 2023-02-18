package provider

import (
	"context"
	"fmt"

	"github.com/go-openapi/strfmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/kelvintaywl/circleci-go-sdk/client/webhook"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &WebhooksDataSource{}

func NewWebhooksDataSource() datasource.DataSource {
	return &WebhooksDataSource{}
}

type WebhooksDataSource struct {
	client *CircleciAPIClient
}

// WebhookDataSourceModel describes the data source data model.
type WebhooksDataSourceModel struct {
	ProjectId types.String   `tfsdk:"project_id"`
	Webhooks  []webhookModel `tfsdk:"webhooks"`
	Id        types.String   `tfsdk:"id"`
}

type webhookModel struct {
	ID            types.String   `tfsdk:"id"`
	Name          types.String   `tfsdk:"name"`
	URL           types.String   `tfsdk:"url"`
	VerifyTLS     types.Bool     `tfsdk:"verify_tls"`
	SigningSecret types.String   `tfsdk:"signing_secret"`
	CreatedAt     types.String   `tfsdk:"created_at"`
	UpdatedAt     types.String   `tfsdk:"updated_at"`
	Scope         scopeModel     `tfsdk:"scope"`
	Events        []types.String `tfsdk:"events"`
}

type scopeModel struct {
	ID   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

func (d *WebhooksDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhooks"
}

func (d *WebhooksDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fetches the list of webhooks for a specific project.",

		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				MarkdownDescription: "CircleCI project ID.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of this data source: project ID.",
				Computed:            true,
			},
			"webhooks": schema.ListNestedAttribute{
				MarkdownDescription: "List of webhooks",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique ID of the webhook",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the webhook",
							Computed:            true,
						},
						"url": schema.StringAttribute{
							MarkdownDescription: "URL to deliver the webhook to",
							Computed:            true,
						},
						"signing_secret": schema.StringAttribute{
							MarkdownDescription: "Masked value of the secret used to build an HMAC hash of the payload",
							Computed:            true,
						},
						"created_at": schema.StringAttribute{
							MarkdownDescription: "The date and time the webhook was created",
							Computed:            true,
						},
						"updated_at": schema.StringAttribute{
							MarkdownDescription: "The date and time the webhook was updated",
							Computed:            true,
						},
						"verify_tls": schema.BoolAttribute{
							MarkdownDescription: "Whether to enforce TLS certificate verification when delivering the webhook",
							Computed:            true,
						},
						"events": schema.ListAttribute{
							MarkdownDescription: "Events that will trigger the webhook",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"scope": schema.SingleNestedAttribute{
							MarkdownDescription: "Scope",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									MarkdownDescription: "Scope ID",
									Computed:            true,
								},
								"type": schema.StringAttribute{
									MarkdownDescription: "Scope type (project)",
									Computed:            true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *WebhooksDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WebhooksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WebhooksDataSourceModel

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

	info := res.GetPayload()

	// token := info.NextPageToken
	// TODO: consider support of pagination
	// However, it seems the ListWebhooks API does not support page-token query params.
	// See https://circleci.com/docs/api/v2/index.html#operation/getWebhooks

	for _, w := range info.Items {
		webhookState := webhookModel{
			ID:            types.StringValue(w.ID.String()),
			Name:          types.StringValue(w.Name),
			URL:           types.StringValue(w.URL),
			VerifyTLS:     types.BoolValue(*w.VerifyTLS),
			SigningSecret: types.StringValue(w.SigningSecret),
			CreatedAt:     types.StringValue(w.CreatedAt.String()),
			UpdatedAt:     types.StringValue(w.UpdatedAt.String()),
			// NOTE: Scope values MUST be returned;
			// we can assume this, based on https://circleci.com/docs/api/v2/index.html#operation/getWebhooks
			Scope: scopeModel{
				ID:   types.StringValue(w.Scope.ID.String()),
				Type: types.StringValue(*w.Scope.Type),
			},
		}
		for _, event := range w.Events {
			webhookState.Events = append(webhookState.Events, types.StringValue(event))
		}
		data.Webhooks = append(data.Webhooks, webhookState)
	}
	data.Id = data.ProjectId

	// Save data into Terraform state
	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
