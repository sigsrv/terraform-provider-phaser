// Copyright (c) EcmaXp.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sigsrv/terraform-provider-phaser/internal/phaser"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SequentialResource{}
var _ resource.ResourceWithImportState = &SequentialResource{}
var _ resource.ResourceWithModifyPlan = &SequentialResource{}

func NewSequentialResource() resource.Resource {
	return &SequentialResource{}
}

type SequentialResource struct{}

type SequentialResourceModel struct {
	Phase  types.String `tfsdk:"phase"`
	Phases types.List   `tfsdk:"phases"`
}

func (r *SequentialResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sequential"
}

func (r *SequentialResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sequential resource",
		Attributes: map[string]schema.Attribute{
			"phase": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Current active phase",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"phases": schema.ListAttribute{
				MarkdownDescription: "List of sequential phases to progress through",
				Required:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.NoNullValues(),
					listvalidator.UniqueValues(),
				},
			},
		},
	}
}

func (r *SequentialResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
}

func (r *SequentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SequentialResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SequentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *SequentialResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var config, data *SequentialResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() || config == nil {
		return
	}

	var phases []string
	resp.Diagnostics.Append(config.Phases.ElementsAs(ctx, &phases, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Phase.IsUnknown() {
		data.Phase = types.StringValue(phases[0])
	} else {
		phase := data.Phase.ValueString()
		nextPhase, err := phaser.GetNextPhaseSequential(phases, phase)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid phase",
				err.Error(),
			)
			return
		}

		data.Phase = types.StringValue(nextPhase)
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &data)...)
}

func (r *SequentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SequentialResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SequentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *SequentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data SequentialResourceModel

	data.Phase = types.StringValue(req.ID)
	data.Phases = types.ListNull(types.StringType)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
