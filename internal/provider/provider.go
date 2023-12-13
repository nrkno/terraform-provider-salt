// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nrkno/terraform-provider-salt/internal/salt"
)

// Ensure SaltProvider satisfies various provider interfaces.
var _ provider.Provider = &SaltProvider{}

// SaltProvider defines the provider implementation.
type SaltProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// SaltProviderModel describes the provider data model.
type SaltProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *SaltProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "salt"
	resp.Version = p.version
}

func (p *SaltProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "API endpoint to use",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username used to authenticate to the Salt API",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password used to authenticate to the Salt API",
				Required:            true,
			},
		},
	}
}

func (p *SaltProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data SaltProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	c, err := salt.New(data.Endpoint.ValueString(), data.Username.ValueString(), data.Password.ValueString(), p.version)
	if err != nil {
		resp.Diagnostics.AddError(
			"API Error Authentication",
			fmt.Sprintf("Could not authenticate to the Salt API: %s", err),
		)
		return
	}

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *SaltProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSaltWrappedPrivateKeyResource,
	}
}

func (p *SaltProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		//
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SaltProvider{
			version: version,
		}
	}
}
