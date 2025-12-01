package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// version info (set via ldflags at build time)
	version = "dev"
	commit  = "none"
	date    = "unknown"

	cfgFile      string
	region       string
	excludeFile  string
	outputFormat string
)

var rootCmd = &cobra.Command{
	Use:     "sglobal",
	Short:   "AWS Security Group Global Access Scanner",
	Version: version,
	Long: `sglobal scans AWS Security Groups for rules that allow global access (0.0.0.0/0).
It helps identify potentially risky security configurations across your AWS infrastructure.`,
	Run: runScan,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Set custom version output
	rootCmd.SetVersionTemplate(fmt.Sprintf("sglobal version %s (commit: %s, built: %s)\n", version, commit, date))

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sglobal.yaml)")
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "", "AWS region to scan (default: current profile region, 'all' for all regions)")
	rootCmd.PersistentFlags().StringVarP(&excludeFile, "exclude-file", "e", "", "file containing security group IDs to exclude (one per line)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json", "output format: text, json, md/markdown")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".sglobal")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
