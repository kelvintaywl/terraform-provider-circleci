package provider

import (
	"context"
	"fmt"

	"github.com/go-openapi/strfmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/kelvintaywl/circleci-webhook-go-sdk/client/webhook"
	"github.com/kelvintaywl/circleci-webhook-go-sdk/models"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &WebhookResource{}

func NewWebhookResource() resource.Resource {
	return &WebhookResource{}
}

type WebhookResource struct {
	client *CircleciAPIClient
}

type WebhookResourceModel struct {
	ID            types.String   `tfsdk:"id"`
	CreatedAt     types.String   `tfsdk:"created_at"`
	UpdatedAt     types.String   `tfsdk:"updated_at"`
	Name          types.String   `tfsdk:"name"`
	URL           types.String   `tfsdk:"url"`
	SigningSecret types.String   `tfsdk:"signing_secret"`
	ProjectID     types.String   `tfsdk:"project_id"`
	VerifyTLS     types.Bool     `tfsdk:"verify_tls"`
	Events        []types.String `tfsdk:"events"`
}

func (r *WebhookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

func (r *WebhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a project webhook",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the webhook",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The date and time the webhook was created",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The date and time the webhook was last updated",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the webhook",
				Required:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL to deliver the webhook to. Note: protocol must be included as well (only https is supported)",
				Required:            true,
			},
			"signing_secret": schema.StringAttribute{
				MarkdownDescription: "Secret used to build an HMAC hash of the payload and passed as a header in the webhook request",
				Required:            true,
				Sensitive:           true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "ID of the project",
				Required:            true,
			},
			"verify_tls": schema.BoolAttribute{
				MarkdownDescription: "Whether to enforce TLS certificate verification when delivering the webhook",
				Required:            true,
			},
			"events": schema.ListAttribute{
				// TODO: consider validation here?
				MarkdownDescription: "Events that will trigger the webhook",
				ElementType:         types.StringType,
				Required:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *WebhookResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

// Read refreshes the Terraform state with the latest data.
func (r *WebhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state WebhookResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	param := webhook.NewGetWebhookParamsWithContext(ctx).WithDefaults()
	param = param.WithID(strfmt.UUID(id))

	res, err := r.client.Client.Webhook.GetWebhook(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Encountered error reading Webhook %s", id), fmt.Sprintf("%s", err))
		return
	}

	w := res.GetPayload()

	state.ID = types.StringValue(w.ID.String())
	state.CreatedAt = types.StringValue(w.CreatedAt.String())
	state.UpdatedAt = types.StringValue(w.UpdatedAt.String())
	state.Name = types.StringValue(w.Name)
	state.URL = types.StringValue(w.URL)
	state.SigningSecret = types.StringValue(w.SigningSecret)
	state.VerifyTLS = types.BoolValue(*w.VerifyTLS)
	state.ProjectID = types.StringValue(w.Scope.ID.String())
	for _, event := range w.Events {
		state.Events = append(state.Events, types.StringValue(event))
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *WebhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan WebhookResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	param := webhook.NewAddWebhookParamsWithContext(ctx).WithDefaults()
	project := "project"
	scope := models.WebhookBasePayloadScope{
		Type: &project,
		ID:   strfmt.UUID(plan.ProjectID.ValueString()),
	}
	verifyTLS := plan.VerifyTLS.ValueBool()
	body := models.WebhookPayloadForRequest{
		WebhookBasePayload: models.WebhookBasePayload{
			Name:  plan.Name.ValueString(),
			URL:   plan.URL.ValueString(),
			Scope: &scope,
		},
		SigningSecret: plan.SigningSecret.ValueString(),
		VerifyTLS:     &verifyTLS,
	}
	for _, event := range plan.Events {
		body.Events = append(body.Events, event.ValueString())
	}
	param = param.WithBody(&body)

	res, err := r.client.Client.Webhook.AddWebhook(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating webhook",
			fmt.Sprintf("Could not create webhook, unexpected error: %s", err.Error()),
		)
		return
	}

	w := res.GetPayload()

	id := w.ID.String()
	plan.ID = types.StringValue(id)
	plan.CreatedAt = types.StringValue(w.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(w.UpdatedAt.String())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *WebhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan WebhookResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()

	param := webhook.NewUpdateWebhookParamsWithContext(ctx).WithDefaults()
	param = param.WithID(strfmt.UUID(id))
	project := "project"
	scope := models.WebhookBasePayloadScope{
		Type: &project,
		ID:   strfmt.UUID(plan.ProjectID.ValueString()),
	}
	verifyTLS := plan.VerifyTLS.ValueBool()
	body := models.WebhookPayloadForRequest{
		WebhookBasePayload: models.WebhookBasePayload{
			Name:  plan.Name.ValueString(),
			URL:   plan.URL.ValueString(),
			Scope: &scope,
		},
		SigningSecret: plan.SigningSecret.ValueString(),
		VerifyTLS:     &verifyTLS,
	}
	for _, event := range plan.Events {
		body.Events = append(body.Events, event.ValueString())
	}

	param = param.WithBody(&body)
	res, err := r.client.Client.Webhook.UpdateWebhook(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating webhook",
			fmt.Sprintf("Could not update webhook %s, unexpected error: %s", id, err.Error()),
		)
		return
	}

	w := res.GetPayload()

	// NOTE: no need to update ID and ProjectID
	plan.CreatedAt = types.StringValue(w.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(w.UpdatedAt.String())
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *WebhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state WebhookResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	param := webhook.NewDeleteWebhookParamsWithContext(ctx).WithDefaults()
	param = param.WithID(strfmt.UUID(id))

	_, err := r.client.Client.Webhook.DeleteWebhook(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting webhook",
			fmt.Sprintf("Could not delete webhook %s, unexpected error: %s", id, err.Error()),
		)
		return
	}
}
