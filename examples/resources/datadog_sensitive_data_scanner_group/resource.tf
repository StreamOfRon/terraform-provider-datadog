
# Create new sensitive_data_scanner_group resource


resource "datadog_sensitive_data_scanner_group" "foo" {
  description = "UPDATE ME"
  filter {
    query = "UPDATE ME"
  }
  is_enabled   = "UPDATE ME"
  name         = "UPDATE ME"
  product_list = "UPDATE ME"
}