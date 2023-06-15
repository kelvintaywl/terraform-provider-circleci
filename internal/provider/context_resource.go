package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-openapi/strfmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/kelvintaywl/circleci-go-sdk/client/contexts"
	"github.com/kelvintaywl/circleci-go-sdk/models"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &ContextResource{}

func NewContextResource() resource.Resource {
	return &ContextResource{}
}

type ContextResource struct {
	client *CircleciAPIClient
}

type ContextResourceModel struct {
	Id        types.String `tfsdk:"id"`
	CreatedAt types.String `tfsdk:"created_at"`
	Name      types.String `tfsdk:"name"`
	Owner     ownerModel   `tfsdk:"owner"`
}

type ownerModel struct {
	Id   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

var vOwnerTypes = []string{
	"account",
	"organization",
}

func (r *ContextResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_context"
}

func (r *ContextResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a context",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the context",
				Computed:            true,
				// unchanged even during updates
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The date and time the schedule was created",
				Computed:            true,
				// unchanged even during updates
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the context",
				Required:            true,
				// if modifed, this requires a replacement instead.
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"owner": schema.SingleNestedAttribute{
				MarkdownDescription: "The owner of the context",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "The unique ID of the owner",
						Required:            true,
					},
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the owner. Accepts `account` or `organization`. Accounts are only used as context owners in **Server**.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(vOwnerTypes...),
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *ContextResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ContextResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ContextResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	param := contexts.NewGetContextParamsWithContext(ctx).WithDefaults()
	param = param.WithID(strfmt.UUID(id))

	res, err := r.client.Client.Contexts.GetContext(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Encountered error reading Context %s", id), fmt.Sprintf("%s", err))
		return
	}

	c := res.GetPayload()

	state.Id = types.StringValue(c.ID.String())
	state.CreatedAt = types.StringValue(c.CreatedAt.String())
	state.Name = types.StringValue(c.Name)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *ContextResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ContextResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	param := contexts.NewAddContextParamsWithContext(ctx).WithDefaults()
	name := plan.Name.ValueString()
	ownerType := plan.Owner.Type.ValueString()
	if ownerType == "account" {
		msgOnlyForServer := "Owner Type: account requested. Make sure this is for a CircleCI Server (self-hosted) instance."
		tflog.Warn(ctx, msgOnlyForServer)
	}

	ownerID := strfmt.UUID(plan.Owner.Id.ValueString())
	owner := models.ContextPayloadOwner{
		Type: &ownerType,
		ID:   &ownerID,
	}
	body := models.ContextPayload{
		Name:  &name,
		Owner: &owner,
	}

	param = param.WithBody(&body)

	res, err := r.client.Client.Contexts.AddContext(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating context",
			fmt.Sprintf("Could not create context, unexpected error: %s", err.Error()),
		)
		return
	}

	c := res.GetPayload()
	plan.Id = types.StringValue(c.ID.String())
	plan.CreatedAt = types.StringValue(c.CreatedAt.String())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ContextResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// not implemented; requires a replacement
}

func (r *ContextResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ContextResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()

	param := contexts.NewDeleteContextParamsWithContext(ctx).WithDefaults()
	param = param.WithID(strfmt.UUID(id))

	_, err := r.client.Client.Contexts.DeleteContext(param, r.client.Auth)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "not found") {
			tflog.Warn(ctx, fmt.Sprintf("Context no longer found: %s", id))
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting context",
			fmt.Sprintf("Could not delete context %s, unexpected error: %s", id, errMsg),
		)
		return
	}
}
