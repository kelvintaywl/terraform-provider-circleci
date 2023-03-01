package provider

import (
	"context"
	"fmt"

	"github.com/go-openapi/strfmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/kelvintaywl/circleci-go-sdk/client/contexts"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &ContextDataSource{}

func NewContextDataSource() datasource.DataSource {
	return &ContextDataSource{}
}

type ContextDataSource struct {
	client *CircleciAPIClient
}

// ContextDataSourceModel describes the data source data model.
type ContextDataSourceModel struct {
	Name      types.String `tfsdk:"name"`
	Owner     ownerModel   `tfsdk:"owner"`
	Id        types.String `tfsdk:"id"`
	CreatedAt types.String `tfsdk:"created_at"`
}

func (d *ContextDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_context"
}

func (d *ContextDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fetches the context",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "context name.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of this context",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The date and time the context was created",
				Computed:            true,
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

func (d *ContextDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ContextDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ContextDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	nextToken := ""
	ownerId := strfmt.UUID(data.Owner.Id.ValueString())
	ownerType := data.Owner.Type.ValueString()
	name := data.Name.ValueString()

	for {
		param := contexts.NewListContextsParamsWithContext(ctx).WithDefaults()
		param = param.WithOwnerID(&ownerId).WithOwnerType(ownerType).WithPageToken(&nextToken)

		res, err := d.client.Client.Contexts.ListContexts(param, d.client.Auth)
		if err != nil {
			resp.Diagnostics.AddError("Encountered error fetching API", fmt.Sprintf("%s", err))
			return
		}

		info := res.GetPayload()
		nextToken = info.NextPageToken
		for _, c := range info.Items {
			if c.Name == name {
				id := c.ID.String()
				data.Id = types.StringValue(id)
				createdAt := c.CreatedAt.String()
				data.CreatedAt = types.StringValue(createdAt)

				// Save data into Terraform state
				diags := resp.State.Set(ctx, &data)
				resp.Diagnostics.Append(diags...)
				if resp.Diagnostics.HasError() {
					return
				}
				return
			}
		}
		if nextToken == "" && data.Id.ValueString() == "" {
			resp.Diagnostics.AddError(fmt.Sprintf("Did not find context with name %s", name), fmt.Sprintf("%s", err))
			return
		}
	}
}
