package property

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
)

func TestResourceEdgeHostname(t *testing.T) {
	testDir := "testdata/TestResourceEdgeHostname"
	tests := map[string]struct {
		givenTF            string
		init               func(*mockpapi)
		expectedAttributes map[string]string
		expectedOutputs    map[string]string
		withError          *regexp.Regexp
	}{
		"edge hostname with .edgesuite.net, create edge hostname": {
			givenTF: "new_edgesuite_net.tf",
			init: func(m *mockpapi) {
				m.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "eh_123",
							Domain:       "test2.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test",
							DomainSuffix: "edgesuite.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test3.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
					}},
				}, nil).Once()
				m.On("CreateEdgeHostname", mock.Anything, papi.CreateEdgeHostnameRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostname: papi.EdgeHostnameCreate{
						ProductID:         "prd_2",
						DomainPrefix:      "test2",
						DomainSuffix:      "edgesuite.net",
						SecureNetwork:     "STANDARD_TLS",
						IPVersionBehavior: "IPV6_COMPLIANCE",
						CertEnrollmentID:  123,
						SlotNumber:        123,
					},
				}).Return(&papi.CreateEdgeHostnameResponse{
					EdgeHostnameID: "eh_123",
				}, nil)
				m.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:                "eh_123",
							Domain:            "test2.edgesuite.net",
							ProductID:         "prd_2",
							DomainPrefix:      "test2",
							DomainSuffix:      "edgesuite.net",
							IPVersionBehavior: "IPV6_COMPLIANCE",
						},
						{
							ID:           "eh_2",
							Domain:       "test.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
						{
							ID:           "eh_123",
							Domain:       "test.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test",
							DomainSuffix: "edgesuite.net",
						},
					}},
				}, nil)
			},
			expectedAttributes: map[string]string{
				"id":            "eh_123",
				"ip_behavior":   "IPV6_COMPLIANCE",
				"contract":      "ctr_2",
				"group":         "grp_2",
				"edge_hostname": "test2.edgesuite.net",
			},
			expectedOutputs: map[string]string{
				"edge_hostname": "test2.edgesuite.net",
			},
		},
		"edge hostname with .edgekey.net, create edge hostname": {
			givenTF: "new_edgekey_net.tf",
			init: func(m *mockpapi) {
				m.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "eh_123",
							Domain:       "test.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test2",
							DomainSuffix: "edgesuite.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
					}},
				}, nil).Once()
				m.On("CreateEdgeHostname", mock.Anything, papi.CreateEdgeHostnameRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostname: papi.EdgeHostnameCreate{
						ProductID:         "prd_2",
						DomainPrefix:      "test",
						DomainSuffix:      "edgekey.net",
						SecureNetwork:     "ENHANCED_TLS",
						IPVersionBehavior: "IPV6_PERFORMANCE",
						CertEnrollmentID:  123,
						SlotNumber:        123,
					},
				}).Return(&papi.CreateEdgeHostnameResponse{
					EdgeHostnameID: "eh_123",
				}, nil)
				m.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "eh_123",
							Domain:       "test.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test2",
							DomainSuffix: "edgesuite.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test3.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
						{
							ID:           "eh_123",
							Domain:       "test.edgekey.net",
							ProductID:    "prd_2",
							DomainPrefix: "test",
							DomainSuffix: "edgekey.net",
						},
					}},
				}, nil)
			},
			expectedAttributes: map[string]string{
				"id":            "eh_123",
				"ip_behavior":   "IPV6_PERFORMANCE",
				"contract":      "ctr_2",
				"group":         "grp_2",
				"edge_hostname": "test.edgekey.net",
			},
			expectedOutputs: map[string]string{
				"edge_hostname": "test.edgekey.net",
			},
		},
		"edge hostname with .akamaized.net, create edge hostname": {
			givenTF: "new_akamaized_net.tf",
			init: func(m *mockpapi) {
				m.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "eh_123",
							Domain:       "test.akamaized.net",
							ProductID:    "prd_2",
							DomainPrefix: "test2",
							DomainSuffix: "akamaized.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test.akamaized.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "akamaized.net",
						},
					}},
				}, nil).Once()
				m.On("CreateEdgeHostname", mock.Anything, papi.CreateEdgeHostnameRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostname: papi.EdgeHostnameCreate{
						ProductID:         "prd_2",
						DomainPrefix:      "test",
						DomainSuffix:      "akamaized.net",
						SecureNetwork:     "SHARED_CERT",
						IPVersionBehavior: "IPV6_COMPLIANCE",
					},
				}).Return(&papi.CreateEdgeHostnameResponse{
					EdgeHostnameID: "eh_123",
				}, nil)
				m.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "eh_123",
							Domain:       "test.akamaized.net",
							ProductID:    "prd_2",
							DomainPrefix: "test2",
							DomainSuffix: "akamaized.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test.akamaized.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "akamaized.net",
						},
						{
							ID:           "eh_123",
							Domain:       "test.akamaized.net",
							ProductID:    "prd_2",
							DomainPrefix: "test",
							DomainSuffix: "akamaized.net",
						},
					}},
				}, nil)
			},
			expectedAttributes: map[string]string{
				"id":            "eh_123",
				"ip_behavior":   "IPV6_COMPLIANCE",
				"contract":      "ctr_2",
				"group":         "grp_2",
				"edge_hostname": "test.akamaized.net",
			},
		},
		"different edge hostname, create": {
			givenTF: "new.tf",
			init: func(m *mockpapi) {
				m.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "eh_123",
							Domain:       "test.aka.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test2",
							DomainSuffix: "aka.net.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
					}},
				}, nil).Once()
				m.On("CreateEdgeHostname", mock.Anything, papi.CreateEdgeHostnameRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostname: papi.EdgeHostnameCreate{
						ProductID:         "prd_2",
						DomainPrefix:      "test.aka",
						DomainSuffix:      "edgesuite.net",
						SecureNetwork:     "",
						IPVersionBehavior: "IPV4",
						UseCases: []papi.UseCase{
							{
								UseCase: "Download_Mode",
								Option:  "BACKGROUND",
								Type:    "GLOBAL",
							},
						},
					},
				}).Return(&papi.CreateEdgeHostnameResponse{
					EdgeHostnameID: "eh_123",
				}, nil)
				m.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "eh_1",
							Domain:       "test2.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test2",
							DomainSuffix: "edgesuite.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test3.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
						{
							ID:           "eh_123",
							Domain:       "test.aka.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test.aka",
							DomainSuffix: "edgesuite.net",
							UseCases: []papi.UseCase{
								{
									UseCase: "Download_Mode",
									Option:  "BACKGROUND",
									Type:    "GLOBAL",
								},
							},
						},
					}},
				}, nil)
			},
			expectedAttributes: map[string]string{
				"id":            "eh_123",
				"ip_behavior":   "IPV4",
				"contract":      "ctr_2",
				"group":         "grp_2",
				"edge_hostname": "test.aka.edgesuite.net",
				"use_cases":     loadFixtureString(fmt.Sprintf("%s/use_cases/use_cases_new.json", testDir)),
			},
			expectedOutputs: map[string]string{
				"edge_hostname": "test.aka.edgesuite.net",
			},
		},
		"edge hostname exists": {
			givenTF: "new_akamaized_net.tf",
			init: func(m *mockpapi) {
				m.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "eh_123",
							Domain:       "test.akamaized.net",
							ProductID:    "prd_2",
							DomainPrefix: "test",
							DomainSuffix: "akamaized.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test.akamaized.net",
							ProductID:    "prd_2",
							DomainPrefix: "test",
							DomainSuffix: "akamaized.net",
						},
					}},
				}, nil)
			},
			expectedAttributes: map[string]string{
				"id":            "eh_123",
				"contract":      "ctr_2",
				"group":         "grp_2",
				"edge_hostname": "test.akamaized.net",
			},
		},
		"error fetching edge hostnames": {
			givenTF: "new_akamaized_net.tf",
			init: func(m *mockpapi) {
				m.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(nil, fmt.Errorf("oops"))
			},
			withError: regexp.MustCompile("oops"),
		},
		"invalid IP version behavior": {
			givenTF:   "invalid_ip.tf",
			init:      func(m *mockpapi) {},
			withError: regexp.MustCompile("ip_behavior must be one of IPV4, IPV6_PERFORMANCE, IPV6_COMPLIANCE"),
		},
		"certificate required for ENHANCED_TLS": {
			givenTF: "missing_certificate.tf",
			init: func(m *mockpapi) {
				m.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{}, nil)
			},
			withError: regexp.MustCompile("a certificate enrollment ID is required for Enhanced TLS edge hostnames with 'edgekey.net' suffix"),
		},
		"error creating edge hostname": {
			givenTF: "new_akamaized_net.tf",
			init: func(m *mockpapi) {
				m.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "eh_1",
							Domain:       "test2.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test2",
							DomainSuffix: "edgesuite.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test3.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
					}},
				}, nil)
				m.On("CreateEdgeHostname", mock.Anything, papi.CreateEdgeHostnameRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostname: papi.EdgeHostnameCreate{
						ProductID:         "prd_2",
						DomainPrefix:      "test",
						DomainSuffix:      "akamaized.net",
						SecureNetwork:     "SHARED_CERT",
						IPVersionBehavior: "IPV6_COMPLIANCE",
					},
				}).Return(nil, fmt.Errorf("oops"))
			},
			withError: regexp.MustCompile("oops"),
		},
		"error edge hostname not found": {
			givenTF: "new_akamaized_net.tf",
			init: func(m *mockpapi) {
				m.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
				}).Return(&papi.GetEdgeHostnamesResponse{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
						{
							ID:           "eh_1",
							Domain:       "test2.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test2",
							DomainSuffix: "edgesuite.net",
						},
						{
							ID:           "eh_2",
							Domain:       "test3.edgesuite.net",
							ProductID:    "prd_2",
							DomainPrefix: "test3",
							DomainSuffix: "edgesuite.net",
						},
					}},
				}, nil).Twice()
				m.On("CreateEdgeHostname", mock.Anything, papi.CreateEdgeHostnameRequest{
					ContractID: "ctr_2",
					GroupID:    "grp_2",
					EdgeHostname: papi.EdgeHostnameCreate{
						ProductID:         "prd_2",
						DomainPrefix:      "test",
						DomainSuffix:      "akamaized.net",
						SecureNetwork:     "SHARED_CERT",
						IPVersionBehavior: "IPV6_COMPLIANCE",
					},
				}).Return(&papi.CreateEdgeHostnameResponse{
					EdgeHostnameID: "eh_123",
				}, nil)
			},
			withError: regexp.MustCompile("unable to find edge hostname"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mockpapi{}
			test.init(client)
			var checkFuncs []resource.TestCheckFunc
			for k, v := range test.expectedAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_edge_hostname.edgehostname", k, v))
			}
			for k, v := range test.expectedOutputs {
				checkFuncs = append(checkFuncs, resource.TestCheckOutput(k, v))
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config:      loadFixtureString(fmt.Sprintf("%s/%s", testDir, test.givenTF)),
							Check:       resource.ComposeAggregateTestCheckFunc(checkFuncs...),
							ExpectError: test.withError,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func TestResourceEdgeHostnames_WithImport(t *testing.T) {
	expectGetEdgeHostname := func(m *mockpapi, edgehostID, ContractID, GroupID string) *mock.Call {
		return m.On("GetEdgeHostname", mock.Anything, papi.GetEdgeHostnameRequest{
			EdgeHostnameID: edgehostID,
			ContractID:     ContractID,
			GroupID:        GroupID,
		}).Return(&papi.GetEdgeHostnamesResponse{
			ContractID: "ctr_1",
			GroupID:    "grp_2",
			EdgeHostname: papi.EdgeHostnameGetItem{
				ID:                "eh_1",
				Domain:            "test.akamaized.net",
				ProductID:         "prd_2",
				DomainPrefix:      "test",
				DomainSuffix:      "akamaized.net",
				IPVersionBehavior: "IPV4",
			},
			EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:                "eh_1",
					Domain:            "test.akamaized.net",
					ProductID:         "prd_2",
					DomainPrefix:      "test",
					DomainSuffix:      "akamaized.net",
					IPVersionBehavior: "IPV4",
				},
				{
					ID:                "eh_2",
					Domain:            "test3.edgesuite.net",
					ProductID:         "prd_2",
					DomainPrefix:      "test3",
					DomainSuffix:      "edgesuite.net",
					IPVersionBehavior: "IPV4",
				},
			}},
		}, nil)
	}

	expectGetEdgeHostnames := func(m *mockpapi, ContractID, GroupID string) *mock.Call {
		return m.On("GetEdgeHostnames", mock.Anything, papi.GetEdgeHostnamesRequest{
			ContractID: ContractID,
			GroupID:    GroupID,
		}).Return(&papi.GetEdgeHostnamesResponse{
			ContractID: "ctr_1",
			GroupID:    "grp_2",
			EdgeHostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:                "eh_1",
					Domain:            "test.akamaized.net",
					ProductID:         "prd_2",
					DomainPrefix:      "test",
					DomainSuffix:      "akamaized.net",
					IPVersionBehavior: "IPV4",
				},
				{
					ID:                "eh_2",
					Domain:            "test3.edgesuite.net",
					ProductID:         "prd_2",
					DomainPrefix:      "test3",
					DomainSuffix:      "edgesuite.net",
					IPVersionBehavior: "IPV4",
				},
			}},
		}, nil)
	}

	t.Run("import existing edgehostname code", func(t *testing.T) {
		client := &mockpapi{}
		id := "eh_1,1,2"

		expectGetEdgeHostname(client, "eh_1", "ctr_1", "grp_2")
		expectGetEdgeHostnames(client, "ctr_1", "grp_2")
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResourceEdgeHostname/import_edgehostname.tf"),
					},
					{
						Config:      loadFixtureString("testdata/TestResourceEdgeHostname/import_edgehostname.tf"),
						ImportState: true,
						ImportStateCheck: func(s []*terraform.InstanceState) error {
							assert.Len(t, s, 1)
							rs := s[0]
							assert.Equal(t, "grp_2", rs.Attributes["group_id"])
							assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
							assert.Equal(t, "eh_1", rs.Attributes["id"])
							return nil
						},
						ImportStateId:           id,
						ResourceName:            "akamai_edge_hostname.importedgehostname",
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"ip_behavior"},
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
}

