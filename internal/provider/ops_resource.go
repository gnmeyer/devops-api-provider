// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &OpsResource{}
var _ resource.ResourceWithImportState = &OpsResource{}

func NewOpsResource() resource.Resource {
	return &OpsResource{}
}

// OpsResource defines the resource implementation.
type OpsResource struct {
	client *http.Client
}

func (r *OpsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ops-resource"
}

func (r *OpsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Op resource",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"engineers": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"id": schema.StringAttribute{
							Required: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"email": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (r *OpsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *OpsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OpsTFModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	/* Step 1 Create Engineers */

	var apiEngineers []EngineerAPIModel

	tflog.Debug(ctx, "Begin checking engineers", map[string]interface{}{})

	// Convert engineers to API model and make a POST request for each engineer
	for _, engineer := range data.Engineers {

		tflog.Debug(ctx, "Checking Engineer", map[string]interface{}{"engineerID": engineer.Id})

		// Check if the engineer already exists
		httpEngineerResp, err := r.client.Get("http://localhost:8080/engineers/id/" + engineer.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get engineer, got error: %s", err))
			continue
		}

		// Read the HTTP response body
		bodyEngineerBytes, err := ioutil.ReadAll(httpEngineerResp.Body)
		if err != nil {
			resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read response body, got error: %s", err))
			continue
		}

		// Unmarshal the response into an Engineer struct
		var engineerRespObject EngineerAPIModel
		err = json.Unmarshal(bodyEngineerBytes, &engineerRespObject)
		if err != nil {
			resp.Diagnostics.AddError("JSON Error", fmt.Sprintf("Unable to unmarshal response body, got error: %s", err))
			continue
		}

		// If the engineer exists, add it to apiEngineers and skip to the next engineer
		if httpEngineerResp.StatusCode == http.StatusOK {
			tflog.Debug(ctx, "Engineer exists, moving on", map[string]interface{}{"engineerID": engineer.Id})
			apiEngineers = append(apiEngineers, engineerRespObject)
			continue
		}

		tflog.Debug(ctx, "Engineer does not exist", map[string]interface{}{"engineerID": engineer.Id})

		return
	}

	/* Step 2 Create Ops */

	var OpsObject OpsAPIModel
	OpsObject.Name = data.Name.ValueString()
	OpsObject.Id = data.Id.ValueString()
	OpsObject.Engineers = apiEngineers

	// Convert data to JSON
	jsonData, err := json.Marshal(OpsObject)
	if err != nil {
		tflog.Debug(ctx, "Error marshalling Engineer JSON", map[string]interface{}{"error": err.Error()})
		return
	} else {
		tflog.Debug(ctx, "Marshalled Engineer JSON", map[string]interface{}{"engineerData": string(jsonData)})
	}

	// Make a POST request with JSON data
	httpResp, err := r.client.Post("http://localhost:8080/op", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
		return
	}

	// Read the HTTP response body
	bodyBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read response body, got error: %s", err))
		return
	}

	// Log the response body
	tflog.Debug(ctx, "Response Body", map[string]interface{}{"body": string(bodyBytes)})

	// Unmarshal the response into an Ops struct
	var OpsRespObject OpsAPIModel
	err = json.Unmarshal(bodyBytes, &OpsRespObject)
	if err != nil {
		resp.Diagnostics.AddError("JSON Error", fmt.Sprintf("Unable to unmarshal response body, got error: %s", err))
		return
	}

	// Clear existing engineers before appending new ones
	data.Engineers = []EngineerTFModel{} // Reinitialize the slice

	for _, tfEngineer := range OpsRespObject.Engineers {
		data.Engineers = append(data.Engineers, EngineerTFModel{
			Name:  types.StringValue(tfEngineer.Name),
			Id:    types.StringValue(tfEngineer.Id),
			Email: types.StringValue(tfEngineer.Email),
		})
	}

	data.Name = types.StringValue(OpsRespObject.Name)
	data.Id = types.StringValue(OpsRespObject.Id)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OpsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OpsTFModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Make a call to your API to fetch the Ops data by ID
	httpResp, err := r.client.Get("http://localhost:8080/op/id/" + data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Ops, got error: %s", err))
		return
	}
	defer httpResp.Body.Close()

	// Read the HTTP response body into a byte slice
	bodyBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read response body, got error: %s", err))
		return
	}

	// Unmarshal the response into an Ops struct
	var OpsRespObject OpsAPIModel
	err = json.Unmarshal(bodyBytes, &OpsRespObject)
	if err != nil {
		resp.Diagnostics.AddError("JSON Error", fmt.Sprintf("Unable to unmarshal response body, got error: %s", err))
		return
	}

	// // Update the data object with the response data
	data.Name = types.StringValue(OpsRespObject.Name)
	data.Id = types.StringValue(OpsRespObject.Id)

	for _, tfEngineer := range OpsRespObject.Engineers {
		data.Engineers = append(data.Engineers, EngineerTFModel{
			Name:  types.StringValue(tfEngineer.Name),
			Id:    types.StringValue(tfEngineer.Id),
			Email: types.StringValue(tfEngineer.Email),
		})
	}

	// // Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OpsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data OpsTFModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Log the entire data object
	log.Printf("Data: %+v", data)

	if resp.Diagnostics.HasError() {
		return
	}

	log.Printf("ID: %s", data.Id.ValueString())

	var apiEngineers []EngineerAPIModel

	// Convert engineers to API model and make a POST request for each engineer
	for _, engineer := range data.Engineers {

		tflog.Debug(ctx, "Checking Engineer", map[string]interface{}{"engineerID": engineer.Id})

		// Check if the engineer already exists
		httpEngineerResp, err := r.client.Get("http://localhost:8080/engineers/id/" + engineer.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get engineer, got error: %s", err))
			continue
		}

		// Read the HTTP response body
		bodyEngineerBytes, err := ioutil.ReadAll(httpEngineerResp.Body)
		if err != nil {
			resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read response body, got error: %s", err))
			continue
		}

		// Unmarshal the response into an Engineer struct
		var engineerRespObject EngineerAPIModel
		err = json.Unmarshal(bodyEngineerBytes, &engineerRespObject)
		if err != nil {
			resp.Diagnostics.AddError("JSON Error", fmt.Sprintf("Unable to unmarshal response body, got error: %s", err))
			continue
		}

		// If the engineer exists, add it to apiEngineers and skip to the next engineer
		if httpEngineerResp.StatusCode == http.StatusOK {
			tflog.Debug(ctx, "Engineer exists, moving on", map[string]interface{}{"engineerID": engineer.Id})
			apiEngineers = append(apiEngineers, engineerRespObject)
			continue
		}

		tflog.Debug(ctx, "Engineer does not exist", map[string]interface{}{"engineerID": engineer.Id})

		return
	}

	/* Step 2 Create Ops */

	var OpsObject OpsAPIModel
	OpsObject.Name = data.Name.ValueString()
	OpsObject.Id = data.Id.ValueString()
	OpsObject.Engineers = apiEngineers

	// Convert data to JSON
	jsonData, err := json.Marshal(OpsObject)
	if err != nil {
		tflog.Debug(ctx, "Error marshalling Engineer JSON", map[string]interface{}{"error": err.Error()})
		return
	} else {
		tflog.Debug(ctx, "Marshalled Engineer JSON", map[string]interface{}{"engineerData": string(jsonData)})
	}

	// Create a new HTTP request
	newReq, err := http.NewRequest(http.MethodPut, "http://localhost:8080/op/"+OpsObject.Id, bytes.NewBuffer(jsonData))
	if err != nil {
		resp.Diagnostics.AddError("Request Error", fmt.Sprintf("Unable to create request, got error: %s", err))
		return
	}

	// Set the Content-Type header
	newReq.Header.Set("Content-Type", "application/json")

	// Send the request
	httpResp, err := r.client.Do(newReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
		return
	}
	// Read the HTTP response body
	bodyBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read response body, got error: %s", err))
		return
	}

	// Log the response body
	tflog.Debug(ctx, "Response Body", map[string]interface{}{"body": string(bodyBytes)})

	// Unmarshal the response into an Ops struct
	var OpsRespObject OpsAPIModel
	err = json.Unmarshal(bodyBytes, &OpsRespObject)
	if err != nil {
		resp.Diagnostics.AddError("JSON Error", fmt.Sprintf("Unable to unmarshal response body, got error: %s", err))
		return
	}

	// Clear existing engineers before appending new ones
	data.Engineers = []EngineerTFModel{} // Reinitialize the slice

	for _, tfEngineer := range OpsRespObject.Engineers {
		data.Engineers = append(data.Engineers, EngineerTFModel{
			Name:  types.StringValue(tfEngineer.Name),
			Id:    types.StringValue(tfEngineer.Id),
			Email: types.StringValue(tfEngineer.Email),
		})
	}

	data.Name = types.StringValue(OpsRespObject.Name)
	data.Id = types.StringValue(OpsRespObject.Id)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *OpsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OpsTFModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var OpsObject OpsAPIModel
	OpsObject.Name = ""
	OpsObject.Engineers = []EngineerAPIModel{}
	OpsObject.Id = data.Id.ValueString()

	// Convert data to JSON
	jsonData, err := json.Marshal(OpsObject)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
	} else {
		log.Printf("Marshalled JSON: %s", string(jsonData))
	}

	// Create a new HTTP request
	newReq, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/op/"+data.Id.ValueString(), bytes.NewBuffer(jsonData))
	if err != nil {
		resp.Diagnostics.AddError("Request Error", fmt.Sprintf("Unable to create request, got error: %s", err))
		return
	}

	// Set the Content-Type header
	newReq.Header.Set("Content-Type", "application/json")

	// Send the request
	httpResp, err := r.client.Do(newReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
		return
	}

	// Read the HTTP response body
	bodyBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read response body, got error: %s", err))
		return
	}

	// Log the response body
	log.Printf("Response Body: %s", string(bodyBytes))
}

func (r *OpsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
