package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-datadog/datadog"
	"github.com/terraform-providers/terraform-provider-datadog/datadog/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSensitiveDataScannerGroupSimple(t *testing.T) {
	t.Parallel()
	ctx, accProviders := testAccProviders(context.Background(), t)
	uniq := uniqueEntityName(ctx, t)
	accProvider := testAccProvider(t, accProviders)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: accProviders,
		CheckDestroy:      testAccCheckDatadogSensitiveDataScannerGroupDestroy(accProvider),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDatadogSensitiveDataScannerGroup(uniq),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatadogSensitiveDataScannerGroupExists(accProvider),
					resource.TestCheckResourceAttr(
						"datadog_sensitive_data_scanner_group.foo", "description", "UPDATE ME"),
					resource.TestCheckResourceAttr(
						"datadog_sensitive_data_scanner_group.foo", "is_enabled", "UPDATE ME"),
					resource.TestCheckResourceAttr(
						"datadog_sensitive_data_scanner_group.foo", "name", "UPDATE ME"),
				),
			},
		},
	})
}

func testAccCheckDatadogSensitiveDataScannerGroup(uniq string) string {
	// Update me to make use of the unique value
	return fmt.Sprintf(`
resource "datadog_sensitive_data_scanner_group" "foo" {
    description = "UPDATE ME"
    filter {
    query = "UPDATE ME"
    }
    is_enabled = "UPDATE ME"
    name = "UPDATE ME"
    product_list = "UPDATE ME"
}`, uniq)
}

func testAccCheckDatadogSensitiveDataScannerGroupDestroy(accProvider func() (*schema.Provider, error)) func(*terraform.State) error {
	return func(s *terraform.State) error {
		provider, _ := accProvider()
		providerConf := provider.Meta().(*datadog.ProviderConfiguration)
		apiInstances := providerConf.DatadogApiInstances
		auth := providerConf.Auth

		if err := SensitiveDataScannerGroupDestroyHelper(auth, s, apiInstances); err != nil {
			return err
		}
		return nil
	}
}

func SensitiveDataScannerGroupDestroyHelper(auth context.Context, s *terraform.State, apiInstances *utils.ApiInstances) error {
	err := utils.Retry(2, 10, func() error {
		for _, r := range s.RootModule().Resources {
			if r.Type != "resource_datadog_sensitive_data_scanner_group" {
				continue
			}

			_, httpResp, err := apiInstances.GetSensitiveDataScannerApiV2().ListScanningGroups(auth)
			if err != nil {
				if httpResp != nil && httpResp.StatusCode == 404 {
					return nil
				}
				return &utils.RetryableError{Prob: fmt.Sprintf("received an error retrieving Monitor %s", err)}
			}
			return &utils.RetryableError{Prob: "Monitor still exists"}
		}
		return nil
	})
	return err
}

func testAccCheckDatadogSensitiveDataScannerGroupExists(accProvider func() (*schema.Provider, error)) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		provider, _ := accProvider()
		providerConf := provider.Meta().(*datadog.ProviderConfiguration)
		apiInstances := providerConf.DatadogApiInstances
		auth := providerConf.Auth

		if err := sensitiveDataScannerGroupExistsHelper(auth, s, apiInstances); err != nil {
			return err
		}
		return nil
	}
}

func sensitiveDataScannerGroupExistsHelper(auth context.Context, s *terraform.State, apiInstances *utils.ApiInstances) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "resource_datadog_sensitive_data_scanner_group" {
			continue
		}

		_, httpResp, err := apiInstances.GetSensitiveDataScannerApiV2().ListScanningGroups(auth)
		if err != nil {
			return utils.TranslateClientError(err, httpResp, "error retrieving monitor")
		}
	}
	return nil
}
