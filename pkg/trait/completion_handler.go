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

package trait

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"

	v1 "github.com/apache/camel-k/pkg/apis/camel/v1"
	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	"github.com/apache/camel-k/pkg/util"
)

// The completion-handler is a platform trait used to inject On Completion
// Handler sources into the integration runtime.
//
// +camel-k:trait=completion-handler.
type completionHandlerTrait struct {
	BaseTrait `property:",squash"`
	// The completion handler ref name provided or found in application properties
	CompletionHandlerRef string `property:"ref" json:"ref,omitempty"`
}

func newCompletionHandlerTrait() Trait {
	return &completionHandlerTrait{
		// NOTE: Must run before dependency trait
		BaseTrait: NewBaseTrait("completion-handler", 470),
	}
}

// IsPlatformTrait overrides base class method.
func (t *completionHandlerTrait) IsPlatformTrait() bool {
	return true
}

func (t *completionHandlerTrait) Configure(e *Environment) (bool, error) {
	if IsFalse(t.Enabled) {
		return false, nil
	}

	if !e.IntegrationInPhase(v1.IntegrationPhaseInitialization) && !e.IntegrationInRunningPhases() {
		return false, nil
	}

	if t.CompletionHandlerRef == "" {
		t.CompletionHandlerRef = e.Integration.Spec.GetConfigurationProperty(v1alpha1.CompletionHandlerRefName)
	}

	return t.CompletionHandlerRef != "", nil
}

func (t *completionHandlerTrait) Apply(e *Environment) error {
	if e.IntegrationInPhase(v1.IntegrationPhaseInitialization) {
		// If the user configure directly the URI, we need to auto-discover the underlying component
		// and add the related dependency
		defaultCompletionHandlerURI := e.Integration.Spec.GetConfigurationProperty(
			fmt.Sprintf("%s.deadLetterUri", v1alpha1.CompletionHandlerAppPropertiesPrefix))
		if defaultCompletionHandlerURI != "" && !strings.HasPrefix(defaultCompletionHandlerURI, "kamelet:") {
			t.addCompletionHandlerDependencies(e, defaultCompletionHandlerURI)
		}

		return t.addCompletionHandlerAsSource(e)
	}
	return nil
}

func (t *completionHandlerTrait) addCompletionHandlerDependencies(e *Environment, uri string) {
	candidateComp, scheme := e.CamelCatalog.DecodeComponent(uri)
	if candidateComp != nil {
		util.StringSliceUniqueAdd(&e.Integration.Status.Dependencies, candidateComp.GetDependencyID())
		if scheme != nil {
			for _, dep := range candidateComp.GetProducerDependencyIDs(scheme.ID) {
				util.StringSliceUniqueAdd(&e.Integration.Status.Dependencies, dep)
			}
		}
	}
}

func (t *completionHandlerTrait) addCompletionHandlerAsSource(e *Environment) error {
//    if (t == nil) {
//      return nil
//    }

	flowCompletionHandler := map[string]interface{}{
		"completion-handler": map[string]string{
			"ref": t.CompletionHandlerRef,
		},
	}
	encodedFlowCompletionHandler, err := yaml.Marshal([]map[string]interface{}{flowCompletionHandler})
	if err != nil {
		return err
	}
	completionHandlerSource := v1.SourceSpec{
		DataSpec: v1.DataSpec{
			Name:    "camel-k-embedded-completion-handler.yaml",
			Content: string(encodedFlowCompletionHandler),
		},
		Language: v1.LanguageYaml,
		Type:     v1.SourceTypeCompletionHandler,
	}

	e.Integration.Status.AddOrReplaceGeneratedSources(completionHandlerSource)

	return nil
}
