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
var _ datasource.DataSource = &OpsDataSource{}

func NewOpsDataSource() datasource.DataSource {
	return &OpsDataSource{}
}

// OpsDataSource defines the data source implementation.
type OpsDataSource struct {
	client *http.Client
}

// OpsDataSourceModel describes the data source data model.
type OpsDataSourceModel struct {
	Ops []OpsTFModel `tfsdk:"ops"`
}

func (d *OpsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ops"
}

func (d *OpsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
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
	}
}

func (d *OpsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OpsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var apiOps []OpsAPIModel
	var state OpsDataSourceModel

	// Make a call to your API to fetch the ops data
	httpResp, err := d.client.Get("http://localhost:8080/op")
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ops, got error: %s", err))
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
	err = json.Unmarshal(bodyBytes, &apiOps)
	if err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Unable to decode response body, got error: %s", err))
		return
	}

	// Convert API model to Terraform schema model and set in state
	for _, apiOp := range apiOps {
		opState := OpsTFModel{
			Name: types.StringValue(apiOp.Name),
			Id:   types.StringValue(apiOp.Id),
		}

		for _, apiEngineer := range apiOp.Engineers {
			opState.Engineers = append(opState.Engineers, EngineerTFModel{
				Name:  types.StringValue(apiEngineer.Name),
				Id:    types.StringValue(apiEngineer.Id),
				Email: types.StringValue(apiEngineer.Email),
			})
		}

		state.Ops = append(state.Ops, opState)

	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
