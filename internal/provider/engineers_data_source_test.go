package provider

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/tfsdk"
    //"github.com/hashicorp/terraform-plugin-framework/types"

)

func TestEngineerDataSource(t *testing.T) {
    // Create a test server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`[
            {"name": "John Doe", "id": "1", "email": "john@example.com"},
            {"name": "Jane Doe", "id": "2", "email": "jane@example.com"}
        ]`))
    }))
    defer server.Close()

    // Initialize the data source with the test server's client
    ds := EngineerDataSource{
        client: server.Client(),
    }

    // Create a context and request/response objects
    ctx := context.Background()
    req := datasource.ReadRequest{
        Config: tfsdk.Config{
            Config: schema.Config{
                Schema: ds.Schema().Schema,
                Data:   nil,
            },
        },
    }
    resp := datasource.NewReadResponse()

    // Call the Read method
    ds.Read(ctx, req, &resp)

    // Check for errors in the response
    if resp.Diagnostics.HasError() {
        t.Errorf("unexpected error in Read: %s", resp.Diagnostics)
    }

    // Check the values
    if len(resp.State.Attributes) != 2 {
        t.Fatalf("expected 2 engineers, got %d", len(resp.State.Attributes))
    }

    // You can add more specific checks here for attributes values
}
