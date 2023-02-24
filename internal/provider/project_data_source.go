package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/kelvintaywl/circleci-go-sdk/client/project"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &ProjectDataSource{}

func NewProjectDataSource() datasource.DataSource {
	return &ProjectDataSource{}
}

type ProjectDataSource struct {
	client *CircleciAPIClient
}

// ProjectDataSourceModel describes the data source data model.
type ProjectDataSourceModel struct {
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

func (d *ProjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *ProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fetches the information for a project",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of this project",
				Computed:            true,
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

func (d *ProjectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	param := project.NewGetProjectParams().WithDefaults()
	param = param.WithProjectSlug(data.Slug.ValueString())

	res, err := d.client.Client.Project.GetProject(param, d.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError("Encountered error fetching API", fmt.Sprintf("%s", err))
		return
	}

	pj := res.GetPayload()
	data.Id = types.StringValue(pj.ID.String())
	data.Name = types.StringValue(pj.Name)
	data.OrganizationName = types.StringValue(pj.OrganizationName)
	data.OrganizationSlug = types.StringValue(pj.OrganizationSlug)
	data.OrganizationId = types.StringValue(pj.OrganizationID.String())
	data.VcsProvider = types.StringValue(pj.VcsInfo.Provider)
	data.VcsDefaultBranch = types.StringValue(pj.VcsInfo.DefaultBranch)
	data.VcsURL = types.StringValue(pj.VcsInfo.VcsURL)

	// Save data into Terraform state
	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
