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

// Engineer is used for Terraform schema.
type OpsTFModel struct {
	Name      types.String      `tfsdk:"name"`
	Id        types.String      `tfsdk:"id"`
	Engineers []EngineerTFModel `tfsdk:"engineers"`
}

type OpsAPIModel struct {
	Name      string             `json:"name"`
	Id        string             `json:"id"`
	Engineers []EngineerAPIModel `json:"engineers"`
}

type DevTFModel struct {
	Name      types.String      `tfsdk:"name"`
	Id        types.String      `tfsdk:"id"`
	Engineers []EngineerTFModel `tfsdk:"engineers"`
}

type DevAPIModel struct {
	Name      string             `json:"name"`
	Id        string             `json:"id"`
	Engineers []EngineerAPIModel `json:"engineers"`
}

type DevopsTFModel struct {
	Id  types.String `tfsdk:"id"`
	Dev []DevTFModel `tfsdk:"dev"`
	Ops []OpsTFModel `tfsdk:"ops"`
}

type DevopsAPIModel struct {
	Id  string        `tfsdk:"id"`
	Dev []DevAPIModel `tfsdk:"dev"`
	Ops []OpsAPIModel `tfsdk:"ops"`
}
