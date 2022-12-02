package provider

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/runtime"
	rtc "github.com/go-openapi/runtime/client"

	api "github.com/kelvintaywl/circleci-schedule-go-sdk/client"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &ScheduleResource{}
var _ resource.ResourceWithImportState = &ScheduleResource{}

var DayOfWeek = []string{
	"MON",
	"TUE",
	"WED",
	"THU",
	"FRI",
	"SAT",
	"SUN",
}
var Month = []string{
	"JAN",
	"FEB",
	"MAR",
	"APR",
	"MAY",
	"JUN",
	"JUL",
	"AUG",
	"SEP",
	"OCT",
	"NOV",
	"DEC",
}
var AttributionActor = []string{
	"system",
	"current",
}

func NewScheduleResource() resource.Resource {
	return &ScheduleResource{}
}

// ScheduleResource defines the resource implementation.
type ScheduleResource struct {
	client *api.Circleci
	auth runtime.ClientAuthInfoWriter
}

// ScheduleResourceModel describes the resource data model.
type ScheduleResourceModel struct {
	Id               types.String `tfsdk:"id"`
	CreatedAt        types.String `tfsdk:"created_at"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	ProjectSlug      types.String `tfsdk:"project_slug"`
	AttributionActor types.String `tfsdk:"attribution_actor"`
	Actor            types.Object `tfsdk:"actor"`
	Parameters       types.Object `tfsdk:"parameters"`
	Timetable        types.Object `tfsdk:"timetable"`
}

func (r *ScheduleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schedule"
}

func (r *ScheduleResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	s := tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "CircleCI schedule resource",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				// read-only
				Computed:            true,
				Optional:            false,
				MarkdownDescription: "The unique ID of the schedule",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"created_at": {
				// read-only
				Computed:            true,
				Optional:            false,
				MarkdownDescription: "The date and time the schedule was created",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"updated_at": {
				// read-only
				Computed:            true,
				Optional:            false,
				MarkdownDescription: "The date and time the schedule was last updated",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"actor": {
				MarkdownDescription: "The actor for the schedule",
				// read-only
				Computed: true,
				Optional: false,
				Attributes: tfsdk.SingleNestedAttributes(
					map[string]tfsdk.Attribute{
						"id": {
							Type:                types.StringType,
							MarkdownDescription: "The user ID of the actor",
						},
						"login": {
							Type:                types.StringType,
							MarkdownDescription: "The user login name of the actor",
						},
						"name": {
							Type:                types.StringType,
							MarkdownDescription: "The name of the actor",
						},
					},
				),
			},
			"name": {
				MarkdownDescription: "Name of the schedule",
				Required:            true,
				Type:                types.StringType,
			},
			"description": {
				MarkdownDescription: "Description of the schedule",
				Required:            true,
				Type:                types.StringType,
			},
			"project_slug": {
				MarkdownDescription: "The project-slug for the schedule",
				Required:            true,
				Type:                types.StringType,
			},
			"parameters": {
				MarkdownDescription: "Pipeline parameters for the schedule",
				Required:            true,
				Attributes: tfsdk.SingleNestedAttributes(
					map[string]tfsdk.Attribute{
						"branch": {
							Type:                types.StringType,
							Optional:            true,
							MarkdownDescription: "Branch name for schedule to trigger on",
						},
						"tag": {
							Type:                types.StringType,
							Optional:            true,
							MarkdownDescription: "Git tag for schedule to trigger on",
						},
						"json": {
							Type:                types.StringType,
							Optional:            true,
							MarkdownDescription: "Additional parameters for pipeline (in JSON string)",
						},
					},
				),
			},
			"timetable": {
				MarkdownDescription: "The timetable for the schedule",
				Required:            true,
				Attributes: tfsdk.SingleNestedAttributes(
					map[string]tfsdk.Attribute{
						"per_hour": {
							Type:                types.Int64Type,
							Required:            true,
							MarkdownDescription: "Repeats Per Hour",
							Validators: []tfsdk.AttributeValidator{
								int64validator.Between(1, 12),
							},
						},
						"days_of_week": {
							Type:                types.ListType{ElemType: types.StringType},
							Optional:            true,
							MarkdownDescription: "Repeats on these days of the week",
							Validators: []tfsdk.AttributeValidator{
								listvalidator.ValuesAre(
									stringvalidator.OneOf(DayOfWeek...),
								),
							},
						},
						"days_of_month": {
							Type:                types.ListType{ElemType: types.Int64Type},
							Optional:            true,
							MarkdownDescription: "Repeats on these days of the month",
							Validators: []tfsdk.AttributeValidator{
								listvalidator.ValuesAre(
									int64validator.Between(1, 31),
								),
							},
						},
						"months": {
							Type:                types.ListType{ElemType: types.Int64Type},
							Optional:            true,
							MarkdownDescription: "Repeats on these months",
							Validators: []tfsdk.AttributeValidator{
								listvalidator.ValuesAre(
									stringvalidator.OneOf(Month...),
								),
							},
						},
					},
				),
			},
			"attribution_actor": {
				MarkdownDescription: "The attribution-actor of the scheduled pipeline",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf(AttributionActor...),
				},
			},
		},
	}
	return s, nil
}

func (r *ScheduleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	// cfg := api.DefaultTransportConfig().WithHost(req.ProviderData.Hostname)
	cfg := api.DefaultTransportConfig().WithHost("req.ProviderData.Hostname")
	client := api.NewHTTPClientWithConfig(strfmt.Default, cfg)
	// auth := rtc.APIKeyAuth("Circle-Token", "header", req.ProviderData.ApiToken)
	auth := rtc.APIKeyAuth("Circle-Token", "header", "req.ProviderData.ApiToken")

	r.client = client
	r.auth = auth
}

func (r *ScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ScheduleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	//     return
	// }

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.StringValue("example-id")

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ScheduleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Prevent panic if the provider has not been configured.
	if r.client == nil {
		resp.Diagnostics.AddError(
			"Unconfigured CircleCI API Client",
			"Expected configured CircleCI API client. Please report this issue to the provider developers.",
		)

		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ScheduleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ScheduleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *ScheduleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
