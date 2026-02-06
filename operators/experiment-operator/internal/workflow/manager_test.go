package workflow

import (
	"testing"
)

func TestIsTerminal(t *testing.T) {
	tests := []struct {
		phase    string
		expected bool
	}{
		{"Succeeded", true},
		{"Failed", true},
		{"Error", true},
		{"Running", false},
		{"Pending", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.phase, func(t *testing.T) {
			if got := IsTerminal(tt.phase); got != tt.expected {
				t.Errorf("IsTerminal(%q) = %v, want %v", tt.phase, got, tt.expected)
			}
		})
	}
}

func TestIsSucceeded(t *testing.T) {
	tests := []struct {
		phase    string
		expected bool
	}{
		{"Succeeded", true},
		{"Failed", false},
		{"Error", false},
		{"Running", false},
		{"Pending", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.phase, func(t *testing.T) {
			if got := IsSucceeded(tt.phase); got != tt.expected {
				t.Errorf("IsSucceeded(%q) = %v, want %v", tt.phase, got, tt.expected)
			}
		})
	}
}

func TestNewManager_Defaults(t *testing.T) {
	m := NewManager(nil)
	if m.Namespace != DefaultNamespace {
		t.Errorf("expected namespace %s, got %s", DefaultNamespace, m.Namespace)
	}
	if DefaultNamespace != "argo" {
		t.Errorf("expected default namespace to be 'argo', got %s", DefaultNamespace)
	}
}
