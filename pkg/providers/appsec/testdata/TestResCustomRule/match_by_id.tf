provider "akamai" {
  edgerc        = "~/.edgerc"
  cache_enabled = false
}

resource "akamai_appsec_custom_rule" "test" {
  config_id   = 43253
  custom_rule = <<-EOF
{
    "name": "Rule Test New",
    "description": "Can I create all conditions?",
    "tag": [
        "test"
    ],
    "conditions": [{
            "type": "requestMethodMatch",
            "positiveMatch": true,
            "value": [
                "GET",
                "CONNECT",
                "TRACE",
                "PUT",
                "POST",
                "OPTIONS",
                "DELETE",
                "HEAD"
            ]
        },
        {
            "type": "pathMatch",
            "positiveMatch": true,
            "value": [
                "/H",
                "/Li",
                "/He"
            ]
        },
        {
            "type": "extensionMatch",
            "positiveMatch": true,
            "valueWildcard": true,
            "valueCase": true,
            "value": [
                "Li",
                "He",
                "H"
            ]
        }
    ]
}
EOF
}

