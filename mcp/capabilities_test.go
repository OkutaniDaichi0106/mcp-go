package mcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapabilities_EnableCapability(t *testing.T) {
	tests := map[string]struct {
		initial  Capabilities
		feature  string
		method   string
		expected Capabilities
	}{
		"enable capability on empty map": {
			initial:  Capabilities{},
			feature:  "feature1",
			method:   "method1",
			expected: Capabilities{"feature1": {"method1": true}},
		},
		"enable capability on existing feature": {
			initial:  Capabilities{"feature1": {"method2": true}},
			feature:  "feature1",
			method:   "method1",
			expected: Capabilities{"feature1": {"method1": true, "method2": true}},
		},
		"enable already enabled capability": {
			initial:  Capabilities{"feature1": {"method1": true}},
			feature:  "feature1",
			method:   "method1",
			expected: Capabilities{"feature1": {"method1": true}},
		},
		"enable capability on different feature": {
			initial:  Capabilities{"feature1": {"method1": true}},
			feature:  "feature2",
			method:   "method1",
			expected: Capabilities{"feature1": {"method1": true}, "feature2": {"method1": true}},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Create a copy of the initial capabilities to avoid test interference
			c := make(Capabilities)
			for k, v := range tc.initial {
				c[k] = make(map[string]bool)
				for mk, mv := range v {
					c[k][mk] = mv
				}
			}

			result := c.EnableCapability(tc.feature, tc.method)

			// Verify all expected features and methods
			assert.Equal(t, tc.expected, result, "Capabilities after EnableCapability should match expected")
		})
	}
}

func TestCapabilities_DisableCapability(t *testing.T) {
	tests := map[string]struct {
		initial  Capabilities
		feature  string
		method   string
		expected Capabilities
	}{
		"disable capability on empty map": {
			initial:  Capabilities{},
			feature:  "feature1",
			method:   "method1",
			expected: Capabilities{"feature1": {"method1": false}},
		},
		"disable capability on existing feature": {
			initial:  Capabilities{"feature1": {"method1": true, "method2": true}},
			feature:  "feature1",
			method:   "method1",
			expected: Capabilities{"feature1": {"method1": false, "method2": true}},
		},
		"disable already disabled capability": {
			initial:  Capabilities{"feature1": {"method1": false}},
			feature:  "feature1",
			method:   "method1",
			expected: Capabilities{"feature1": {"method1": false}},
		},
		"disable capability on different feature": {
			initial:  Capabilities{"feature1": {"method1": true}},
			feature:  "feature2",
			method:   "method1",
			expected: Capabilities{"feature1": {"method1": true}, "feature2": {"method1": false}},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Create a copy of the initial capabilities to avoid test interference
			c := make(Capabilities)
			for k, v := range tc.initial {
				c[k] = make(map[string]bool)
				for mk, mv := range v {
					c[k][mk] = mv
				}
			}

			result := c.DisableCapability(tc.feature, tc.method)

			// Verify the result is the same instance (method returns the receiver)
			assert.Same(t, c, result, "DisableCapability should return the receiver")

			// Verify all expected features and methods
			assert.Equal(t, tc.expected, result, "Capabilities after DisableCapability should match expected")
		})
	}
}

