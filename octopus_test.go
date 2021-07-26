// Copyright 2021 Aadhav Vignesh

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package octopool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test for checking the behavior when a new octopus is created.
func TestOctopus(t *testing.T) {
	testOctopus := NewOctopus(20)

	assert.NotNil(t, testOctopus, "octopus should not be nil")
	assert.Equal(t, 20, testOctopus.workerPool.capacity, "octopus should have the capacity as mentioned (20).")
}

// Test for checking the behavior when an octopus with an invalid capacity is created.
func TestInvalidOctopus(t *testing.T) {
	testOctopus := NewOctopus(0)

	assert.NotNil(t, testOctopus, "octopus should not be nil")
	assert.Equal(t, 1000, testOctopus.workerPool.capacity, "octopus should have the default capacity (1000).")
}
