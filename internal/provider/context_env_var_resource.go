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
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/kelvintaywl/circleci-go-sdk/client/contexts"
	"github.com/kelvintaywl/circleci-go-sdk/models"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &ContextEnvVarResource{}

func NewContextEnvVarResource() resource.Resource {
	return &ContextEnvVarResource{}
}

type ContextEnvVarResource struct {
	client *CircleciAPIClient
}

type ContextEnvVarResourceModel struct {
	ContextId types.String `tfsdk:"context_id"`
	Name      types.String `tfsdk:"name"`
	Value     types.String `tfsdk:"value"`
	Id        types.String `tfsdk:"id"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

func (r *ContextEnvVarResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_context_env_var"
}

func (r *ContextEnvVarResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a context environment variable",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Read-only unique identifier, set as {context_id}/{name}",
				Computed:            true,
				// unchanged even during updates
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The date and time the context environment variable was created",
				Computed:            true,
				// unchanged even during updates
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The date and time the context environment variable was last updated",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the context environment variable",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "The value of the context environment variable",
				Required:            true,
				Sensitive:           true,
			},
			"context_id": schema.StringAttribute{
				MarkdownDescription: "ID of the context",
				Required:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *ContextEnvVarResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ContextEnvVarResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ContextEnvVarResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := state.Name.ValueString()
	contextId := state.ContextId.ValueString()
	nextToken := ""

	for {
		param := contexts.NewListContextEnvVarsParamsWithContext(ctx).WithDefaults()
		param = param.WithID(strfmt.UUID(contextId)).WithPageToken(&nextToken)

		res, err := r.client.Client.Contexts.ListContextEnvVars(param, r.client.Auth)
		if err != nil {
			resp.Diagnostics.AddError("Encountered error fetching API", fmt.Sprintf("%s", err))
			return
		}

		info := res.GetPayload()
		nextToken = info.NextPageToken
		for _, ev := range info.Items {
			if ev.Variable == name {
				id := fmt.Sprintf("%s/%s", contextId, name)
				state.Id = types.StringValue(id)
				createdAt := ev.CreatedAt.String()
				state.CreatedAt = types.StringValue(createdAt)
				updatedAt := ev.UpdatedAt.String()
				state.UpdatedAt = types.StringValue(updatedAt)

				// Save data into Terraform state
				diags := resp.State.Set(ctx, &state)
				resp.Diagnostics.Append(diags...)
				if resp.Diagnostics.HasError() {
					return
				}
				return
			}
		}
		if nextToken == "" && state.Id.ValueString() == "" {
			resp.Diagnostics.AddError(fmt.Sprintf("Did not find context env var with name %s", name), fmt.Sprintf("%s", err))
			return
		}
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *ContextEnvVarResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ContextEnvVarResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	contextId := plan.ContextId.ValueString()
	name := plan.Name.ValueString()
	param := contexts.NewUpdateContextEnvVarParamsWithContext(ctx).WithDefaults()
	param = param.WithID(strfmt.UUID(contextId)).WithName(name)

	value := plan.Value.ValueString()
	body := models.ContextEnvVarPayload{
		Value: &value,
	}

	param = param.WithBody(&body)

	res, err := r.client.Client.Contexts.UpdateContextEnvVar(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating context env var",
			fmt.Sprintf("Could not create context env var, unexpected error: %s", err.Error()),
		)
		return
	}

	ev := res.GetPayload()

	plan.Id = types.StringValue(fmt.Sprintf("%s/%s", contextId, name))
	plan.CreatedAt = types.StringValue(ev.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(ev.UpdatedAt.String())
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ContextEnvVarResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ContextEnvVarResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	contextId := plan.ContextId.ValueString()
	name := plan.Name.ValueString()
	param := contexts.NewUpdateContextEnvVarParamsWithContext(ctx).WithDefaults()
	param = param.WithID(strfmt.UUID(contextId)).WithName(name)

	value := plan.Value.ValueString()
	body := models.ContextEnvVarPayload{
		Value: &value,
	}

	param = param.WithBody(&body)

	res, err := r.client.Client.Contexts.UpdateContextEnvVar(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating context env var",
			fmt.Sprintf("Could not update context env var, unexpected error: %s", err.Error()),
		)
		return
	}

	ev := res.GetPayload()

	plan.Id = types.StringValue(fmt.Sprintf("%s/%s", contextId, name))
	plan.CreatedAt = types.StringValue(ev.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(ev.UpdatedAt.String())
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ContextEnvVarResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ContextEnvVarResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := state.Name.ValueString()
	contextId := state.ContextId.ValueString()

	param := contexts.NewDeleteContextEnvVarParamsWithContext(ctx).WithDefaults()
	param = param.WithID(strfmt.UUID(contextId)).WithName(name)

	_, err := r.client.Client.Contexts.DeleteContextEnvVar(param, r.client.Auth)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "not found") {
			tflog.Warn(ctx, fmt.Sprintf("Context env var no longer found: %s/%s", contextId, name))
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting context env var",
			fmt.Sprintf("Could not delete context env var %s/%s, unexpected error: %s", contextId, name, errMsg),
		)
		return
	}
}
