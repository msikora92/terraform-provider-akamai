provider "akamai" {
  edgerc        = "~/.edgerc"
  cache_enabled = false
}

data "akamai_appsec_hostname_coverage_match_targets" "test" {
  config_id = 43253
  hostname  = "rinaldi.sandbox.akamaideveloper.com"
}

