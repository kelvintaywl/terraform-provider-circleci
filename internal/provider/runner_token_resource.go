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

	"github.com/kelvintaywl/circleci-runner-go-sdk/client/token"
	"github.com/kelvintaywl/circleci-runner-go-sdk/models"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &RunnerTokenResource{}

func NewRunnerTokenResource() resource.Resource {
	return &RunnerTokenResource{}
}

type RunnerTokenResource struct {
	client *CircleciAPIClient
}

type RunnerTokenResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Nickname      types.String `tfsdk:"nickname"`
	ResourceClass types.String `tfsdk:"resource_class"`
	Token         types.String `tfsdk:"token"`
	CreatedAt     types.String `tfsdk:"created_at"`
}

func (r *RunnerTokenResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_runner_token"
}

func (r *RunnerTokenResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Runner token",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the Runner token.",
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
			"nickname": schema.StringAttribute{
				MarkdownDescription: "The Runner token alias.",
				Required:            true,
				// if modifed, this requires a replacement instead.
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "The Runner token value.",
				Computed:            true,
				Sensitive:           true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Date and time the token was created",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *RunnerTokenResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *RunnerTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state RunnerTokenResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	resourceClass := state.ResourceClass.ValueString()

	param := token.NewListTokensParamsWithContext(ctx).WithDefaults()
	param = param.WithResourceClass(resourceClass)

	res, err := r.client.RunnerClient.Token.ListTokens(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError("Encountered error fetching API", fmt.Sprintf("%s", err))
		return
	}

	info := res.GetPayload()
	for _, tk := range info.Items {
		if id == tk.ID.String() {
			state.Id = types.StringValue(tk.ID.String())
			state.Nickname = types.StringValue(*tk.Nickname)
			state.ResourceClass = types.StringValue(*tk.ResourceClass)
			state.CreatedAt = types.StringValue(tk.CreatedAt.String())
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
func (r *RunnerTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan RunnerTokenResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	param := token.NewCreateTokenParamsWithContext(ctx).WithDefaults()
	resourceClass := plan.ResourceClass.ValueString()
	nickname := plan.Nickname.ValueString()

	body := models.TokenPayload{
		ResourceClass: &resourceClass,
		Nickname:      &nickname,
	}

	param = param.WithBody(&body)

	res, err := r.client.RunnerClient.Token.CreateToken(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating token",
			fmt.Sprintf("Could not create Runner token, unexpected error: %s", err.Error()),
		)
		return
	}

	rc := res.GetPayload()
	plan.Id = types.StringValue(rc.ID.String())
	plan.Token = types.StringValue(rc.Token)
	plan.CreatedAt = types.StringValue(rc.CreatedAt.String())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *RunnerTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// not implemented; requires a replacement
}

func (r *RunnerTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state RunnerTokenResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()

	param := token.NewDeleteTokenParamsWithContext(ctx).WithDefaults()
	param = param.WithID(strfmt.UUID(id))

	_, err := r.client.RunnerClient.Token.DeleteToken(param, r.client.Auth)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "not found") {
			tflog.Warn(ctx, fmt.Sprintf("Runner token no longer found: %s", id))
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting Runner token",
			fmt.Sprintf("Could not delete token %s, unexpected error: %s", id, errMsg),
		)
		return
	}
}
