package provider

import (
	"context"
	"fmt"
	"regexp"
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

	"github.com/kelvintaywl/circleci-runner-go-sdk/client/resource_class"
	"github.com/kelvintaywl/circleci-runner-go-sdk/models"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &RunnerResourceClassResource{}

func NewRunnerResourceClassResource() resource.Resource {
	return &RunnerResourceClassResource{}
}

type RunnerResourceClassResource struct {
	client *CircleciAPIClient
}

type RunnerResourceClassResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Description   types.String `tfsdk:"description"`
	ResourceClass types.String `tfsdk:"resource_class"`
}

func (r *RunnerResourceClassResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_runner_resource_class"
}

func (r *RunnerResourceClassResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Runner resource-class",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the Runner resource-class",
				Computed:            true,
				// unchanged even during updates
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_class": schema.StringAttribute{
				MarkdownDescription: "The name of the Runner resource-class (should include namespace)",
				Required:            true,
				// if modifed, this requires a replacement instead.
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^.+\/.+`),
						"must follow <namespace>/<name>"),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description for the Runner resource-class",
				Required:            true,
				// if modifed, this requires a replacement instead.
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *RunnerResourceClassResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *RunnerResourceClassResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state RunnerResourceClassResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	resourceClass := state.ResourceClass.ValueString()
	namespaceName := strings.SplitN(resourceClass, "/", 2)
	if len(namespaceName) != 2 {
		resp.Diagnostics.AddError("Invalid resource class", resourceClass)
		return
	}
	namespace := namespaceName[0]

	param := resource_class.NewListResourceClassesParamsWithContext(ctx).WithDefaults()
	param = param.WithNamespace(namespace)

	res, err := r.client.RunnerClient.ResourceClass.ListResourceClasses(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError("Encountered error fetching API", fmt.Sprintf("%s", err))
		return
	}

	info := res.GetPayload()
	for _, rc := range info.Items {
		if id == rc.ID.String() {
			state.Id = types.StringValue(rc.ID.String())
			state.Description = types.StringValue(*rc.Description)
			state.ResourceClass = types.StringValue(*rc.ResourceClass)
			// Set refreshed state
			diags = resp.State.Set(ctx, &state)
			resp.Diagnostics.Append(diags...)
			break
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *RunnerResourceClassResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan RunnerResourceClassResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	param := resource_class.NewCreateResourceClassParamsWithContext(ctx).WithDefaults()
	resourceClass := plan.ResourceClass.ValueString()
	desc := plan.Description.ValueString()

	body := models.ResourceClassPayload{
		ResourceClass: &resourceClass,
		Description:   &desc,
	}

	param = param.WithBody(&body)

	res, err := r.client.RunnerClient.ResourceClass.CreateResourceClass(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating resource-class",
			fmt.Sprintf("Could not create resource-class, unexpected error: %s", err.Error()),
		)
		return
	}

	rc := res.GetPayload()
	plan.Id = types.StringValue(rc.ID.String())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *RunnerResourceClassResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// not implemented; requires a replacement
}

func (r *RunnerResourceClassResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state RunnerResourceClassResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()

	param := resource_class.NewDeleteResourceClassParamsWithContext(ctx).WithDefaults()
	param = param.WithID(strfmt.UUID(id))

	_, err := r.client.RunnerClient.ResourceClass.DeleteResourceClass(param, r.client.Auth)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "not found") {
			tflog.Warn(ctx, fmt.Sprintf("Resource-class no longer found: %s", id))
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting resource-class",
			fmt.Sprintf("Could not delete resource-class %s, unexpected error: %s", id, errMsg),
		)
		return
	}
}
