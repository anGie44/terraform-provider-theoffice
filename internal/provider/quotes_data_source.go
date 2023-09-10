// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/anGie44/terraform-provider-theoffice/internal/theoffice"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource = &QuotesDataSource{}
)

func NewQuotesDataSource() datasource.DataSource {
	return &QuotesDataSource{}
}

// QuotesDataSource defines the data source implementation.
type QuotesDataSource struct {
	client *theoffice.Client
}

// QuotesDataSourceModel describes the data source data model.
type QuotesDataSourceModel struct {
	Episode types.Int64   `tfsdk:"episode"`
	Season  types.Int64   `tfsdk:"season"`
	Quotes  []quotesModel `tfsdk:"quotes"`
	ID      types.String  `tfsdk:"id"`
}

type quotesModel struct {
	Season      types.Int64  `tfsdk:"season"`
	Episode     types.Int64  `tfsdk:"episode"`
	Scene       types.Int64  `tfsdk:"scene"`
	EpisodeName types.String `tfsdk:"episode_name"`
	Character   types.String `tfsdk:"character"`
	Quote       types.String `tfsdk:"quote"`
}

func (d *QuotesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_quotes"
}

func (d *QuotesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fetches a list of quotes",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier attribute.",
				Computed:    true,
			},
			"episode": schema.Int64Attribute{
				Optional:    true,
				Description: "Episode number to filter results by",
			},
			"season": schema.Int64Attribute{
				Required:    true,
				Description: "Season number to filter results by",
			},
			"quotes": schema.ListNestedAttribute{
				Description: "List of quotes",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"season": schema.Int64Attribute{
							Description: "The season the quote occurred in.",
							Computed:    true,
						},
						"episode": schema.Int64Attribute{
							Description: "The episode the quote occurred in.",
							Computed:    true,
						},
						"scene": schema.Int64Attribute{
							Description: "The scene the quote occurred in.",
							Computed:    true,
						},
						"episode_name": schema.StringAttribute{
							Description: "The name of the episode the quote occurred in.",
							Computed:    true,
						},
						"character": schema.StringAttribute{
							Description: "The character who said the quote.",
							Computed:    true,
						},
						"quote": schema.StringAttribute{
							Description: "The quote as a string",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *QuotesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*theoffice.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *theoffice.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *QuotesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data QuotesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Read Terraform configuration data into the model
	quotes, err := d.client.GetQuotes(ctx, int(data.Season.ValueInt64()), int(data.Episode.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read theOffice Quotes",
			err.Error(),
		)
		return
	}

	for _, quote := range quotes.Quotes {
		quoteState := quotesModel{
			Season:      types.Int64Value(int64(quote.Season)),
			Episode:     types.Int64Value(int64(quote.Episode)),
			Scene:       types.Int64Value(int64(quote.Scene)),
			EpisodeName: types.StringValue(quote.EpisodeName),
			Character:   types.StringValue(quote.Character),
			Quote:       types.StringValue(quote.Quote),
		}

		data.Quotes = append(data.Quotes, quoteState)
	}

	data.ID = types.StringValue("placeholder")

	// Save data into Terraform state
	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "read quotes data source")
}
