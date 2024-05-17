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
var _ datasource.DataSource = &DevDataSource{}

func NewDevDataSource() datasource.DataSource {
	return &DevDataSource{}
}

// DevDataSource defines the data source implementation.
type DevDataSource struct {
	client *http.Client
}

// DevDataSourceModel describes the data source data model.
type DevDataSourceModel struct {
	Dev []DevTFModel `tfsdk:"dev"`
}

func (d *DevDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dev"
}

func (d *DevDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
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
		},
	}
}

func (d *DevDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DevDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiDev []DevAPIModel
	var state DevDataSourceModel

	// Make a call to your API to fetch the dev data
	httpResp, err := d.client.Get("http://localhost:8080/dev")
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read dev, got error: %s", err))
		return
	}
	defer httpResp.Body.Close()

	// Read the response body
	bodyBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read response body, got error: %s", err))
		return
	}

	// Print the response body
	fmt.Println(string(bodyBytes))

	// Decode the response body into the API model struct
	err = json.Unmarshal(bodyBytes, &apiDev)
	if err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Unable to decode response body, got error: %s", err))
		return
	}

	// Convert API model to Terraform schema model and set in state
	for _, apiDev := range apiDev {
		devState := DevTFModel{
			Name: types.StringValue(apiDev.Name),
			Id:   types.StringValue(apiDev.Id),
		}

		for _, apiEngineer := range apiDev.Engineers {
			devState.Engineers = append(devState.Engineers, EngineerTFModel{
				Name:  types.StringValue(apiEngineer.Name),
				Id:    types.StringValue(apiEngineer.Id),
				Email: types.StringValue(apiEngineer.Email),
			})
		}

		state.Dev = append(state.Dev, devState)

	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
