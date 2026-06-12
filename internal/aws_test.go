package internal

import (
	"net"
	"testing"
)

func TestIsPublicIPv4CIDR(t *testing.T) {
	tests := []struct {
		cidr     string
		expected bool
	}{
		// Private Ranges
		{"10.0.0.0/8", false},
		{"10.255.255.255/32", false},
		{"172.16.0.0/12", false},
		{"172.31.255.255/32", false},
		{"192.168.0.0/16", false},
		{"127.0.0.0/8", false},
		{"169.254.0.0/16", false},
		{"224.0.0.0/4", false}, // Multicast

		// Public Ranges
		{"0.0.0.0/0", true},
		{"8.8.8.8/32", true},
		{"11.0.0.0/8", true},
		{"192.167.0.0/16", true},
		{"172.15.255.255/32", true},
		{"172.32.0.0/16", true},

		// Invalid CIDR
		{"invalid-cidr", false},
	}

	for _, tt := range tests {
		t.Run(tt.cidr, func(t *testing.T) {
			result := isPublicIPv4CIDR(tt.cidr)
			if result != tt.expected {
				t.Errorf("isPublicIPv4CIDR(%q) = %v, expected %v", tt.cidr, result, tt.expected)
			}
		})
	}
}

func TestIsPublicIPv6CIDR(t *testing.T) {
	tests := []struct {
		cidr     string
		expected bool
	}{
		// Private Ranges
		{"fc00::/7", false},
		{"fd00::/8", false},
		{"fe80::/10", false},
		{"::1/128", false},
		{"ff00::/8", false}, // Multicast

		// Public Ranges
		{"::/0", true},
		{"2001:db8::/32", true}, // Documentation (technically public/non-private in the code)
		{"2001:4860:4860::8888/128", true},

		// Invalid CIDR
		{"invalid-cidr", false},
	}

	for _, tt := range tests {
		t.Run(tt.cidr, func(t *testing.T) {
			result := isPublicIPv6CIDR(tt.cidr)
			if result != tt.expected {
				t.Errorf("isPublicIPv6CIDR(%q) = %v, expected %v", tt.cidr, result, tt.expected)
			}
		})
	}
}

func TestIsNetworkSubsetOf(t *testing.T) {
	tests := []struct {
		cidr1    string
		cidr2    string
		expected bool
	}{
		{"10.0.0.0/8", "10.0.0.0/8", true},
		{"10.0.0.0/9", "10.0.0.0/8", true},
		{"10.128.0.0/9", "10.0.0.0/8", true},
		{"10.0.0.0/8", "10.0.0.0/9", false}, // Less specific is not a subset of more specific
		{"172.16.0.0/12", "10.0.0.0/8", false},
		{"fc00::/8", "fc00::/7", true},
		{"fc00::/7", "fc00::/8", false},
	}

	for _, tt := range tests {
		t.Run(tt.cidr1+"_in_"+tt.cidr2, func(t *testing.T) {
			_, net1, _ := net.ParseCIDR(tt.cidr1)
			_, net2, _ := net.ParseCIDR(tt.cidr2)
			result := isNetworkSubsetOf(net1, net2)
			if result != tt.expected {
				t.Errorf("isNetworkSubsetOf(%s, %s) = %v, expected %v", tt.cidr1, tt.cidr2, result, tt.expected)
			}
		})
	}
}
