provider "akamai" {
  edgerc        = "~/.edgerc"
  cache_enabled = false
}

data "akamai_appsec_match_targets" "test" {
  config_id = 43253
}

