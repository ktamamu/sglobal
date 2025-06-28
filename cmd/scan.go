package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ktamamu/sglobal/internal"
	"github.com/spf13/cobra"
)

func runScan(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	// Load exclude list
	excludeIDs, err := internal.LoadExcludeList(excludeFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading exclude file: %v\n", err)
		os.Exit(1)
	}

	// Initialize AWS client
	awsClient, err := internal.NewAWSClient(ctx, region)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing AWS client: %v\n", err)
		os.Exit(1)
	}

	// Scan security groups
	results, err := awsClient.ScanSecurityGroups(ctx, excludeIDs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning security groups: %v\n", err)
		os.Exit(1)
	}

	// Format and output results
	formatter, err := internal.NewFormatter(outputFormat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating formatter: %v\n", err)
		os.Exit(1)
	}

	if err := formatter.Format(results); err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
		os.Exit(1)
	}
}
