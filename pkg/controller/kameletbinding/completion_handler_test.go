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

package kameletbinding

import (
	"fmt"
	"testing"

	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestParseCompletionHandlerSinkDoesSucceed(t *testing.T) {
	fmt.Println("Test")
	sinkCompletionHandler, err := parseCompletionHandler(
		[]byte(`{"sink": {"endpoint": {"uri": "someUri"}}}`),
	)
	assert.Nil(t, err)
	assert.NotNil(t, sinkCompletionHandler)
	assert.Equal(t, v1alpha1.CompletionHandlerTypeSink, sinkCompletionHandler.Type())
	assert.Equal(t, "someUri", *sinkCompletionHandler.Endpoint().URI)
	parameters, err := sinkCompletionHandler.Configuration()
	assert.Nil(t, err)
	assert.Equal(t, "#class:org.apache.camel.builder.DeadLetterChannelBuilder", parameters[v1alpha1.CompletionHandlerAppPropertiesPrefix])
	assert.Equal(t, v1alpha1.CompletionHandlerRefDefaultName, parameters[v1alpha1.CompletionHandlerRefName])
}

func TestParseCompletionHandlerSinkWithParametersDoesSucceed(t *testing.T) {
	sinkCompletionHandler, err := parseCompletionHandler(
		[]byte(`{
			"sink": {
				"endpoint": {
					"uri": "someUri"
					}, 
				"parameters": 
					{"param1": "value1", "param2": "value2"}
			}
		}`),
	)
	assert.Nil(t, err)
	assert.NotNil(t, sinkCompletionHandler)
	assert.Equal(t, v1alpha1.CompletionHandlerTypeSink, sinkCompletionHandler.Type())
	assert.Equal(t, "someUri", *sinkCompletionHandler.Endpoint().URI)
	parameters, err := sinkCompletionHandler.Configuration()
	assert.Nil(t, err)
	assert.Equal(t, "#class:org.apache.camel.builder.DeadLetterChannelBuilder", parameters[v1alpha1.CompletionHandlerAppPropertiesPrefix])
	assert.Equal(t, v1alpha1.CompletionHandlerRefDefaultName, parameters[v1alpha1.CompletionHandlerRefName])
	assert.Equal(t, "value1", parameters["camel.beans.defaultCompletionHandler.param1"])
	assert.Equal(t, "value2", parameters["camel.beans.defaultCompletionHandler.param2"])
}
