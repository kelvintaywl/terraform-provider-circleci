package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/kelvintaywl/circleci-runner-go-sdk/client/token"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &RunnerTokensDataSource{}

func NewRunnerTokensDataSource() datasource.DataSource {
	return &RunnerTokensDataSource{}
}

type RunnerTokensDataSource struct {
	client *CircleciAPIClient
}

type RunnerTokensDataSourceModel struct {
	ResourceClass types.String `tfsdk:"resource_class"`
	Tokens        []tokenModel `tfsdk:"tokens"`
	Id            types.String `tfsdk:"id"`
}

type tokenModel struct {
	Id            types.String `tfsdk:"id"`
	ResourceClass types.String `tfsdk:"resource_class"`
	Nickname      types.String `tfsdk:"nickname"`
	CreatedAt     types.String `tfsdk:"created_at"`
}

func (d *RunnerTokensDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_runner_tokens"
}

func (d *RunnerTokensDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fetches the list of tokens for a specific Runner resource-class.",
		Attributes: map[string]schema.Attribute{
			"resource_class": schema.StringAttribute{
				MarkdownDescription: "The Runner resource-class name.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of this data source: Runner resource-class.",
				Computed:            true,
			},
			"tokens": schema.ListNestedAttribute{
				MarkdownDescription: "List of tokens",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique ID of the Runner token.",
							Computed:            true,
						},
						"resource_class": schema.StringAttribute{
							MarkdownDescription: "The Runner resource-class name.",
							Computed:            true,
						},
						"nickname": schema.StringAttribute{
							MarkdownDescription: "The Runner token alias.",
							Computed:            true,
						},
						"created_at": schema.StringAttribute{
							MarkdownDescription: "Date and time the token was created",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *RunnerTokensDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RunnerTokensDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RunnerTokensDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	param := token.NewListTokensParamsWithContext(ctx).WithDefaults()
	param = param.WithResourceClass(data.ResourceClass.ValueString())

	res, err := d.client.RunnerClient.Token.ListTokens(param, d.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError("Encountered error fetching API", fmt.Sprintf("%s", err))
		return
	}

	info := res.GetPayload()
	for _, tk := range info.Items {
		tokenState := tokenModel{
			Id:            types.StringValue(tk.ID.String()),
			ResourceClass: types.StringValue(*tk.ResourceClass),
			Nickname:      types.StringValue(*tk.Nickname),
			CreatedAt:     types.StringValue(tk.CreatedAt.String()),
		}
		data.Tokens = append(data.Tokens, tokenState)
	}

	data.Id = data.ResourceClass

	// Save data into Terraform state
	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
