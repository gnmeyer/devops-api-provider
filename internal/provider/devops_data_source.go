// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DevopsDataSource{}

func NewDevopsDataSource() datasource.DataSource {
	return &DevopsDataSource{}
}

// DevopsDataSource defines the data source implementation.
type DevopsDataSource struct {
	client *http.Client
}

// DevopsDataSourceModel describes the data source data model.
type DevopsDataSourceModel struct {
	Devops []DevopsTFModel `tfsdk:"devops"`
}

func (d *DevopsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devops"
}

func (d *DevopsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"devops": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"dev": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Computed: true,
									},
									"id": schema.StringAttribute{
										Computed: true,
									},
									"engineers": schema.ListNestedAttribute{
										Computed: true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"name": schema.StringAttribute{
													Computed: true,
												},
												"id": schema.StringAttribute{
													Computed: true,
												},
												"email": schema.StringAttribute{
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"ops": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Computed: true,
									},
									"id": schema.StringAttribute{
										Computed: true,
									},
									"engineers": schema.ListNestedAttribute{
										Computed: true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"name": schema.StringAttribute{
													Computed: true,
												},
												"id": schema.StringAttribute{
													Computed: true,
												},
												"email": schema.StringAttribute{
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *DevopsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *DevopsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiDevops []DevopsAPIModel
	var state DevopsDataSourceModel

	// Make a call to your API to fetch the devops data
	httpResp, err := d.client.Get("http://localhost:8080/devops")
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read devops, got error: %s", err))
		return
	}
	defer httpResp.Body.Close()

	// Read the response body
	bodyBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read response body, got error: %s", err))
		return
	}

	// Decode the response body into the API model struct
	err = json.Unmarshal(bodyBytes, &apiDevops)
	if err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Unable to decode response body, got error: %s", err))
		return
	}

	// Convert API model to Terraform schema model and set in state
	for _, apiDevopsItem := range apiDevops {
		devopsState := DevopsTFModel{
			Id:  types.StringValue(apiDevopsItem.Id),
			Dev: convertDevAPIModelToTFModel(apiDevopsItem.Dev),
			Ops: convertOpsAPIModelToTFModel(apiDevopsItem.Ops),
		}

		state.Devops = append(state.Devops, devopsState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Helper function to convert API model to Terraform schema model
func convertDevAPIModelToTFModel(apiModel []DevAPIModel) []DevTFModel {
	var tfModel []DevTFModel
	for _, apiItem := range apiModel {
		tfItem := DevTFModel{
			Name:      types.StringValue(apiItem.Name),
			Id:        types.StringValue(apiItem.Id),
			Engineers: convertEngineerAPIModelToTFModel(apiItem.Engineers),
		}
		tfModel = append(tfModel, tfItem)
	}
	return tfModel
}

// Helper function to convert API model to Terraform schema model
func convertOpsAPIModelToTFModel(apiModel []OpsAPIModel) []OpsTFModel {
	var tfModel []OpsTFModel
	for _, apiItem := range apiModel {
		tfItem := OpsTFModel{
			Name:      types.StringValue(apiItem.Name),
			Id:        types.StringValue(apiItem.Id),
			Engineers: convertEngineerAPIModelToTFModel(apiItem.Engineers),
		}
		tfModel = append(tfModel, tfItem)
	}
	return tfModel
}

// Helper function to convert Engineer API model to Terraform schema model
func convertEngineerAPIModelToTFModel(apiModel []EngineerAPIModel) []EngineerTFModel {
	var tfModel []EngineerTFModel
	for _, apiItem := range apiModel {
		tfItem := EngineerTFModel{
			Name:  types.StringValue(apiItem.Name),
			Id:    types.StringValue(apiItem.Id),
			Email: types.StringValue(apiItem.Email),
		}
		tfModel = append(tfModel, tfItem)
	}
	return tfModel
}
