package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiConfiguration_res_basic(t *testing.T) {
	t.Run("match by Configuration ID", func(t *testing.T) {
		client := &mockappsec{}

		createConfigResponse := appsec.CreateConfigurationResponse{}
		createConfigResponseJSON := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/ConfigurationCreate.json"))
		json.Unmarshal([]byte(createConfigResponseJSON), &createConfigResponse)

		readConfigResponse := appsec.GetConfigurationResponse{}
		readConfigResponseJSON := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/Configuration.json"))
		json.Unmarshal([]byte(readConfigResponseJSON), &readConfigResponse)

		deleteConfigResponse := appsec.RemoveConfigurationResponse{}
		deleteConfigResponseJSON := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/Configuration.json"))
		json.Unmarshal([]byte(deleteConfigResponseJSON), &deleteConfigResponse)

		getConfigurationVersionsResponse := appsec.GetConfigurationVersionsResponse{}
		getConfigurationVersionJSON := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/ConfigurationVersions.json"))
		json.Unmarshal([]byte(getConfigurationVersionJSON), &getConfigurationVersionsResponse)

		hns := appsec.GetSelectedHostnamesResponse{}
		expectJSHN := compactJSON(loadFixtureBytes("testdata/TestResSelectedHostname/SelectedHostname.json"))
		json.Unmarshal([]byte(expectJSHN), &hns)

		client.On("GetSelectedHostnames",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetSelectedHostnamesRequest{ConfigID: 43253, Version: 7},
		).Return(&hns, nil)

		client.On("CreateConfiguration",
			mock.Anything, // ctx is irrelevant for this test
			appsec.CreateConfigurationRequest{Name: "Akamai Tools", Description: "Akamai Tools", ContractID: "C-1FRYVV3", GroupID: 64867, Hostnames: []string{"rinaldi.sandbox.akamaideveloper.com", "sujala.sandbox.akamaideveloper.com"}},
		).Return(&createConfigResponse, nil)

		client.On("GetConfiguration",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&readConfigResponse, nil)

		client.On("RemoveConfiguration",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveConfigurationRequest{ConfigID: 43253},
		).Return(&deleteConfigResponse, nil)

		client.On("GetConfigurationVersions",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetConfigurationVersionsRequest{ConfigID: 43253},
		).Return(&getConfigurationVersionsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResConfiguration/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_configuration.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestAccAkamaiConfiguration_res_error_updating_configuration(t *testing.T) {
	t.Run("match by Configuration ID", func(t *testing.T) {
		client := &mockappsec{}

		createConfigResponse := appsec.CreateConfigurationResponse{}
		createConfigResponseJSON := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/ConfigurationCreate.json"))
		json.Unmarshal([]byte(createConfigResponseJSON), &createConfigResponse)

		readConfigResponse := appsec.GetConfigurationResponse{}
		readConfigResponseJSON := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/Configuration.json"))
		json.Unmarshal([]byte(readConfigResponseJSON), &readConfigResponse)

		deleteConfigResponse := appsec.RemoveConfigurationResponse{}
		deleteConfigResponseJSON := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/Configuration.json"))
		json.Unmarshal([]byte(deleteConfigResponseJSON), &deleteConfigResponse)

		getConfigurationVersionsResponse := appsec.GetConfigurationVersionsResponse{}
		getConfigurationVersionJSON := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/ConfigurationVersions.json"))
		json.Unmarshal([]byte(getConfigurationVersionJSON), &getConfigurationVersionsResponse)

		hns := appsec.GetSelectedHostnamesResponse{}
		expectJSHN := compactJSON(loadFixtureBytes("testdata/TestResSelectedHostname/SelectedHostname.json"))
		json.Unmarshal([]byte(expectJSHN), &hns)

		client.On("GetSelectedHostnames",
			mock.Anything,
			appsec.GetSelectedHostnamesRequest{ConfigID: 43253, Version: 7},
		).Return(&hns, nil)

		client.On("CreateConfiguration",
			mock.Anything,
			appsec.CreateConfigurationRequest{Name: "Akamai Tools", Description: "Akamai Tools", ContractID: "C-1FRYVV3", GroupID: 64867, Hostnames: []string{"rinaldi.sandbox.akamaideveloper.com", "sujala.sandbox.akamaideveloper.com"}},
		).Return(&createConfigResponse, nil)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&readConfigResponse, nil)

		client.On("UpdateConfiguration",
			mock.Anything,
			appsec.UpdateConfigurationRequest{ConfigID: 43253, Name: "Akamai Tools", Description: "Akamai Tools"},
		).Return(nil, fmt.Errorf("UpdateConfiguration failed"))

		client.On("RemoveConfiguration",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveConfigurationRequest{ConfigID: 43253},
		).Return(&deleteConfigResponse, nil)

		client.On("GetConfigurationVersions",
			mock.Anything,
			appsec.GetConfigurationVersionsRequest{ConfigID: 43253},
		).Return(&getConfigurationVersionsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResConfiguration/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_configuration.test", "id", "43253"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResConfiguration/modify_contract.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_configuration.test", "id", "43253"),
						),
						ExpectError: regexp.MustCompile(`UpdateConfiguration failed`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
