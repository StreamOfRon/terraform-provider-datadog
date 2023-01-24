package datadog

import (
	"context"
	"strconv"

	"github.com/terraform-providers/terraform-provider-datadog/datadog/internal/utils"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatadogSensitiveDataScannerGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Datadog SensitiveDataScannerGroup resource. This can be used to create and manage Datadog sensitive_data_scanner_group.",
		ReadContext:   resourceDatadogSensitiveDataScannerGroupRead,
		CreateContext: resourceDatadogSensitiveDataScannerGroupCreate,
		UpdateContext: resourceDatadogSensitiveDataScannerGroupUpdate,
		DeleteContext: resourceDatadogSensitiveDataScannerGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the group.",
			},
			"filter": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Filter for the Scanning Group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Query to filter the events.",
						},
					},
				},
			},
			"is_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether or not the group is enabled.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the group.",
			},
			"product_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of products the scanning group applies.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceDatadogSensitiveDataScannerGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConf := meta.(*ProviderConfiguration)
	apiInstances := providerConf.DatadogApiInstances
	auth := providerConf.Auth

	resp, httpResp, err := apiInstances.GetSensitiveDataScannerApiV2().ListScanningGroups(auth)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			// Delete the resource from the local state since it doesn't exist anymore in the actual state
			d.SetId("")
			return nil
		}
		return utils.TranslateClientErrorDiag(err, httpResp, "error calling ListScanningGroups")
	}
	if err := utils.CheckForUnparsed(resp); err != nil {
		return diag.FromErr(err)
	}

	return updateSensitiveDataScannerGroupState(d, &resp)
}

func resourceDatadogSensitiveDataScannerGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConf := meta.(*ProviderConfiguration)
	apiInstances := providerConf.DatadogApiInstances
	auth := providerConf.Auth

	body := buildSensitiveDataScannerGroupRequestBody(d)

	resp, httpResp, err := apiInstances.GetSensitiveDataScannerApiV2().CreateScanningGroup(auth, *body)
	if err != nil {
		return utils.TranslateClientErrorDiag(err, httpResp, "error creating SensitiveDataScannerGroup")
	}
	if err := utils.CheckForUnparsed(resp); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.Data.GetId())

	return updateSensitiveDataScannerGroupState(d, &resp)
}

func buildSensitiveDataScannerGroupRequestBody(d *schema.ResourceData) *datadogV2.SensitiveDataScannerGroupCreateRequest {
	attributes := datadogV2.NewSensitiveDataScannerGroupAttributesWithDefaults()

	if description, ok := d.GetOk("description"); ok {
		attributes.SetDescription(description.(string))
	}
	filter := datadogV2.NewSensitiveDataScannerFilterWithDefaults()

	if query, ok := d.GetOk("query"); ok {
		filter.SetQuery(query.(string))
	}
	attributes.SetFilter(*filter)

	if isEnabled, ok := d.GetOk("is_enabled"); ok {
		attributes.SetIsEnabled(isEnabled.(bool))
	}

	if name, ok := d.GetOk("name"); ok {
		attributes.SetName(name.(string))
	}
	productList := []datadogV2.SensitiveDataScannerProduct{}
	for _, s := range d.Get("product_list").([]interface{}) {
		sensitiveDataScannerProductItem, _ := datadogV2.NewSensitiveDataScannerProductFromValue(s.(string))
		productList = append(productList, *sensitiveDataScannerProductItem)
	}
	attributes.SetProductList(productList)

	req := datadogV2.NewSensitiveDataScannerGroupCreateRequestWithDefaults()
	req.Data = datadogV2.NewSensitiveDataScannerGroupCreateWithDefaults()
	req.Data.SetAttributes(*attributes)

	return req
}

func resourceDatadogSensitiveDataScannerGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConf := meta.(*ProviderConfiguration)
	apiInstances := providerConf.DatadogApiInstances
	auth := providerConf.Auth

	id := d.Id()

	body := buildSensitiveDataScannerGroupUpdateRequestBody(d)

	resp, httpResp, err := apiInstances.GetSensitiveDataScannerApiV2().UpdateScanningGroup(auth, id, *body)
	if err != nil {
		return utils.TranslateClientErrorDiag(err, httpResp, "error creating SensitiveDataScannerGroup")
	}
	if err := utils.CheckForUnparsed(resp); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.Data.GetId())

	return updateSensitiveDataScannerGroupState(d, &resp)
}

func buildSensitiveDataScannerGroupUpdateRequestBody(d *schema.ResourceData) *datadogV2.SensitiveDataScannerGroupUpdateRequest {
	attributes := datadogV2.NewSensitiveDataScannerGroupAttributesWithDefaults()

	if description, ok := d.GetOk("description"); ok {
		attributes.SetDescription(description.(string))
	}
	filter := datadogV2.NewSensitiveDataScannerFilterWithDefaults()

	if query, ok := d.GetOk("query"); ok {
		filter.SetQuery(query.(string))
	}
	attributes.SetFilter(*filter)

	if isEnabled, ok := d.GetOk("is_enabled"); ok {
		attributes.SetIsEnabled(isEnabled.(bool))
	}

	if name, ok := d.GetOk("name"); ok {
		attributes.SetName(name.(string))
	}
	productList := []datadogV2.SensitiveDataScannerProduct{}
	for _, s := range d.Get("product_list").([]interface{}) {
		sensitiveDataScannerProductItem, _ := datadogV2.NewSensitiveDataScannerProductFromValue(s.(string))
		productList = append(productList, *sensitiveDataScannerProductItem)
	}
	attributes.SetProductList(productList)

	req := datadogV2.NewSensitiveDataScannerGroupUpdateRequestWithDefaults()
	req.Data = *datadogV2.NewSensitiveDataScannerGroupUpdateWithDefaults()
	req.Data.SetAttributes(*attributes)

	return req
}

func resourceDatadogSensitiveDataScannerGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConf := meta.(*ProviderConfiguration)
	apiInstances := providerConf.DatadogApiInstances
	auth := providerConf.Auth

	id := d.Id()
	body := datadogV2.NewSensitiveDataScannerGroupDeleteRequestWithDefaults()
	metaVar := datadogV2.NewSensitiveDataScannerMetaVersionOnlyWithDefaults()

	if version, ok := d.GetOk("version"); ok {
		versionInt, _ := strconv.ParseInt(version.(string), 10, 64)
		metaVar.SetVersion(versionInt)
	}
	body.SetMeta(*metaVar)

	_, httpResp, err := apiInstances.GetSensitiveDataScannerApiV2().DeleteScanningGroup(auth, id, body)
	if err != nil {
		// The resource is assumed to still exist, and all prior state is preserved.
		return utils.TranslateClientErrorDiag(err, httpResp, "error deleting SensitiveDataScannerGroup")
	}

	return nil
}

func updateSensitiveDataScannerGroupState(d *schema.ResourceData, resp *datadogV2.SensitiveDataScannerGetConfigResponse) diag.Diagnostics {

	return nil
}
