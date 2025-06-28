package internal

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type SecurityGroupResult struct {
	Region          string        `json:"region"`
	SecurityGroupID string        `json:"security_group_id"`
	GroupName       string        `json:"group_name"`
	Description     string        `json:"description"`
	VpcID           string        `json:"vpc_id"`
	RiskyRules      []InboundRule `json:"risky_rules"`
}

type InboundRule struct {
	FromPort   *int32   `json:"from_port"`
	ToPort     *int32   `json:"to_port"`
	Protocol   string   `json:"protocol"`
	CidrBlocks []string `json:"cidr_blocks"`
}

type AWSClient struct {
	ec2Clients map[string]*ec2.Client
	regions    []string
}

func NewAWSClient(ctx context.Context, targetRegion string) (*AWSClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	client := &AWSClient{
		ec2Clients: make(map[string]*ec2.Client),
	}

	switch targetRegion {
	case "all":
		regions, err := client.getAllRegions(ctx, &cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to get regions: %w", err)
		}
		client.regions = regions
	case "":
		client.regions = []string{cfg.Region}
	default:
		client.regions = []string{targetRegion}
	}

	for _, region := range client.regions {
		regionalCfg := cfg.Copy()
		regionalCfg.Region = region
		client.ec2Clients[region] = ec2.NewFromConfig(regionalCfg)
	}

	return client, nil
}

func (c *AWSClient) getAllRegions(ctx context.Context, cfg *aws.Config) ([]string, error) {
	ec2Client := ec2.NewFromConfig(*cfg)

	result, err := ec2Client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, err
	}

	var regions []string
	for _, region := range result.Regions {
		if region.RegionName != nil {
			regions = append(regions, *region.RegionName)
		}
	}

	return regions, nil
}

func (c *AWSClient) ScanSecurityGroups(ctx context.Context, excludeIDs map[string]bool) ([]SecurityGroupResult, error) {
	var results []SecurityGroupResult

	for _, region := range c.regions {
		client := c.ec2Clients[region]

		input := &ec2.DescribeSecurityGroupsInput{}
		paginator := ec2.NewDescribeSecurityGroupsPaginator(client, input)

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to describe security groups in region %s: %w", region, err)
			}

			for _, sg := range page.SecurityGroups {
				if sg.GroupId == nil {
					continue
				}

				groupID := *sg.GroupId
				if excludeIDs[groupID] {
					continue
				}

				riskyRules := c.findRiskyInboundRules(sg.IpPermissions)
				if len(riskyRules) > 0 {
					result := SecurityGroupResult{
						Region:          region,
						SecurityGroupID: groupID,
						GroupName:       aws.ToString(sg.GroupName),
						Description:     aws.ToString(sg.Description),
						VpcID:           aws.ToString(sg.VpcId),
						RiskyRules:      riskyRules,
					}
					results = append(results, result)
				}
			}
		}
	}

	return results, nil
}

func (c *AWSClient) findRiskyInboundRules(permissions []types.IpPermission) []InboundRule {
	var riskyRules []InboundRule

	for i := range permissions {
		perm := &permissions[i]
		var riskyCidrBlocks []string
		hasPublicAccess := false

		// Check IPv4 ranges
		for _, ipRange := range perm.IpRanges {
			if ipRange.CidrIp != nil {
				cidr := *ipRange.CidrIp
				if isPublicIPv4CIDR(cidr) {
					riskyCidrBlocks = append(riskyCidrBlocks, cidr)
					hasPublicAccess = true
				}
			}
		}

		// Check IPv6 ranges
		for _, ipv6Range := range perm.Ipv6Ranges {
			if ipv6Range.CidrIpv6 != nil {
				cidr := *ipv6Range.CidrIpv6
				if isPublicIPv6CIDR(cidr) {
					riskyCidrBlocks = append(riskyCidrBlocks, cidr)
					hasPublicAccess = true
				}
			}
		}

		// Check prefix list IDs (for managed prefix lists that might include global ranges)
		for _, prefixList := range perm.PrefixListIds {
			if prefixList.PrefixListId != nil {
				riskyCidrBlocks = append(riskyCidrBlocks, "prefix:"+*prefixList.PrefixListId)
				hasPublicAccess = true // Assume prefix lists might contain public IPs
			}
		}

		if hasPublicAccess {
			rule := InboundRule{
				FromPort:   perm.FromPort,
				ToPort:     perm.ToPort,
				Protocol:   aws.ToString(perm.IpProtocol),
				CidrBlocks: riskyCidrBlocks,
			}
			riskyRules = append(riskyRules, rule)
		}
	}

	return riskyRules
}

// isPublicIPv4CIDR checks if the CIDR contains public IPv4 addresses
func isPublicIPv4CIDR(cidr string) bool {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}

	// Check if the network contains any public IP
	return containsPublicIPv4(ipNet)
}

// isPublicIPv6CIDR checks if the CIDR contains public IPv6 addresses
func isPublicIPv6CIDR(cidr string) bool {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}

	// Check if the network contains any public IP
	return containsPublicIPv6(ipNet)
}

// containsPublicIPv4 checks if the network contains public IPv4 addresses
func containsPublicIPv4(ipNet *net.IPNet) bool {
	// Private IPv4 ranges:
	// 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
	// Loopback: 127.0.0.0/8
	// Link-local: 169.254.0.0/16
	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"224.0.0.0/4", // Multicast
	}

	for _, privateRange := range privateRanges {
		_, privateNet, _ := net.ParseCIDR(privateRange)
		if isNetworkSubsetOf(ipNet, privateNet) {
			return false // Entirely within private range
		}
	}

	return true // Contains public IPs
}

// containsPublicIPv6 checks if the network contains public IPv6 addresses
func containsPublicIPv6(ipNet *net.IPNet) bool {
	// Private IPv6 ranges:
	// fc00::/7 (unique local), fe80::/10 (link-local), ::1/128 (loopback)
	privateRanges := []string{
		"fc00::/7",  // Unique local
		"fe80::/10", // Link-local
		"::1/128",   // Loopback
		"ff00::/8",  // Multicast
	}

	for _, privateRange := range privateRanges {
		_, privateNet, _ := net.ParseCIDR(privateRange)
		if isNetworkSubsetOf(ipNet, privateNet) {
			return false // Entirely within private range
		}
	}

	return true // Contains public IPs
}

// isNetworkSubsetOf checks if network1 is entirely contained within network2
func isNetworkSubsetOf(network1, network2 *net.IPNet) bool {
	return network2.Contains(network1.IP) &&
		network1.Mask != nil && network2.Mask != nil &&
		compareMask(network1.Mask, network2.Mask) >= 0
}

// compareMask compares two network masks, returns positive if mask1 is more specific than mask2
func compareMask(mask1, mask2 net.IPMask) int {
	ones1, _ := mask1.Size()
	ones2, _ := mask2.Size()
	return ones1 - ones2
}

func LoadExcludeList(filename string) (map[string]bool, error) {
	excludeMap := make(map[string]bool)

	if filename == "" {
		return excludeMap, nil
	}

	content, err := readFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read exclude file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			excludeMap[line] = true
		}
	}

	return excludeMap, nil
}
