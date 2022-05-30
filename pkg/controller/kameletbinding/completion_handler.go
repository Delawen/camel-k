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
	"encoding/json"
	"fmt"

	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	"github.com/apache/camel-k/pkg/util/bindings"
	"github.com/pkg/errors"
)

func maybeCompletionHandler(compHandlConf *v1alpha1.CompletionHandlerSpec, bindingContext bindings.BindingContext) (*bindings.Binding, error) {
	var handlerBinding *bindings.Binding
	if compHandlConf != nil {
		handlerSpec, err := parseCompletionHandler(compHandlConf.RawMessage)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse error handler")
		}
		// We need to get the translated URI from any referenced resource (ie, kamelets)
		if handlerSpec.Type() == v1alpha1.CompletionHandlerTypeSink {
			handlerBinding, err = bindings.Translate(bindingContext, bindings.EndpointContext{Type: v1alpha1.EndpointTypeCompletionHandler}, *handlerSpec.Endpoint())
			if err != nil {
				return nil, errors.Wrap(err, "could not determine error handler URI")
			}
		} else {
			// Create a new binding otherwise in order to store error handler application properties
			handlerBinding = &bindings.Binding{
				ApplicationProperties: make(map[string]string),
			}
		}

		err = setCompletionHandlerConfiguration(handlerBinding, handlerSpec)
		if err != nil {
			return nil, errors.Wrap(err, "could not set integration error handler")
		}

		return handlerBinding, nil
	}
	return nil, nil
}

func parseCompletionHandler(rawMessage v1alpha1.RawMessage) (v1alpha1.CompletionHandler, error) {

	fmt.Println("completionHandler::parseCompletionHandler()")
    print("completionHandler::parseCompletionHandler()\n")
	var properties map[v1alpha1.CompletionHandlerType]v1alpha1.RawMessage
	err := json.Unmarshal(rawMessage, &properties)
	if err != nil {
		return nil, err
	}
	if len(properties) > 1 {
		return nil, errors.Errorf("You must provide just 1 error handler, provided %d", len(properties))
	}

	for _, handlValue := range properties {
		var dst v1alpha1.CompletionHandler
		dst = new(v1alpha1.CompletionHandlerSink)
		err := json.Unmarshal(handlValue, dst)

		if err != nil {
			return nil, err
		}

		return dst, nil
	}

	return nil, nil
}

func setCompletionHandlerConfiguration(completionHandlerBinding *bindings.Binding, completionHandler v1alpha1.CompletionHandler) error {

	fmt.Println("completionHandler::setCompletionHandlerConfiguration()")
    print("completionHandler::setCompletionHandlerConfiguration()\n")
	properties, err := completionHandler.Configuration()
	if err != nil {
		return err
	}
	// initialize map if not yet initialized
	if completionHandlerBinding.ApplicationProperties == nil {
		completionHandlerBinding.ApplicationProperties = make(map[string]string)
	}
	for key, value := range properties {
		completionHandlerBinding.ApplicationProperties[key] = fmt.Sprintf("%v", value)
	}
	if completionHandler.Type() == v1alpha1.CompletionHandlerTypeSink && completionHandlerBinding.URI != "" {
		completionHandlerBinding.ApplicationProperties[fmt.Sprintf("%s.deadLetterUri", v1alpha1.CompletionHandlerAppPropertiesPrefix)] = fmt.Sprintf("%v", completionHandlerBinding.URI)
	}

	return nil
}
