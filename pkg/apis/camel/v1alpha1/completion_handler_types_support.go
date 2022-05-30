/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"encoding/json"
)

// +kubebuilder:object:generate=false

// CompletionHandler is a generic interface that represent any type of error handler specification
type CompletionHandler interface {
	Type() CompletionHandlerType
	Endpoint() *Endpoint
	Configuration() (map[string]interface{}, error)
}

// baseCompletionHandler is the base used for the Error Handler hierarchy
type baseCompletionHandler struct {
}

// Type --
func (e baseCompletionHandler) Type() CompletionHandlerType {
	return completionHandlerTypeBase
}

// Endpoint --
func (e baseCompletionHandler) Endpoint() *Endpoint {
	return nil
}

// Configuration --
func (e baseCompletionHandler) Configuration() (map[string]interface{}, error) {
	return nil, nil
}

// CompletionHandlerNone --
type CompletionHandlerNone struct {
	baseCompletionHandler
}

// Type --
func (e CompletionHandlerNone) Type() CompletionHandlerType {
	return CompletionHandlerTypeNone
}

// Configuration --
func (e CompletionHandlerNone) Configuration() (map[string]interface{}, error) {
	return map[string]interface{}{
		CompletionHandlerAppPropertiesPrefix: "#class:org.apache.camel.builder.NoCompletionHandlerBuilder",
		CompletionHandlerRefName:             CompletionHandlerRefDefaultName,
	}, nil
}

// CompletionHandlerLog represent a default (log) error handler type
type CompletionHandlerLog struct {
	CompletionHandlerNone
	Parameters *CompletionHandlerParameters `json:"parameters,omitempty"`
}

// Type --
func (e CompletionHandlerLog) Type() CompletionHandlerType {
	return CompletionHandlerTypeLog
}

// Configuration --
func (e CompletionHandlerLog) Configuration() (map[string]interface{}, error) {
	properties, err := e.CompletionHandlerNone.Configuration()
	if err != nil {
		return nil, err
	}
	properties[CompletionHandlerAppPropertiesPrefix] = "#class:org.apache.camel.builder.DefaultCompletionHandlerBuilder"

	if e.Parameters != nil {
		var parameters map[string]interface{}
		err := json.Unmarshal(e.Parameters.RawMessage, &parameters)
		if err != nil {
			return nil, err
		}
		for key, value := range parameters {
			properties[CompletionHandlerAppPropertiesPrefix+"."+key] = value
		}
	}

	return properties, nil
}

// CompletionHandlerSink represents a sink error handler type which behave like a dead letter channel
type CompletionHandlerSink struct {
	CompletionHandlerLog
	DLCEndpoint *Endpoint `json:"endpoint,omitempty"`
}

// Type --
func (e CompletionHandlerSink) Type() CompletionHandlerType {
	return CompletionHandlerTypeSink
}

// Endpoint --
func (e CompletionHandlerSink) Endpoint() *Endpoint {
	return e.DLCEndpoint
}

// Configuration --
func (e CompletionHandlerSink) Configuration() (map[string]interface{}, error) {
	properties, err := e.CompletionHandlerLog.Configuration()
	if err != nil {
		return nil, err
	}
	properties[CompletionHandlerAppPropertiesPrefix] = "#class:org.apache.camel.builder.DeadLetterChannelBuilder"

	return properties, err
}