func TestFindEdgeHostname(t *testing.T) {
	tests := map[string]struct {
		hostnames papi.EdgeHostnameItems
		domain    string
		expected  *papi.EdgeHostnameGetItem
		withError error
	}{
		"edge hostname found, different domain": {
			hostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:           "eh_1",
					Domain:       "some.domain.edgesuite.net",
					DomainPrefix: "some.domain",
					DomainSuffix: "edgesuite.net",
				},
				{
					ID:           "eh_2",
					Domain:       "test.domain.edgesuite.net",
					DomainPrefix: "test.domain",
					DomainSuffix: "edgesuite.net",
				},
			}},
			domain: "some.domain",
			expected: &papi.EdgeHostnameGetItem{
				ID:           "eh_1",
				Domain:       "some.domain.edgesuite.net",
				DomainPrefix: "some.domain",
				DomainSuffix: "edgesuite.net",
			},
		},
		"edge hostname found, edgesuite domain": {
			hostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:           "eh_1",
					Domain:       "some.domain.edgesuite.net",
					DomainPrefix: "some.domain",
					DomainSuffix: "edgesuite.net",
				},
				{
					ID:           "eh_2",
					Domain:       "test.domain.edgesuite.net",
					DomainPrefix: "test.domain",
					DomainSuffix: "edgesuite.net",
				},
			}},
			domain: "some.domain.edgesuite.net",
			expected: &papi.EdgeHostnameGetItem{
				ID:           "eh_1",
				Domain:       "some.domain.edgesuite.net",
				DomainPrefix: "some.domain",
				DomainSuffix: "edgesuite.net",
			},
		},
		"edge hostname found, edgekey domain": {
			hostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:           "eh_1",
					Domain:       "some.domain.edgekey.net",
					DomainPrefix: "some.domain",
					DomainSuffix: "edgekey.net",
				},
				{
					ID:           "eh_2",
					Domain:       "test.domain.edgesuite.net",
					DomainPrefix: "test.domain",
					DomainSuffix: "edgesuite.net",
				},
			}},
			domain: "some.domain.edgekey.net",
			expected: &papi.EdgeHostnameGetItem{
				ID:           "eh_1",
				Domain:       "some.domain.edgekey.net",
				DomainPrefix: "some.domain",
				DomainSuffix: "edgekey.net",
			},
		},
		"edge hostname found, akamaized domain": {
			hostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:           "eh_1",
					Domain:       "some.domain.akamaized.net",
					DomainPrefix: "some.domain",
					DomainSuffix: "akamaized.net",
				},
				{
					ID:           "eh_2",
					Domain:       "test.domain.edgesuite.net",
					DomainPrefix: "test.domain",
					DomainSuffix: "edgesuite.net",
				},
			}},
			domain: "some.domain.akamaized.net",
			expected: &papi.EdgeHostnameGetItem{
				ID:           "eh_1",
				Domain:       "some.domain.akamaized.net",
				DomainPrefix: "some.domain",
				DomainSuffix: "akamaized.net",
			},
		},
		"edge hostname not found": {
			hostnames: papi.EdgeHostnameItems{Items: []papi.EdgeHostnameGetItem{
				{
					ID:           "eh_1",
					Domain:       "some.domain.akamaized.net",
					DomainPrefix: "some.domain",
					DomainSuffix: "akamaized.net",
				},
				{
					ID:           "eh_2",
					Domain:       "test.domain.edgesuite.net",
					DomainPrefix: "test.domain",
					DomainSuffix: "edgesuite.net",
				},
			}},
			domain:    "other.domain.akamaized.net",
			withError: ErrEdgeHostnameNotFound,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := findEdgeHostname(test.hostnames, test.domain)
			if test.withError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, test.withError), "want: %v; got: %v", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestSuppressEdgeHostnameDomain(t *testing.T) {
	tests := map[string]struct {
		old, new string
		expected bool
	}{
		"equal domains": {
			old:      "test.com",
			new:      "test.com",
			expected: true,
		},
		"domain defaulting to edgesuite.net": {
			old:      "test.com.edgesuite.net",
			new:      "test.com",
			expected: true,
		},
		"different domains": {
			old:      "test.com.akamaized.net",
			new:      "test1.com.akamaized.net",
			expected: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, suppressEdgeHostnameDomain("", test.old, test.new, nil))
		})
	}
}

