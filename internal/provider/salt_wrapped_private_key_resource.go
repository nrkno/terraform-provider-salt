package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/nrkno/terraform-provider-salt/internal/salt"
)

var (
	_ resource.Resource                = &SaltWrappedPrivateKeyResource{}
	_ resource.ResourceWithImportState = &SaltWrappedPrivateKeyResource{}
)

func NewSaltWrappedPrivateKeyResource() resource.Resource {
	return &SaltWrappedPrivateKeyResource{}
}

type SaltWrappedPrivateKeyResource struct {
	client *salt.Client
}

type WrappedPrivateKeyResourceModel struct {
	MinionId          types.String `tfsdk:"minion_id"`
	WrappedPrivateKey types.String `tfsdk:"wrapped_private_key"`
}

func (n *SaltWrappedPrivateKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Trace(ctx, "SaltWrappedPrivateKeyResource.Metadata")
	resp.TypeName = req.ProviderTypeName + "_wrapped_private_key"
}

// Schema should return the schema for this resource.
func (n *SaltWrappedPrivateKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tflog.Trace(ctx, "SaltWrappedPrivateKeyResource.Schema")
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generate keypair, accept public key and return a wrapped private key",
		Attributes: map[string]schema.Attribute{
			"minion_id": schema.StringAttribute{
				MarkdownDescription: "Minion ID",
				Description:         "The minion ID to create a key pair for",
				Required:            true,
			},
			"wrapped_private_key": schema.StringAttribute{
				MarkdownDescription: "Vault wrapped private key",
				Description:         "A Vault wrapped token to get the private key with",
				Computed:            true,
			},
		},
	}
}

func (n *SaltWrappedPrivateKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*salt.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	n.client = client
}

// Create is called when the provider must create a new resource. Config
// and planned state values should be read from the
// CreateRequest and new state values set on the CreateResponse.
func (n *SaltWrappedPrivateKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "SaltWrappedPrivateKeyResource.Create")
	var plan WrappedPrivateKeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	wrappedToken := n.client.WrappedPrivateKey(plan.MinionId.String())

	resourceModel := WrappedPrivateKeyResourceModel{
		MinionId:          plan.MinionId,
		WrappedPrivateKey: types.StringValue(wrappedToken),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, resourceModel)...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("error updating terraform state %v", resp.Diagnostics.Errors()))
		return
	}
}

// Read is called when the provider must read resource values in order
// to update state. Planned state values should be read from the
// ReadRequest and new state values set on the ReadResponse.
func (n *SaltWrappedPrivateKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Trace(ctx, "SaltWrappedPrivateKeyResource.Read")
}

// Update is called to update the state of the resource. Config, planned
// state, and prior state values should be read from the
// UpdateRequest and new state values set on the UpdateResponse.
func (n *SaltWrappedPrivateKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Trace(ctx, "SaltWrappedPrivateKeyResource.Update")
}

// Delete is called when the provider must delete the resource. Config
// values may be read from the DeleteRequest.
//
// If execution completes without error, the framework will automatically
// call DeleteResponse.State.RemoveResource(), so it can be omitted
// from provider logic.
func (n *SaltWrappedPrivateKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "SaltWrappedPrivateKeyResource.Delete")
}

func (n *SaltWrappedPrivateKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
