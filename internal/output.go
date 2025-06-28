package internal

import (
	"encoding/json"
	"fmt"
	"strings"
)

type OutputFormatter interface {
	Format(results []SecurityGroupResult) error
}

type TextFormatter struct{}

func (f *TextFormatter) Format(results []SecurityGroupResult) error {
	if len(results) == 0 {
		fmt.Println("No security groups with global access found.")
		return nil
	}

	fmt.Printf("Found %d security groups with global access:\n\n", len(results))

	for _, result := range results {
		fmt.Printf("Region: %s\n", result.Region)
		fmt.Printf("Security Group ID: %s\n", result.SecurityGroupID)
		fmt.Printf("Group Name: %s\n", result.GroupName)
		fmt.Printf("Description: %s\n", result.Description)
		fmt.Printf("VPC ID: %s\n", result.VpcID)
		fmt.Println("Risky Inbound Rules:")

		for _, rule := range result.RiskyRules {
			fmt.Printf("  - Protocol: %s", rule.Protocol)
			if rule.FromPort != nil && rule.ToPort != nil {
				if *rule.FromPort == *rule.ToPort {
					fmt.Printf(", Port: %d", *rule.FromPort)
				} else {
					fmt.Printf(", Port Range: %d-%d", *rule.FromPort, *rule.ToPort)
				}
			}
			fmt.Printf(", CIDR: %s\n", strings.Join(rule.CidrBlocks, ", "))
		}
		fmt.Println()
	}

	return nil
}

type JSONFormatter struct{}

func (f *JSONFormatter) Format(results []SecurityGroupResult) error {
	output := struct {
		Count   int                   `json:"count"`
		Results []SecurityGroupResult `json:"results"`
	}{
		Count:   len(results),
		Results: results,
	}

	data, err := json.Marshal(output)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

type MarkdownFormatter struct{}

func (f *MarkdownFormatter) Format(results []SecurityGroupResult) error {
	if len(results) == 0 {
		fmt.Println("## Security Group Scan Results\n\nNo security groups with global access found.")
		return nil
	}

	fmt.Printf("# Security Group Scan Results\n\nFound **%d** security groups with global access:\n\n", len(results))

	// Header
	fmt.Println("| Region | Security Group ID | Group Name | Description | VPC ID | Protocol | Port(s) | CIDR |")
	fmt.Println("|--------|------------------|------------|-------------|--------|----------|---------|------|")

	for _, result := range results {
		for i, rule := range result.RiskyRules {
			var ports string
			if rule.FromPort != nil && rule.ToPort != nil {
				if *rule.FromPort == *rule.ToPort {
					ports = fmt.Sprintf("%d", *rule.FromPort)
				} else {
					ports = fmt.Sprintf("%d-%d", *rule.FromPort, *rule.ToPort)
				}
			} else {
				ports = "All"
			}

			// Only show group info on first row for each security group
			if i == 0 {
				fmt.Printf("| %s | %s | %s | %s | %s | %s | %s | %s |\n",
					result.Region,
					result.SecurityGroupID,
					result.GroupName,
					result.Description,
					result.VpcID,
					rule.Protocol,
					ports,
					strings.Join(rule.CidrBlocks, ", "))
			} else {
				fmt.Printf("| | | | | | %s | %s | %s |\n",
					rule.Protocol,
					ports,
					strings.Join(rule.CidrBlocks, ", "))
			}
		}
	}

	return nil
}

func NewFormatter(format string) (OutputFormatter, error) {
	switch format {
	case "text":
		return &TextFormatter{}, nil
	case "json":
		return &JSONFormatter{}, nil
	case "md", "markdown":
		return &MarkdownFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported output format: %s", format)
	}
}
