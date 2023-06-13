package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/kelvintaywl/circleci-runner-go-sdk/client/resource_class"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &RunnerResourceClassesDataSource{}

func NewRunnerResourceClassesDataSource() datasource.DataSource {
	return &RunnerResourceClassesDataSource{}
}

type RunnerResourceClassesDataSource struct {
	client *CircleciAPIClient
}

type RunnerResourceClassesDataSourceModel struct {
	Namespace       types.String         `tfsdk:"namespace"`
	ResourceClasses []resourceClassModel `tfsdk:"resource_classes"`
	Id              types.String         `tfsdk:"id"`
}

type resourceClassModel struct {
	Id            types.String `tfsdk:"id"`
	ResourceClass types.String `tfsdk:"resource_class"`
	Description   types.String `tfsdk:"description"`
}

func (d *RunnerResourceClassesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_runner_resource_classes"
}

func (d *RunnerResourceClassesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fetches the list of Runner resource-classes for a specific namespace.",
		Attributes: map[string]schema.Attribute{
			"namespace": schema.StringAttribute{
				MarkdownDescription: "CircleCI namespace.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of this data source: namespace.",
				Computed:            true,
			},
			"resource_classes": schema.ListNestedAttribute{
				MarkdownDescription: "List of resource-classes",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique ID of the Runner resource-class",
							Computed:            true,
						},
						"resource_class": schema.StringAttribute{
							MarkdownDescription: "The Runner resource-class name",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description of the Runner resource-class",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *RunnerResourceClassesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RunnerResourceClassesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RunnerResourceClassesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	param := resource_class.NewListResourceClassesParamsWithContext(ctx).WithDefaults()
	param = param.WithNamespace(data.Namespace.ValueString())

	res, err := d.client.RunnerClient.ResourceClass.ListResourceClasses(param, d.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError("Encountered error fetching API", fmt.Sprintf("%s", err))
		return
	}

	info := res.GetPayload()
	for _, rc := range info.Items {
		resourceClassState := resourceClassModel{
			Id:            types.StringValue(rc.ID.String()),
			ResourceClass: types.StringValue(*rc.ResourceClass),
			Description:   types.StringValue(*rc.Description),
		}
		data.ResourceClasses = append(data.ResourceClasses, resourceClassState)
	}
	data.Id = data.Namespace

	// Save data into Terraform state
	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
