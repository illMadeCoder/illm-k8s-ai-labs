package crossplane

import (
	"testing"
	"time"
)

func TestCalculateTTL(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		created  time.Time
		ttlDays  int
		expected time.Time
	}{
		{"1 day", now, 1, now.Add(24 * time.Hour)},
		{"7 days", now, 7, now.Add(7 * 24 * time.Hour)},
		{"0 days defaults to 1 day", now, 0, now.Add(24 * time.Hour)},
		{"negative defaults to 1 day", now, -1, now.Add(24 * time.Hour)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateTTL(tt.created, tt.ttlDays)
			if !got.Equal(tt.expected) {
				t.Errorf("CalculateTTL(%v, %d) = %v, want %v", tt.created, tt.ttlDays, got, tt.expected)
			}
		})
	}
}

func TestShouldDeleteCluster(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		created  time.Time
		ttlDays  int
		expected bool
	}{
		{
			name:     "not expired (created just now, 1 day TTL)",
			created:  now,
			ttlDays:  1,
			expected: false,
		},
		{
			name:     "expired (created 2 days ago, 1 day TTL)",
			created:  now.Add(-48 * time.Hour),
			ttlDays:  1,
			expected: true,
		},
		{
			name:     "not expired (created 6 days ago, 7 day TTL)",
			created:  now.Add(-6 * 24 * time.Hour),
			ttlDays:  7,
			expected: false,
		},
		{
			name:     "expired (created 8 days ago, 7 day TTL)",
			created:  now.Add(-8 * 24 * time.Hour),
			ttlDays:  7,
			expected: true,
		},
		{
			name:     "0 ttlDays defaults to 1 day - not expired",
			created:  now,
			ttlDays:  0,
			expected: false,
		},
		{
			name:     "0 ttlDays defaults to 1 day - expired",
			created:  now.Add(-48 * time.Hour),
			ttlDays:  0,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShouldDeleteCluster(tt.created, tt.ttlDays); got != tt.expected {
				t.Errorf("ShouldDeleteCluster(%v, %d) = %v, want %v", tt.created, tt.ttlDays, got, tt.expected)
			}
		})
	}
}
