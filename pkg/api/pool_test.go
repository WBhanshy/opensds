// Copyright 2017 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/astaxie/beego"
	"github.com/opensds/opensds/pkg/db"
	dbtest "github.com/opensds/opensds/pkg/db/testing"
	"github.com/opensds/opensds/pkg/model"
)

func init() {
	var poolPortal PoolPortal
	beego.Router("/v1alpha/block/pools", &poolPortal, "get:ListPools")
	beego.Router("/v1alpha/block/pools/:poolId", &poolPortal, "get:GetPool")
}

var (
	fakePool = &model.StoragePoolSpec{
		BaseModel: &model.BaseModel{
			Id:        "f4486139-78d5-462d-a7b9-fdaf6c797e1b",
			CreatedAt: "2017-10-24T15:04:05",
		},
		Name:             "fakePool",
		Description:      "fake pool for testing",
		Status:           "available",
		AvailabilityZone: "unknown",
		TotalCapacity:    99999,
		FreeCapacity:     6999,
		DockId:           "ccac4f33-e603-425a-8813-371bbe10566e",
		Parameters: map[string]interface{}{
			"key1": "val1",
			"key2": "val2",
			"key3": map[string]string{
				"subKey1": "subVal1",
				"subKey2": "subVal2",
			},
		},
	}
	fakePools = []*model.StoragePoolSpec{fakePool}
)

func TestListPools(t *testing.T) {

	mockClient := new(dbtest.MockClient)
	mockClient.On("ListPools").Return(fakePools, nil)
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1alpha/block/pools", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output []model.StoragePoolSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `[
		{
			"id": "f4486139-78d5-462d-a7b9-fdaf6c797e1b",
			"name": "fakePool",
			"description": "fake pool for testing",
			"createdAt": "2017-10-24T15:04:05",
			"updatedAt": "",
			"status": "available",
			"availabilityZone": "unknown",
			"totalCapacity": 99999,
			"freeCapacity": 6999,
			"dockId": "ccac4f33-e603-425a-8813-371bbe10566e",
			"extras": {
				"key1": "val1",
				"key2": "val2",
				"key3": {
					"subKey1": "subVal1",
					"subKey2": "subVal2"
				}
			}	
		}		
	]`

	var expected []model.StoragePoolSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestListPoolsWithInternalError(t *testing.T) {

	mockClient := new(dbtest.MockClient)
	mockClient.On("ListPools").Return(nil, errors.New("db error"))
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1alpha/block/pools", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 500 {
		t.Errorf("Expected 500, actual %v", w.Code)
	}
}

func TestGetPool(t *testing.T) {

	mockClient := new(dbtest.MockClient)
	mockClient.On("GetPool", "f4486139-78d5-462d-a7b9-fdaf6c797e1b").Return(fakePool, nil)
	db.C = mockClient

	r, _ := http.NewRequest("GET", "/v1alpha/block/pools/f4486139-78d5-462d-a7b9-fdaf6c797e1b", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	var output model.StoragePoolSpec
	json.Unmarshal(w.Body.Bytes(), &output)

	expectedJson := `
		{
			"id": "f4486139-78d5-462d-a7b9-fdaf6c797e1b",
			"name": "fakePool",
			"description": "fake pool for testing",
			"createdAt": "2017-10-24T15:04:05",
			"updatedAt": "",
			"status": "available",
			"availabilityZone": "unknown",
			"totalCapacity": 99999,
			"freeCapacity": 6999,
			"dockId": "ccac4f33-e603-425a-8813-371bbe10566e",
			"extras": {
				"key1": "val1",
				"key2": "val2",
				"key3": {
					"subKey1": "subVal1",
					"subKey2": "subVal2"
				}
			}	
		}`

	var expected model.StoragePoolSpec
	json.Unmarshal([]byte(expectedJson), &expected)

	if w.Code != 200 {
		t.Errorf("Expected 200, actual %v", w.Code)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("Expected %v, actual %v", expected, output)
	}
}

func TestGetPoolWithBadRequest(t *testing.T) {

	mockClient := new(dbtest.MockClient)
	mockClient.On("GetPool", "f4486139-78d5-462d-a7b9-fdaf6c797e1b").Return(
		nil, errors.New("db error"))
	db.C = mockClient

	r, _ := http.NewRequest("GET",
		"/v1alpha/block/pools/f4486139-78d5-462d-a7b9-fdaf6c797e1b", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Errorf("Expected 400, actual %v", w.Code)
	}
}
