// Copyright (c) 2014 Hewlett-Packard Development Company, L.P.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package network_test

import (
	"encoding/json"
	"errors"
	"testing"

	"git.openstack.org/stackforge/golang-client.git/identity/common"
	network "git.openstack.org/stackforge/golang-client.git/network/v2"
	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

var tokn = "eaaafd18-0fed-4b3a-81b4-663c99ec1cbb"
var subnets = []string{"10.3.5.2", "12.34.1.4"}
var sampleNetworkResponse = network.Response{
	ID:                  "16470140hb",
	Name:                "networkName",
	Status:              "active",
	Subnets:             subnets,
	TenantID:            "tenantID",
	RouterExternal:      true,
	AdminStateUp:        false,
	Shared:              true,
	PortSecurityEnabled: false}

func TestGetNetworks(t *testing.T) {
	mockResponseObject := networksContainer{Networks: []network.Response{sampleNetworkResponse}}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, "/networks")
	defer apiServer.Close()

	networkService := CreateNetworkService(apiServer.URL)
	networks, err := networkService.Networks()
	if err != nil {
		t.Error(err)
	}

	if len(networks) != 1 {
		t.Error(errors.New("Error: Expected 2 networks to be listed"))
	}
	testUtil.Equals(t, sampleNetworkResponse, networks[0])
}

func TestQueryNetworks(t *testing.T) {
	mockResponseObject := networksContainer{Networks: []network.Response{sampleNetworkResponse}}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, "/networks?admin_state_up=true&name=foo&router%3Aexternal=true&shared=true&status=ACTIVE")
	defer apiServer.Close()

	networkService := CreateNetworkService(apiServer.URL)
	qp := network.QueryParameters{
		Name:           "foo",
		Status:         "ACTIVE",
		RouterExternal: true,
		AdminStateUp:   true,
		Shared:         true,
	}

	networks, err := networkService.QueryNetworks(qp)
	if err != nil {
		t.Error(err)
	}

	if len(networks) != 1 {
		t.Error(errors.New("Error: Expected 2 networks to be listed"))
	}
	testUtil.Equals(t, sampleNetworkResponse, networks[0])
}

func TestGetNetwork(t *testing.T) {
	mockResponseObject := networkContainer{Network: sampleNetworkResponse}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, "/networks/5270u2tg0")
	defer apiServer.Close()

	networkService := CreateNetworkService(apiServer.URL)
	network, err := networkService.Network("5270u2tg0")
	if err != nil {
		t.Error(err)
	}

	testUtil.Equals(t, sampleNetworkResponse, network)
}

func TestGetNetworkIDsByName(t *testing.T) {
	mockResponseObject := networksContainer{Networks: []network.Response{sampleNetworkResponse}}
	apiServer := testUtil.CreateGetJSONTestRequestServerWithMockObject(t, tokn, mockResponseObject, "/networks?fields=id&name=networkName")
	defer apiServer.Close()

	networkService := CreateNetworkService(apiServer.URL)
	networkIDs, err := networkService.NetworkIDsByName(sampleNetworkResponse.Name)
	if err != nil {
		t.Error(err)
	}

	testUtil.Equals(t, []string{sampleNetworkResponse.ID}, networkIDs)
}

func TestDeleteNetwork(t *testing.T) {
	name := "networkName"
	apiServer := testUtil.CreateDeleteTestRequestServer(t, tokn, "/networks/"+name)
	defer apiServer.Close()

	networkService := CreateNetworkService(apiServer.URL)
	err := networkService.DeleteNetwork(name)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateNetwork(t *testing.T) {
	mockResponse, _ := json.Marshal(networkContainer{sampleNetworkResponse})
	apiServer := testUtil.CreatePostJSONTestRequestServer(t, tokn, string(mockResponse), "/networks",
		`{"network":{"name":"networkName","admin_state_up":false,"shared":true,"tenant_id":"tenantId"}}`)
	defer apiServer.Close()

	networkService := CreateNetworkService(apiServer.URL)
	params := network.CreateNetworkParameters{
		AdminStateUp: false,
		Shared:       true,
		Name:         "networkName",
		TenantID:     "tenantId"}
	actualNetwork, err := networkService.CreateNetwork(params)
	if err != nil {
		t.Error(err)
	}

	testUtil.Equals(t, sampleNetworkResponse, actualNetwork)
}

func CreateNetworkService(url string) network.Service {
	return network.NewService(common.SimpleAuthenticator{Token: tokn, ServiceURL: url})
}

type networksContainer struct {
	Networks []network.Response `json:"networks"`
}

type networkContainer struct {
	Network network.Response `json:"network"`
}
