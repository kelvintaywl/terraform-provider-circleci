package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/kelvintaywl/circleci-go-sdk/client/project"
	"github.com/kelvintaywl/circleci-go-sdk/models"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &EnvVarResource{}

func NewEnvVarResource() resource.Resource {
	return &EnvVarResource{}
}

type EnvVarResource struct {
	client *CircleciAPIClient
}

type EnvVarResourceModel struct {
	ProjectSlug types.String `tfsdk:"project_slug"`
	Name        types.String `tfsdk:"name"`
	Value       types.String `tfsdk:"value"`
	Id          types.String `tfsdk:"id"`
}

func (r *EnvVarResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_env_var"
}

func (r *EnvVarResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a project environment variable",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Read-only unique identifier, set as {project_slug}/{name}",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the environment variable",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "The value of the environment variable",
				Required:            true,
				Sensitive:           true,
			},
			"project_slug": schema.StringAttribute{
				MarkdownDescription: "The project-slug for the environment variable",
				Required:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *EnvVarResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *EnvVarResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state EnvVarResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := state.Name.ValueString()
	projectSlug := state.ProjectSlug.ValueString()
	param := project.NewGetProjectEnvVarParamsWithContext(ctx).WithDefaults()
	param = param.WithProjectSlug(projectSlug).WithName(name)

	_, err := r.client.Client.Project.GetProjectEnvVar(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Encountered error reading Project(%s) Env Var %s", projectSlug, name), fmt.Sprintf("%s", err))
		return
	}

	state.Id = types.StringValue(fmt.Sprintf("%s/%s", projectSlug, name))
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *EnvVarResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan EnvVarResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectSlug := plan.ProjectSlug.ValueString()
	param := project.NewAddProjectEnvVarParamsWithContext(ctx).WithDefaults()
	param = param.WithProjectSlug(projectSlug)

	name := plan.Name.ValueString()
	value := plan.Value.ValueString()
	body := models.ProjectEnvVarPayload{
		Name:  &name,
		Value: &value,
	}

	param = param.WithBody(&body)

	_, err := r.client.Client.Project.AddProjectEnvVar(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project env var",
			fmt.Sprintf("Could not create project env var, unexpected error: %s", err.Error()),
		)
		return
	}

	plan.Id = types.StringValue(fmt.Sprintf("%s/%s", projectSlug, name))
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EnvVarResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan EnvVarResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectSlug := plan.ProjectSlug.ValueString()
	name := plan.Name.ValueString()

	deleteParam := project.NewDeleteProjectEnvVarParamsWithContext(ctx).WithDefaults()
	deleteParam = deleteParam.WithProjectSlug(projectSlug).WithName(name)

	_, err := r.client.Client.Project.DeleteProjectEnvVar(deleteParam, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting project env var",
			fmt.Sprintf("Could not delete project(%s) env var %s, unexpected error: %s", projectSlug, name, err.Error()),
		)
		return
	}

	value := plan.Value.ValueString()
	body := models.ProjectEnvVarPayload{
		Name:  &name,
		Value: &value,
	}

	addParam := project.NewAddProjectEnvVarParamsWithContext(ctx).WithDefaults()
	addParam = addParam.WithProjectSlug(projectSlug)
	addParam = addParam.WithBody(&body)

	_, err = r.client.Client.Project.AddProjectEnvVar(addParam, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error recreating project env var",
			fmt.Sprintf("Could not recreate project env var, unexpected error: %s", err.Error()),
		)
		return
	}

	plan.Id = types.StringValue(fmt.Sprintf("%s/%s", projectSlug, name))
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EnvVarResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state EnvVarResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := state.Name.ValueString()
	projectSlug := state.ProjectSlug.ValueString()

	param := project.NewDeleteProjectEnvVarParamsWithContext(ctx).WithDefaults()
	param = param.WithProjectSlug(projectSlug).WithName(name)

	_, err := r.client.Client.Project.DeleteProjectEnvVar(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting project env var",
			fmt.Sprintf("Could not delete project(%s) env var %s, unexpected error: %s", projectSlug, name, err.Error()),
		)
		return
	}
}
