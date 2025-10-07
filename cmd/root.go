package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	vaultPath string
	verbose   bool

	// Version information (set via ldflags during build)
	version = "dev"
	commit  = "none"
	date    = "unknown"

	rootCmd = &cobra.Command{
		Use:   "pass-cli",
		Short: "A secure CLI password and API key manager",
		Long: `Pass-CLI is a secure, cross-platform command-line password and API key manager
designed for developers. It provides local encrypted storage with optional system
keychain integration, allowing developers to securely manage credentials without
relying on cloud services.

Features:
  • AES-256-GCM encryption with PBKDF2 key derivation
  • Native OS keychain integration (Windows Credential Manager, macOS Keychain, Linux Secret Service)
  • Script-friendly output for CI/CD integration
  • Automatic usage tracking
  • Offline-first design with no cloud dependencies

Examples:
  # Initialize a new vault
  pass-cli init

  # Add a new credential
  pass-cli add github

  # Retrieve a credential
  pass-cli get github

  # List all credentials
  pass-cli list

For more information, visit: https://github.com/username/pass-cli`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pass-cli/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&vaultPath, "vault", "", "vault file path (default is $HOME/.pass-cli/vault.enc)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Bind flags to viper
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("vault", rootCmd.PersistentFlags().Lookup("vault"))
}

// GetVaultPath returns the vault path from flag, config, or default
func GetVaultPath() string {
	// Priority: flag > config > default
	if vaultPath != "" {
		return vaultPath
	}

	if viper.IsSet("vault") {
		return viper.GetString("vault")
	}

	// Default vault path
	home, err := os.UserHomeDir()
	if err != nil {
		return ".pass-cli/vault.enc"
	}

	return filepath.Join(home, ".pass-cli", "vault.enc")
}

// IsVerbose returns whether verbose mode is enabled
func IsVerbose() bool {
	return verbose || viper.GetBool("verbose")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".pass-cli" (without extension).
		viper.AddConfigPath(home + "/.pass-cli")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}
