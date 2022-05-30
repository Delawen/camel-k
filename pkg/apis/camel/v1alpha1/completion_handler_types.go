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

const (
	// CompletionHandlerRefName the reference name to use when looking for an completion handler
	CompletionHandlerRefName = "camel.k.completionHandler.ref"
	// CompletionHandlerRefDefaultName the default name of the completion handler
	CompletionHandlerRefDefaultName = "defaultCompletionHandler"
	// CompletionHandlerAppPropertiesPrefix the prefix used for the completion handler bean
	CompletionHandlerAppPropertiesPrefix = "camel.beans.defaultCompletionHandler"
)

// CompletionHandlerSpec represents an unstructured object for a completion handler
type CompletionHandlerSpec struct {
	RawMessage `json:",omitempty"`
}

// CompletionHandlerParameters represent an unstructured object for completion handler parameters
type CompletionHandlerParameters struct {
	RawMessage `json:",omitempty"`
}

// CompletionHandlerType a type of completion handler (ie, sink)
type CompletionHandlerType string

const (
	completionHandlerTypeBase CompletionHandlerType = ""
	// ErrorHandlerTypeNone used to ignore any error event
	CompletionHandlerTypeNone CompletionHandlerType = "none"
    CompletionHandlerTypeLog CompletionHandlerType = "log"
	// CompletionHandlerTypeSink used to send the event to a further sink (for future processing).
	CompletionHandlerTypeSink CompletionHandlerType = "sink"
)
