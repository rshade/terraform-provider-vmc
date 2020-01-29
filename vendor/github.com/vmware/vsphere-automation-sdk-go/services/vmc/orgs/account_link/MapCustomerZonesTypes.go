/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

// Code generated. DO NOT EDIT.

/*
 * Data type definitions file for service: MapCustomerZones.
 * Includes binding types of a structures and enumerations defined in the service.
 * Shared by client-side stubs and server-side skeletons to ensure type
 * compatibility.
 */

package account_link

import (
	"reflect"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol"
)





func mapCustomerZonesPostInputType() bindings.StructType {
	fields := make(map[string]bindings.BindingType)
	fieldNameMap := make(map[string]string)
	fields["org"] = bindings.NewStringType()
	fields["map_zones_request"] = bindings.NewReferenceType(model.MapZonesRequestBindingType)
	fieldNameMap["org"] = "Org"
	fieldNameMap["map_zones_request"] = "MapZonesRequest"
	var validators = []bindings.Validator{}
	return bindings.NewStructType("operation-input", fields, reflect.TypeOf(data.StructValue{}), fieldNameMap, validators)
}

func mapCustomerZonesPostOutputType() bindings.BindingType {
	return bindings.NewReferenceType(model.TaskBindingType)
}

func mapCustomerZonesPostRestMetadata() protocol.OperationRestMetadata {
	fields := map[string]bindings.BindingType{}
	fieldNameMap := map[string]string{}
	paramsTypeMap := map[string]bindings.BindingType{}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	fields["org"] = bindings.NewStringType()
	fields["map_zones_request"] = bindings.NewReferenceType(model.MapZonesRequestBindingType)
	fieldNameMap["org"] = "Org"
	fieldNameMap["map_zones_request"] = "MapZonesRequest"
	paramsTypeMap["org"] = bindings.NewStringType()
	paramsTypeMap["map_zones_request"] = bindings.NewReferenceType(model.MapZonesRequestBindingType)
	paramsTypeMap["org"] = bindings.NewStringType()
	pathParams["org"] = "org"
	resultHeaders := map[string]string{}
	errorHeaders := map[string]string{}
	errorHeaders["Unauthenticated.challenge"] = "WWW-Authenticate"
	return protocol.NewOperationRestMetadata(
		fields,
		fieldNameMap,
		paramsTypeMap,
		pathParams,
		queryParams,
		headerParams,
		"",
		"map_zones_request",
		"POST",
		"/vmc/api/orgs/{org}/account-link/map-customer-zones",
		resultHeaders,
		200,
		errorHeaders,
		map[string]int{"Unauthenticated": 401,"Unauthorized": 403})
}


