package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/kelvintaywl/circleci-go-sdk/client/project"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &ProjectResource{}

func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

type ProjectResource struct {
	client *CircleciAPIClient
}

type ProjectResourceModel struct {
	Id               types.String `tfsdk:"id"`
	Slug             types.String `tfsdk:"slug"`
	Name             types.String `tfsdk:"name"`
	OrganizationName types.String `tfsdk:"organization_name"`
	OrganizationSlug types.String `tfsdk:"organization_slug"`
	OrganizationId   types.String `tfsdk:"organization_id"`
	VcsProvider      types.String `tfsdk:"vcs_provider"`
	VcsDefaultBranch types.String `tfsdk:"vcs_default_branch"`
	VcsURL           types.String `tfsdk:"vcs_url"`
}

func (r *ProjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a project",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Read-only unique identifier",
				Computed:            true,
				// unchanged even during updates
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "Project slug in the form `vcs-slug/org-name/repo-name`. The / characters may be URL-escaped.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the project",
				Computed:            true,
			},
			"organization_name": schema.StringAttribute{
				MarkdownDescription: "The name of the organization the project belongs to",
				Computed:            true,
			},
			"organization_slug": schema.StringAttribute{
				MarkdownDescription: "The slug of the organization the project belongs to",
				Computed:            true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "The id of the organization the project belongs to",
				Computed:            true,
			},
			"vcs_url": schema.StringAttribute{
				MarkdownDescription: "URL to the repository hosting the project's code",
				Computed:            true,
			},
			"vcs_provider": schema.StringAttribute{
				MarkdownDescription: "VCS provider (either GitHub, Bitbucket or CircleCI)",
				Computed:            true,
			},
			"vcs_default_branch": schema.StringAttribute{
				MarkdownDescription: "Default branch of this project",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *ProjectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ProjectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectSlug := state.Slug.ValueString()
	param := project.NewGetProjectParamsWithContext(ctx).WithDefaults()
	param = param.WithProjectSlug(projectSlug)

	res, err := r.client.Client.Project.GetProject(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Encountered error reading Project(%s)", projectSlug), fmt.Sprintf("%s", err))
		return
	}

	pj := res.GetPayload()
	state.Id = types.StringValue(pj.ID.String())
	state.Name = types.StringValue(pj.Name)
	state.OrganizationName = types.StringValue(pj.OrganizationName)
	state.OrganizationSlug = types.StringValue(pj.OrganizationSlug)
	state.OrganizationId = types.StringValue(pj.OrganizationID.String())
	state.VcsProvider = types.StringValue(pj.VcsInfo.Provider)
	state.VcsDefaultBranch = types.StringValue(pj.VcsInfo.DefaultBranch)
	state.VcsURL = types.StringValue(pj.VcsInfo.VcsURL)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ProjectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectSlug := plan.Slug.ValueString()
	url := fmt.Sprintf("https://%s/api/v1.1/project/%s/follow", r.client.Hostname, projectSlug)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Encountered error setting up API call", fmt.Sprintf("%s", err))
	}
	res, err := r.client.V1Client.Do(request)
	if err != nil || res.StatusCode != 200 {
		resp.Diagnostics.AddError(fmt.Sprintf("Encountered error following project (%s)", projectSlug), fmt.Sprintf("%s", err))
	}

	// read
	readParam := project.NewGetProjectParamsWithContext(ctx).WithDefaults()
	readParam = readParam.WithProjectSlug(projectSlug)

	readRes, err := r.client.Client.Project.GetProject(readParam, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Encountered error reading Project(%s)", projectSlug), fmt.Sprintf("%s", err))
		return
	}

	pj := readRes.GetPayload()
	plan.Id = types.StringValue(pj.ID.String())
	plan.Name = types.StringValue(pj.Name)
	plan.OrganizationName = types.StringValue(pj.OrganizationName)
	plan.OrganizationSlug = types.StringValue(pj.OrganizationSlug)
	plan.OrganizationId = types.StringValue(pj.OrganizationID.String())
	plan.VcsProvider = types.StringValue(pj.VcsInfo.Provider)
	plan.VcsDefaultBranch = types.StringValue(pj.VcsInfo.DefaultBranch)
	plan.VcsURL = types.StringValue(pj.VcsInfo.VcsURL)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// not implemented; not possible to update project
	tflog.Warn(ctx, "Project cannot be updated via this provider.")
}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// not implemented; not possible to delete project
	tflog.Warn(ctx, "Project cannot be deleted via this provider.")
}
