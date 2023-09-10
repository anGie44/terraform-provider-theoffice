// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/anGie44/terraform-provider-theoffice/internal/theoffice"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure theOfficeProvider satisfies various provider interfaces.
var _ provider.Provider = &theOfficeProvider{}

// theOfficeProvider defines the provider implementation.
type theOfficeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// theOfficeProviderModel describes the provider data model.
type theOfficeProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *theOfficeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "theoffice"
	resp.Version = p.version
}

func (p *theOfficeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with theOffice API",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description: "The REST API endpoint to use for reading data (default: http://theofficeapi-angelinepinilla.b4a.run)",
				Optional:    true,
			},
		},
	}
}

func (p *theOfficeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data theOfficeProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	if data.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown theOffice API endpoint",
			"The provider cannot create theOffice API client as there is an unknown configuration value for theOffice API endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the THEOFFICE_ENDPOINT environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := os.Getenv("THEOFFICE_ENDPOINT")

	if !data.Endpoint.IsNull() {
		endpoint = data.Endpoint.ValueString()
	}

	// Example client configuration for data sources and resources
	client, err := theoffice.NewClient(&theoffice.Config{
		Address: endpoint,
	})
	if err != nil {
		resp.Diagnostics.AddError("error configuring theOffice client", err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured theOffice client", map[string]any{"success": true})
}

func (p *theOfficeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *theOfficeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewQuotesDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &theOfficeProvider{
			version: version,
		}
	}
}
