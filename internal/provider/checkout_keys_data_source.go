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
var _ datasource.DataSource = &CheckoutKeysDataSource{}

func NewCheckoutKeysDataSource() datasource.DataSource {
	return &CheckoutKeysDataSource{}
}

type CheckoutKeysDataSource struct {
	client *CircleciAPIClient
}

// CheckoutKeysDataSourceModel describes the data source data model.
type CheckoutKeysDataSourceModel struct {
	ProjectSlug types.String `tfsdk:"project_slug"`
	Keys        []keyModel   `tfsdk:"keys"`
	Id          types.String `tfsdk:"id"`
}

type keyModel struct {
	PublicKey   types.String `tfsdk:"public_key"`
	Type        types.String `tfsdk:"type"`
	Fingerprint types.String `tfsdk:"fingerprint"`
	Preferred   types.Bool   `tfsdk:"preferred"`
	CreatedAt   types.String `tfsdk:"created_at"`
}

func (d *CheckoutKeysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_checkout_keys"
}

func (d *CheckoutKeysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fetches the checkout keys of a project",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of this data source: project slug.",
				Computed:            true,
			},
			"project_slug": schema.StringAttribute{
				MarkdownDescription: "Project slug in the form `vcs-slug/org-name/repo-name`. The / characters may be URL-escaped.",
				Required:            true,
			},
			"keys": schema.ListNestedAttribute{
				MarkdownDescription: "List of checkout keys",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"public_key": schema.StringAttribute{
							MarkdownDescription: "A public SSH key",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of checkout key. This may be either `deploy-key` or `github-user-key`",
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
						},
					},
				},
			},
		},
	}
}

func (d *CheckoutKeysDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CheckoutKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CheckoutKeysDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	param := project.NewListProjectCheckoutKeysParamsWithContext(ctx).WithDefaults()
	param = param.WithProjectSlug(data.ProjectSlug.ValueString())

	res, err := d.client.Client.Project.ListProjectCheckoutKeys(param, d.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError("Encountered error fetching API", fmt.Sprintf("%s", err))
		return
	}

	info := res.GetPayload()
	// token := info.NextPageToken
	// TODO: consider support of pagination
	for _, ck := range info.Items {
		keyState := keyModel{
			PublicKey:   types.StringValue(ck.PublicKey),
			Fingerprint: types.StringValue(ck.Fingerprint),
			Type:        types.StringValue(ck.Type),
			Preferred:   types.BoolValue(*ck.Preferred),
			CreatedAt:   types.StringValue(ck.CreatedAt.String()),
		}
		data.Keys = append(data.Keys, keyState)
	}
	data.Id = data.ProjectSlug

	// Save data into Terraform state
	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
