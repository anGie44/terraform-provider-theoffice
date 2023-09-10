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

var (
	_ datasource.DataSource = &ConnectionsDataSource{}
)

func NewConnectionsDataSource() datasource.DataSource {
	return &ConnectionsDataSource{}
}

type ConnectionsDataSource struct {
	client *theoffice.Client
}

type ConnectionsDataSourceModel struct {
	Connections []connectionsModel `tfsdk:"connections"`
	Season      types.Int64        `tfsdk:"season"`
	ID          types.String       `tfsdk:"id"`
}

type connectionsModel struct {
	Episode     types.Int64            `tfsdk:"episode"`
	EpisodeName types.String           `tfsdk:"episode_name"`
	Links       []connectionsLinkModel `tfsdk:"links"`
	Nodes       []connectionsNodeModel `tfsdk:"nodes"`
}

type connectionsLinkModel struct {
	Source types.String `tfsdk:"source"`
	Target types.String `tfsdk:"target"`
	Value  types.Int64  `tfsdk:"value"`
}

type connectionsNodeModel struct {
	ID types.String `tfsdk:"id"`
}

func (d *ConnectionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connections"
}

func (d *ConnectionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fetches a list of character connections",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier attribute.",
				Computed:    true,
			},
			"season": schema.Int64Attribute{
				Required:    true,
				Description: "Season number to filter results by",
			},
			"connections": schema.ListNestedAttribute{
				Description: "List of character connections",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"episode": schema.Int64Attribute{
							Description: "The episode the connection, i.e. dialogue between characters, occurred in.",
							Computed:    true,
						},
						"episode_name": schema.StringAttribute{
							Description: "The name of the episode the connection, i.e. dialogue between characters, occurred in.",
							Computed:    true,
						},
						"links": schema.ListNestedAttribute{
							Description: "The list of links between characters",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"source": schema.StringAttribute{
										Description: "The name of the source character.",
										Computed:    true,
									},
									"target": schema.StringAttribute{
										Description: "The name of the target character.",
										Computed:    true,
									},
									"value": schema.Int64Attribute{
										Description: "The value of the link.",
										Computed:    true,
									},
								},
							},
						},
						"nodes": schema.ListNestedAttribute{
							Description: "The list of nodes i.e. characters in the episode",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "The name of the target character.",
										Computed:    true,
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

func (d *ConnectionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ConnectionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ConnectionsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Read Terraform configuration data into the model
	connResp, err := d.client.GetConnections(ctx, int(data.Season.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read theOffice Connections",
			err.Error(),
		)
		return
	}

	for _, conn := range connResp.Connections {
		connectionsState := connectionsModel{
			Episode:     types.Int64Value(int64(conn.Episode)),
			EpisodeName: types.StringValue(conn.EpisodeName),
		}

		for _, link := range conn.Links {
			connectionsState.Links = append(connectionsState.Links, connectionsLinkModel{
				Source: types.StringValue(link.Source),
				Target: types.StringValue(link.Target),
				Value:  types.Int64Value(int64(link.Value)),
			})
		}

		for _, node := range conn.Nodes {
			connectionsState.Nodes = append(connectionsState.Nodes, connectionsNodeModel{
				ID: types.StringValue(node.ID),
			})
		}

		data.Connections = append(data.Connections, connectionsState)
	}

	data.ID = types.StringValue("placeholder")

	// Save data into Terraform state
	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "read connections data source")
}