func TestCapabilities_Merge(t *testing.T) {
	tests := map[string]struct {
		initial  Capabilities
		other    Capabilities
		expected Capabilities
	}{
		"merge with empty map": {
			initial:  Capabilities{},
			other:    Capabilities{},
			expected: Capabilities{},
		},
		"merge empty map with non-empty map": {
			initial:  Capabilities{},
			other:    Capabilities{"feature1": {"method1": true}},
			expected: Capabilities{"feature1": {"method1": true}},
		},
		"merge non-empty map with empty map": {
			initial:  Capabilities{"feature1": {"method1": true}},
			other:    Capabilities{},
			expected: Capabilities{"feature1": {"method1": true}},
		},
		"merge with different features": {
			initial:  Capabilities{"feature1": {"method1": true}},
			other:    Capabilities{"feature2": {"method2": true}},
			expected: Capabilities{"feature1": {"method1": true}, "feature2": {"method2": true}},
		},
		"merge with overlapping features (initial takes precedence)": {
			initial:  Capabilities{"feature1": {"method1": true}},
			other:    Capabilities{"feature1": {"method2": true}},
			expected: Capabilities{"feature1": {"method1": true}},
		},
		"merge with multiple features": {
			initial:  Capabilities{"feature1": {"method1": true}, "feature3": {"method3": true}},
			other:    Capabilities{"feature2": {"method2": true}, "feature4": {"method4": true}},
			expected: Capabilities{"feature1": {"method1": true}, "feature2": {"method2": true}, "feature3": {"method3": true}, "feature4": {"method4": true}},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Create a copy of the initial capabilities to avoid test interference
			c := make(Capabilities)
			for k, v := range tc.initial {
				c[k] = make(map[string]bool)
				for mk, mv := range v {
					c[k][mk] = mv
				}
			}

			// Create a copy of the other capabilities
			other := make(Capabilities)
			for k, v := range tc.other {
				other[k] = make(map[string]bool)
				for mk, mv := range v {
					other[k][mk] = mv
				}
			}

			result := c.Merge(other)

			// Verify all expected features and methods
			assert.Equal(t, tc.expected, result, "Capabilities after Merge should match expected")
		})
	}
}

func TestCapabilities_HasFeature(t *testing.T) {
	tests := map[string]struct {
		caps     Capabilities
		feature  string
		expected bool
	}{
		"empty capabilities": {
			caps:     Capabilities{},
			feature:  "feature1",
			expected: false,
		},
		"feature exists": {
			caps:     Capabilities{"feature1": {"method1": true}},
			feature:  "feature1",
			expected: true,
		},
		"feature doesn't exist": {
			caps:     Capabilities{"feature1": {"method1": true}},
			feature:  "feature2",
			expected: false,
		},
		"feature exists with empty methods": {
			caps:     Capabilities{"feature1": {}},
			feature:  "feature1",
			expected: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Create a copy of the capabilities to avoid test interference
			c := make(Capabilities)
			for k, v := range tc.caps {
				c[k] = make(map[string]bool)
				for mk, mv := range v {
					c[k][mk] = mv
				}
			}

			result := c.HasFeature(tc.feature)

			assert.Equal(t, tc.expected, result, "HasFeature result should match expected")
		})
	}
}

func TestCapabilities_HasCapability(t *testing.T) {
	tests := map[string]struct {
		caps     Capabilities
		feature  string
		method   string
		expected bool
	}{
		"empty capabilities": {
			caps:     Capabilities{},
			feature:  "feature1",
			method:   "method1",
			expected: false,
		},
		"feature exists but method doesn't": {
			caps:     Capabilities{"feature1": {"method2": true}},
			feature:  "feature1",
			method:   "method1",
			expected: false,
		},
		"feature doesn't exist": {
			caps:     Capabilities{"feature2": {"method1": true}},
			feature:  "feature1",
			method:   "method1",
			expected: false,
		},
		"capability exists and is enabled": {
			caps:     Capabilities{"feature1": {"method1": true}},
			feature:  "feature1",
			method:   "method1",
			expected: true,
		},
		"capability exists but is disabled": {
			caps:     Capabilities{"feature1": {"method1": false}},
			feature:  "feature1",
			method:   "method1",
			expected: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Create a copy of the capabilities to avoid test interference
			c := make(Capabilities)
			for k, v := range tc.caps {
				c[k] = make(map[string]bool)
				for mk, mv := range v {
					c[k][mk] = mv
				}
			}

			result := c.HasCapability(tc.feature, tc.method)

			assert.Equal(t, tc.expected, result, "HasCapability result should match expected")
		})
	}
}
