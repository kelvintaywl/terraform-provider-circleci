package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/kelvintaywl/circleci-go-sdk/client/project"
	"github.com/kelvintaywl/circleci-go-sdk/models"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &CheckoutKeyResource{}

func NewCheckoutKeyResource() resource.Resource {
	return &CheckoutKeyResource{}
}

type CheckoutKeyResource struct {
	client *CircleciAPIClient
}

type CheckoutKeyResourceModel struct {
	ProjectSlug types.String `tfsdk:"project_slug"`
	PublicKey   types.String `tfsdk:"public_key"`
	Fingerprint types.String `tfsdk:"fingerprint"`
	Type        types.String `tfsdk:"type"`
	Preferred   types.Bool   `tfsdk:"preferred"`
	CreatedAt   types.String `tfsdk:"created_at"`
	Id          types.String `tfsdk:"id"`
}

var vKeyTypes = []string{
	"deploy-key",
	"user-key",
}

func (r *CheckoutKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_checkout_key"
}

func (r *CheckoutKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a project checkout key",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Read-only unique identifier: uses fingerprint",
				Computed:            true,
				// unchanged even during updates
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_slug": schema.StringAttribute{
				MarkdownDescription: "The project-slug for the environment variable",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of checkout key to create. This may be either `deploy-key` or `user-key`",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(vKeyTypes...),
				},
				// if modifed, this requires a replacement instead.
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"public_key": schema.StringAttribute{
				MarkdownDescription: "A public SSH key",
				Computed:            true,
			},
			"fingerprint": schema.StringAttribute{
				MarkdownDescription: "An SSH key fingerprint",
				Computed:            true,
			},
			"preferred": schema.BoolAttribute{
				MarkdownDescription: "A boolean value that indicates if this key is preferred",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The date and time the checkout key was created",
				Computed:            true,
				// unchanged even during updates
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *CheckoutKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *CheckoutKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state CheckoutKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	fingerprint := state.Fingerprint.ValueString()
	projectSlug := state.ProjectSlug.ValueString()
	param := project.NewGetProjectCheckoutKeyParamsWithContext(ctx).WithDefaults()
	param = param.WithProjectSlug(projectSlug).WithFingerprint(fingerprint)

	res, err := r.client.Client.Project.GetProjectCheckoutKey(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Encountered error reading Project(%s) checkout key %s", projectSlug, fingerprint), fmt.Sprintf("%s", err))
		return
	}

	ck := res.GetPayload()

	state.Id = types.StringValue(ck.Fingerprint)
	state.Fingerprint = types.StringValue(ck.Fingerprint)
	state.PublicKey = types.StringValue(ck.PublicKey)
	keyType := ck.Type
	if ck.Type == "github-user-key" {
		keyType = "user-key"
	}
	state.Type = types.StringValue(keyType)
	state.Preferred = types.BoolValue(*ck.Preferred)
	state.CreatedAt = types.StringValue(ck.CreatedAt.String())

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *CheckoutKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan CheckoutKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectSlug := plan.ProjectSlug.ValueString()
	param := project.NewAddProjectCheckoutKeyParamsWithContext(ctx).WithDefaults()
	param = param.WithProjectSlug(projectSlug)

	keyType := plan.Type.ValueString()
	body := models.ProjectCheckoutKeyPayload{
		Type: keyType,
	}

	param = param.WithBody(&body)

	res, err := r.client.Client.Project.AddProjectCheckoutKey(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project checkout key",
			fmt.Sprintf("Could not create project checkout key, unexpected error: %s", err.Error()),
		)
		return
	}

	ck := res.GetPayload()
	plan.Id = types.StringValue(ck.Fingerprint)
	plan.Fingerprint = types.StringValue(ck.Fingerprint)
	plan.PublicKey = types.StringValue(ck.PublicKey)
	keyType = ck.Type
	if ck.Type == "github-user-key" {
		keyType = "user-key"
	}
	plan.Type = types.StringValue(keyType)
	plan.Preferred = types.BoolValue(*ck.Preferred)
	plan.CreatedAt = types.StringValue(ck.CreatedAt.String())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *CheckoutKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// not implemented; requires a replacement
}

func (r *CheckoutKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state CheckoutKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectSlug := state.ProjectSlug.ValueString()
	fingerprint := state.Fingerprint.ValueString()

	param := project.NewDeleteProjectCheckoutKeyParamsWithContext(ctx).WithDefaults()
	param = param.WithProjectSlug(projectSlug).WithFingerprint(fingerprint)

	_, err := r.client.Client.Project.DeleteProjectCheckoutKey(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting project checkout key",
			fmt.Sprintf("Could not delete project(%s) checkout key %s, unexpected error: %s", projectSlug, fingerprint, err.Error()),
		)
		return
	}
}
