package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-openapi/strfmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/kelvintaywl/circleci-go-sdk/client/schedule"
	"github.com/kelvintaywl/circleci-go-sdk/models"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &ScheduleResource{}

func NewScheduleResource() resource.Resource {
	return &ScheduleResource{}
}

type ScheduleResource struct {
	client *CircleciAPIClient
}

type timetableModel struct {
	PerHour     types.Int64    `tfsdk:"per_hour"`
	HoursOfDay  []types.Int64  `tfsdk:"hours_of_day"`
	DaysOfWeek  []types.String `tfsdk:"days_of_week"`
	DaysOfMonth []types.Int64  `tfsdk:"days_of_month"`
	Months      []types.String `tfsdk:"months"`
}

type ScheduleResourceModel struct {
	ID                 types.String   `tfsdk:"id"`
	CreatedAt          types.String   `tfsdk:"created_at"`
	UpdatedAt          types.String   `tfsdk:"updated_at"`
	ProjectSlug        types.String   `tfsdk:"project_slug"`
	Name               types.String   `tfsdk:"name"`
	Description        types.String   `tfsdk:"description"`
	AttributionActor   types.String   `tfsdk:"actor"`
	PipelineParameters types.String   `tfsdk:"parameters"`
	Branch             types.String   `tfsdk:"branch"`
	Tag                types.String   `tfsdk:"tag"`
	Timetable          timetableModel `tfsdk:"timetable"`
}

func IsSystemActor(a *models.User) bool {
	login := *a.Login
	name := *a.Name

	switch {
	case login == "system-actor":
	case name == "Scheduled":
	case a.ID.String() == "d9b3fcaa-6032-405a-8c75-40079ce33c3e":
		return true
	}
	return false
}

var vAttributionActors []string = []string{
	"current",
	"system",
}

var vDaysOfWeek []string = []string{
	string(models.DayOfAWeekMON),
	string(models.DayOfAWeekTUE),
	string(models.DayOfAWeekWED),
	string(models.DayOfAWeekTHU),
	string(models.DayOfAWeekFRI),
	string(models.DayOfAWeekSAT),
	string(models.DayOfAWeekSUN),
}

var vMonths []string = []string{
	string(models.MonthJAN),
	string(models.MonthFEB),
	string(models.MonthMAR),
	string(models.MonthAPR),
	string(models.MonthMAY),
	string(models.MonthJUN),
	string(models.MonthJUL),
	string(models.MonthAUG),
	string(models.MonthSEP),
	string(models.MonthOCT),
	string(models.MonthNOV),
	string(models.MonthDEC),
}

func (r *ScheduleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schedule"
}

