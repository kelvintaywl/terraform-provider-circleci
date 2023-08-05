package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-openapi/strfmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/kelvintaywl/circleci-go-sdk/client/webhook"
	"github.com/kelvintaywl/circleci-go-sdk/models"
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
	Id            types.String `tfsdk:"id"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
	Name          types.String `tfsdk:"name"`
	URL           types.String `tfsdk:"url"`
	SigningSecret types.String `tfsdk:"signing_secret"`
	ProjectID     types.String `tfsdk:"project_id"`
	VerifyTLS     types.Bool   `tfsdk:"verify_tls"`
	Events        types.Set    `tfsdk:"events"`
}

var vEvents = []string{
	"job-completed",
	"workflow-completed",
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
				// unchanged even during updates
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^https:\/\/.+`),
						"must start with https://"),
				},
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
			"events": schema.SetAttribute{
				MarkdownDescription: fmt.Sprintf("Events that will trigger the webhook. Allowed values: %v", vEvents),
				ElementType:         types.StringType,
				Required:            true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.OneOf(vEvents...)),
				},
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
			"Unexpected Resource Configure Type",
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

	id := state.Id.ValueString()
	param := webhook.NewGetWebhookParamsWithContext(ctx).WithDefaults()
	param = param.WithID(strfmt.UUID(id))

	res, err := r.client.Client.Webhook.GetWebhook(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Encountered error reading Webhook %s", id), fmt.Sprintf("%s", err))
		return
	}

	w := res.GetPayload()

	state.Id = types.StringValue(w.ID.String())
	state.CreatedAt = types.StringValue(w.CreatedAt.String())
	state.UpdatedAt = types.StringValue(w.UpdatedAt.String())
	state.Name = types.StringValue(w.Name)
	state.URL = types.StringValue(w.URL)
	state.SigningSecret = types.StringValue(w.SigningSecret)
	state.VerifyTLS = types.BoolValue(*w.VerifyTLS)
	state.ProjectID = types.StringValue(w.Scope.ID.String())
	state.Events, _ = types.SetValueFrom(ctx, types.StringType, w.Events)

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

	var events []string
	diags.Append(plan.Events.ElementsAs(ctx, &events, false)...)
	body.Events = events
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
	plan.Id = types.StringValue(id)
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

	id := plan.Id.ValueString()

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
	var events []string
	diags.Append(plan.Events.ElementsAs(ctx, &events, false)...)
	body.Events = events

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

	id := state.Id.ValueString()

	param := webhook.NewDeleteWebhookParamsWithContext(ctx).WithDefaults()
	param = param.WithID(strfmt.UUID(id))

	_, err := r.client.Client.Webhook.DeleteWebhook(param, r.client.Auth)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "not found") {
			tflog.Warn(ctx, fmt.Sprintf("Webhook no longer found: %s", id))
			return
		}

		resp.Diagnostics.AddError(
			"Error deleting webhook",
			fmt.Sprintf("Could not delete webhook %s, unexpected error: %s", id, errMsg),
		)
		return
	}
}

func (r *WebhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
