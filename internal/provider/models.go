package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// Engineer is used for Terraform schema.
type EngineerTFModel struct {
	Name  types.String `tfsdk:"name"`
	Id    types.String `tfsdk:"id"`
	Email types.String `tfsdk:"email"`
}

// EngineerAPIModel is used to unmarshal JSON data from the API.
type EngineerAPIModel struct {
	Name  string `json:"name"`
	Id    string `json:"id"`
	Email string `json:"email"`
}