func TestSuppressEdgeHostnameUseCases(t *testing.T) {
	testDir := "testdata/TestResourceEdgeHostname/use_cases"
	tests := map[string]struct {
		oldPath, newPath string
		expected         bool
	}{
		"equal use cases in the same order": {
			oldPath:  "use_cases1.json",
			newPath:  "use_cases2.json",
			expected: true,
		},
		"equal use cases in different order": {
			oldPath:  "use_cases1.json",
			newPath:  "use_cases1_mixed.json",
			expected: true,
		},
		"not equal use cases": {
			oldPath:  "use_cases1.json",
			newPath:  "use_cases3.json",
			expected: false,
		},
		"error unmarshalling new": {
			oldPath:  "use_cases1.json",
			newPath:  "invalid_use_cases.json",
			expected: false,
		},
		"error unmarshalling old": {
			oldPath:  "invalid_use_cases.json",
			newPath:  "use_cases1.json",
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			old := loadFixtureString(fmt.Sprintf("%s/%s", testDir, test.oldPath))
			new := loadFixtureString(fmt.Sprintf("%s/%s", testDir, test.newPath))

			assert.Equal(t, test.expected, suppressEdgeHostnameUseCases("", old, new, nil))
		})
	}
}

func TestConvertingUseCases2JSON(t *testing.T) {
	testDir := "testdata/TestResourceEdgeHostname/use_cases"
	tests := map[string]struct {
		useCases []papi.UseCase
		expected []byte
	}{
		"no use cases": {
			useCases: []papi.UseCase{},
			expected: []byte{},
		},
		"two use cases": {
			useCases: []papi.UseCase{
				{
					Option:  "BACKGROUND",
					Type:    "GLOBAL",
					UseCase: "Download_Mode",
				},
				{
					Option:  "FOREGROUND",
					Type:    "GLOBAL",
					UseCase: "Download_Mode",
				},
			},
			expected: loadFixtureBytes(fmt.Sprintf("%s/%s", testDir, "use_cases1.json")),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			useCasesJSON, err := useCases2JSON(test.useCases)
			assert.NoError(t, err)

			if len(useCasesJSON) > 0 {
				expected := new(bytes.Buffer)
				err = json.Compact(expected, test.expected)
				assert.NoError(t, err)

				actual := new(bytes.Buffer)
				err = json.Compact(actual, useCasesJSON)
				assert.NoError(t, err)

				assert.Equal(t, expected.String(), actual.String())
			} else {
				assert.Equal(t, test.expected, useCasesJSON)
			}
		})
	}
}
