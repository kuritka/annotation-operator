package utils

/*
Copyright 2022 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapHasOnlyKeys(t *testing.T) {
	var tests = []struct {
		name     string
		m        map[string]string
		a        []string
		expected bool
	}{
		{
			name:     "valid",
			m:        map[string]string{"a": "1", "b": "3", "c": "1"},
			a:        []string{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "repeating items",
			m:        map[string]string{"a": "1", "b": "3", "c": "1"},
			a:        []string{"a", "a", "c"},
			expected: false,
		},
		{
			name:     "less items",
			m:        map[string]string{"a": "1", "b": "3", "c": "1"},
			a:        []string{"a", "c"},
			expected: false,
		},
		{
			name:     "empty",
			m:        map[string]string{},
			a:        []string{},
			expected: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b := MapHasOnlyKeys(test.m, test.a...)
			assert.Equal(t, test.expected, b)
		})
	}
}