func (r *ScheduleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a project's schedule",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the schedule",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The date and time the schedule was created",
				Computed:            true,
				// unchanged even during updates
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The date and time the schedule was last updated",
				Computed:            true,
			},
			"project_slug": schema.StringAttribute{
				MarkdownDescription: "The project-slug for the schedule",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the schedule",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the schedule",
				Required:            true,
			},
			"actor": schema.StringAttribute{
				MarkdownDescription: "The actor to attribute as author of the scheduled pipeline (accepts 'current' or 'system')",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(vAttributionActors...),
				},
			},
			"branch": schema.StringAttribute{
				MarkdownDescription: "Branch name to trigger scheduled pipeline from (mutually exclusive to tag)",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("tag"),
					}...),
				},
			},
			"tag": schema.StringAttribute{
				MarkdownDescription: "Tag name to trigger scheduled pipeline from (mutually exclusive to branch)",
				Optional:            true,
			},
			"parameters": schema.StringAttribute{
				MarkdownDescription: "Pipeline parameters represented in a JSON string",
				Optional:            true,
			},
			"timetable": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"per_hour": schema.Int64Attribute{
						MarkdownDescription: "Number of times a schedule triggers per hour, value must be between 1 and 60",
						Required:            true,
						Validators: []validator.Int64{
							int64validator.Between(1, 60),
						},
					},
					"hours_of_day": schema.ListAttribute{
						ElementType:         types.Int64Type,
						MarkdownDescription: "Hours in a day in which the schedule triggers.",
						Required:            true,
						Validators: []validator.List{
							listvalidator.ValueInt64sAre(int64validator.Between(0, 23)),
							listvalidator.UniqueValues(),
						},
					},
					"days_of_week": schema.ListAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "Days in a week in which the schedule triggers.",
						Optional:            true,
						Validators: []validator.List{
							listvalidator.ValueStringsAre(stringvalidator.OneOf(vDaysOfWeek...)),
							listvalidator.UniqueValues(),
							listvalidator.ExactlyOneOf(path.Expressions{
								path.MatchRoot("timetable").AtName("days_of_month"),
							}...),
						},
					},
					"days_of_month": schema.ListAttribute{
						ElementType:         types.Int64Type,
						MarkdownDescription: "Days in a month in which the schedule triggers. This is mutually exclusive with days in a week.",
						Optional:            true,
						Validators: []validator.List{
							listvalidator.ValueInt64sAre(int64validator.Between(1, 31)),
							listvalidator.UniqueValues(),
						},
					},
					"months": schema.ListAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "Months in which the schedule triggers. Defaults to all months if not set.",
						Optional:            true,
						Validators: []validator.List{
							listvalidator.ValueStringsAre(stringvalidator.OneOf(vMonths...)),
							listvalidator.UniqueValues(),
						},
					},
				},
				MarkdownDescription: "Timetable that specifies when a schedule triggers.",
				Required:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *ScheduleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ScheduleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	param := schedule.NewGetScheduleParamsWithContext(ctx).WithDefaults()
	param = param.WithID(strfmt.UUID(id))

	res, err := r.client.Client.Schedule.GetSchedule(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Encountered error reading Schedule %s", id), fmt.Sprintf("%s", err))
		return
	}

	sc := res.GetPayload()

	state.ID = types.StringValue(sc.ID.String())
	state.CreatedAt = types.StringValue(sc.CreatedAt.String())
	state.UpdatedAt = types.StringValue(sc.UpdatedAt.String())
	state.ProjectSlug = types.StringValue(*sc.ProjectSlug)
	state.Name = types.StringValue(sc.Name)
	state.Description = types.StringValue(sc.Description)

	if IsSystemActor(sc.Actor) {
		state.AttributionActor = types.StringValue("system")
	} else {
		state.AttributionActor = types.StringValue("current")
	}

	if sc.Parameters.Branch != "" {
		state.Branch = types.StringValue(sc.Parameters.Branch)
	} else {
		state.Tag = types.StringValue(sc.Parameters.Tag)
	}
	jsn, err := json.Marshal(sc.Parameters.ScheduleBaseDataParameters)
	if err != nil {
		resp.Diagnostics.AddError("Encountered error marshalling parameters to JSON", fmt.Sprintf("%s", err))
		return
	}
	state.PipelineParameters = types.StringValue(string(jsn[:]))

	state.Timetable = timetableModel{
		PerHour: types.Int64Value(*sc.Timetable.PerHour),
	}
	for _, h := range sc.Timetable.HoursOfDay {
		state.Timetable.HoursOfDay = append(state.Timetable.HoursOfDay, types.Int64Value(int64(*h)))
	}
	for _, dw := range sc.Timetable.DaysOfWeek {
		state.Timetable.DaysOfWeek = append(state.Timetable.DaysOfWeek, types.StringValue(string(dw)))
	}
	for _, dm := range sc.Timetable.DaysOfMonth {
		state.Timetable.DaysOfMonth = append(state.Timetable.DaysOfMonth, types.Int64Value(int64(dm)))
	}
	for _, m := range sc.Timetable.Months {
		state.Timetable.Months = append(state.Timetable.Months, types.StringValue(string(m)))
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func makeUpsertBodyPayload(plan ScheduleResourceModel) (models.SchedulePayload, string, error) {
	name := plan.Name.ValueString()
	desc := plan.Description.ValueString()
	actor := plan.AttributionActor.ValueString()

	// timetable
	perhour := plan.Timetable.PerHour.ValueInt64()
	tt := models.ScheduleBaseDataTimetable{
		PerHour: &perhour,
	}
	for _, h := range plan.Timetable.HoursOfDay {
		hourOfDay := models.HourOfADay(h.ValueInt64())
		tt.HoursOfDay = append(tt.HoursOfDay, &hourOfDay)
	}

	// optional attributes
	if len(plan.Timetable.Months) > 0 {
		for _, m := range plan.Timetable.Months {
			month := models.Month(m.ValueString())
			tt.Months = append(tt.Months, month)
		}
	} else {
		tt.Months = make([]models.Month, 0)
	}
	if len(plan.Timetable.DaysOfWeek) > 0 {
		for _, dw := range plan.Timetable.DaysOfWeek {
			dayOfWeek := models.DayOfAWeek(dw.ValueString())
			tt.DaysOfWeek = append(tt.DaysOfWeek, dayOfWeek)
		}
	} else {
		tt.DaysOfWeek = make([]models.DayOfAWeek, 0)
	}

	if len(plan.Timetable.DaysOfMonth) > 0 {
		for _, dm := range plan.Timetable.DaysOfMonth {
			dayOfMonth := models.DayOfAMonth(dm.ValueInt64())
			tt.DaysOfMonth = append(tt.DaysOfMonth, dayOfMonth)
		}
	} else {
		tt.DaysOfMonth = make([]models.DayOfAMonth, 0)
	}

	// pipeline params
	branch := plan.Branch.ValueString()
	tag := plan.Tag.ValueString()
	blob := []byte(plan.PipelineParameters.ValueString())
	p := make(map[string]interface{})
	if len(blob) > 0 {
		err := json.Unmarshal(blob, &p)
		if err != nil {
			return models.SchedulePayload{}, "Encountered error unmarshalling parameters from JSON", err
		}
	}
	pipelineParams := models.ScheduleBaseDataParameters{
		Branch:                     branch,
		Tag:                        tag,
		ScheduleBaseDataParameters: p,
	}
	payload := models.SchedulePayload{
		AttributionActor: actor,
		ScheduleBaseData: models.ScheduleBaseData{
			Description: desc,
			Name:        name,
			Parameters:  &pipelineParams,
			Timetable:   &tt,
		},
	}
	return payload, "", nil
}

// Create creates the resource and sets the initial Terraform state.
func (r *ScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ScheduleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	param := schedule.NewAddScheduleParamsWithContext(ctx).WithDefaults()
	project := plan.ProjectSlug.ValueString()
	param = param.WithProjectSlug(project)

	payload, errStr, err := makeUpsertBodyPayload(plan)
	if err != nil {
		resp.Diagnostics.AddError(errStr, fmt.Sprintf("%s", err))
		return
	}

	param = param.WithBody(&payload)
	res, err := r.client.Client.Schedule.AddSchedule(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating schedule",
			fmt.Sprintf("Could not create schedule, unexpected error: %s", err.Error()),
		)
		return
	}
	sc := res.GetPayload()

	id := sc.ID.String()
	plan.ID = types.StringValue(id)
	plan.CreatedAt = types.StringValue(sc.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(sc.UpdatedAt.String())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan ScheduleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	param := schedule.NewUpdateScheduleParamsWithContext(ctx).WithDefaults()
	param = param.WithID(strfmt.UUID(id))

	payload, errStr, err := makeUpsertBodyPayload(plan)
	if err != nil {
		resp.Diagnostics.AddError(errStr, fmt.Sprintf("%s", err))
		return
	}

	param = param.WithBody(&payload)
	res, err := r.client.Client.Schedule.UpdateSchedule(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating schedule",
			fmt.Sprintf("Could not update schedule, unexpected error: %s", err.Error()),
		)
		return
	}

	sc := res.GetPayload()
	plan.CreatedAt = types.StringValue(sc.CreatedAt.String())
	plan.UpdatedAt = types.StringValue(sc.UpdatedAt.String())
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ScheduleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	param := schedule.NewDeleteScheduleParamsWithContext(ctx).WithDefaults()
	param = param.WithID(strfmt.UUID(id))

	_, err := r.client.Client.Schedule.DeleteSchedule(param, r.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting schedule",
			fmt.Sprintf("Could not delete schedule %s, unexpected error: %s", id, err.Error()),
		)
		return
	}
}
