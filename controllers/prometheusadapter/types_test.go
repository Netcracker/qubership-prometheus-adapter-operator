// Copyright NetCracker Technology Corporation
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

package prometheusadapter

import (
	"testing"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestNewPrometheusAdapterManager(t *testing.T) {
	// Create a fake client
	client := fake.NewClientBuilder().Build()

	// Create a logger
	logger := logr.Discard()

	// Test creating a new manager
	manager := NewPrometheusAdapterManager(client, logger)

	if manager == nil {
		t.Fatal("expected manager to be created")
	}

	if manager.client == nil {
		t.Error("expected client to be set")
	}

	// Test that the manager has a logger set
	// The logger is created with WithName in the constructor
}
