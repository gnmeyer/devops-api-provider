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
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &EngineerResource{}
var _ resource.ResourceWithImportState = &EngineerResource{}

func NewEngineerResource() resource.Resource {
	return &EngineerResource{}
}

// EngineerResource defines the resource implementation.
type EngineerResource struct {
	client *http.Client
}

func (r *EngineerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engineer-resource"
}

func (r *EngineerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Engineer resource",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Engineer name",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Engineer id",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Engineer email",
			},
		},
	}
}

func (r *EngineerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EngineerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EngineerTFModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var engineerObject EngineerAPIModel
	engineerObject.Name = data.Name.ValueString()
	engineerObject.Id = data.Id.ValueString()
	engineerObject.Email = data.Email.ValueString()

	// Convert data to JSON
	jsonData, err := json.Marshal(engineerObject)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
	} else {
		log.Printf("Marshalled JSON: %s", string(jsonData))
	}

	// Make a POST request with JSON data
	httpResp, err := r.client.Post("http://localhost:8080/engineers", "application/json", bytes.NewBuffer(jsonData))
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
	log.Printf("Response Body: %s", string(bodyBytes))

	// Unmarshal the response into an Engineer struct
	var engineerRespObject EngineerAPIModel
	err = json.Unmarshal(bodyBytes, &engineerRespObject)
	if err != nil {
		resp.Diagnostics.AddError("JSON Error", fmt.Sprintf("Unable to unmarshal response body, got error: %s", err))
		return
	}

	// Map response body to schema and populate Computed attribute values
	data.Name = types.StringValue(engineerRespObject.Name)
	data.Id = types.StringValue(engineerRespObject.Id)
	data.Email = types.StringValue(engineerRespObject.Email)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EngineerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EngineerTFModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Make a call to your API to fetch the engineer data by ID
	httpResp, err := r.client.Get("http://localhost:8080/engineers/id/" + data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read engineer, got error: %s", err))
		return
	}
	defer httpResp.Body.Close()

	// Read the HTTP response body into a byte slice
	bodyBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read response body, got error: %s", err))
		return
	}

	// Unmarshal the response into an Engineer struct
	var engineerRespObject EngineerAPIModel
	err = json.Unmarshal(bodyBytes, &engineerRespObject)
	if err != nil {
		resp.Diagnostics.AddError("JSON Error", fmt.Sprintf("Unable to unmarshal response body, got error: %s", err))
		return
	}

	// Update the data object with the response data
	data.Name = types.StringValue(engineerRespObject.Name)
	data.Email = types.StringValue(engineerRespObject.Email)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EngineerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data EngineerTFModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EngineerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EngineerTFModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *EngineerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
