package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"vaultguard/pkg/api"
	"vaultguard/pkg/scanner"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func main() {
	var configPath string
	var repoPath string
	var port string
	var deepScan bool
	var jsonOutput bool
	var excludePatterns []string

	var rootCmd = &cobra.Command{
		Use:   "vaultguard",
		Short: "VaultGuard is an advanced secret scanner for Git repositories",
		Long:  `VaultGuard scans your Git history for API keys, passwords, and other sensitive information.`,
	}

	var scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "Perform a security scan on a repository",
		Run: func(cmd *cobra.Command, args []string) {
			if repoPath == "" {
				repoPath, _ = os.Getwd()
			}
			if configPath == "" {
				configPath = "pkg/scanner/rules.yaml"
			}

			s, err := scanner.NewScanner(configPath)
			if err != nil {
				fmt.Printf("Error initializing scanner: %v\n", err)
				os.Exit(1)
			}
			
			// Append manual exclusions
			if len(excludePatterns) > 0 {
				s.Config.ExcludePaths = append(s.Config.ExcludePaths, excludePatterns...)
			}

			if !jsonOutput {
				color.Cyan("Scanning repository: %s", repoPath)
				if deepScan {
					color.Yellow("Mode: DEEP SCAN (full history)")
				}
			}

			var findings []scanner.Finding
			
			// Detect if repoPath is a URL
			isURL := strings.HasPrefix(repoPath, "http://") || 
					 strings.HasPrefix(repoPath, "https://") || 
					 strings.HasPrefix(repoPath, "git@")

			if isURL {
				if !jsonOutput {
					color.Yellow("Remote repository detected. Cloning...")
				}
				findings, _, err = s.ScanRemote(repoPath, deepScan)
			} else {
				if deepScan {
					findings, err = s.ScanRepoFull(repoPath)
				} else {
					findings, err = s.ScanRepo(repoPath)
				}
			}
			if err != nil {
				fmt.Printf("Error scanning repo: %v\n", err)
				os.Exit(1)
			}

			if jsonOutput {
				data, _ := json.MarshalIndent(findings, "", "  ")
				fmt.Println(string(data))
				return
			}

			if len(findings) == 0 {
				color.Green("No sensitive information found. Your repository is safe!")
				return
			}

			color.Red("Found %d potential leaks!\n", len(findings))
			for _, f := range findings {
				sevColor := color.New(color.FgYellow)
				switch f.Severity {
				case "CRITICAL":
					sevColor = color.New(color.FgRed, color.Bold)
				case "HIGH":
					sevColor = color.New(color.FgHiRed)
				case "MEDIUM":
					sevColor = color.New(color.FgYellow)
				case "LOW":
					sevColor = color.New(color.FgCyan)
				}
				sevColor.Printf("[%s] [%s] %s\n", f.Severity, f.RuleID, f.Description)
				fmt.Printf("  File: %s:%d\n", f.File, f.LineNumber)
				fmt.Printf("  Commit: %s\n", f.Commit)
				fmt.Printf("  Match: %s\n", f.Match)
				fmt.Println("--------------------------------------------------")
			}
		},
	}

	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Start the VaultGuard web dashboard",
		Run: func(cmd *cobra.Command, args []string) {
			if configPath == "" {
				configPath = "pkg/scanner/rules.yaml"
			}
			color.Cyan("Starting VaultGuard dashboard on port %s...", port)
			color.Green("Dashboard: http://localhost:3000")
			color.Yellow("API Server: http://localhost:%s", port)

			if err := api.StartServer(port, configPath); err != nil {
				fmt.Printf("Error starting server: %v\n", err)
				os.Exit(1)
			}
		},
	}

	scanCmd.Flags().StringVarP(&configPath, "config", "c", "", "path to rules.yaml")
	scanCmd.Flags().StringVarP(&repoPath, "path", "p", "", "path to git repository")
	scanCmd.Flags().BoolVarP(&deepScan, "deep", "d", false, "deep scan: scan full git history")
	scanCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "output results as JSON")
	scanCmd.Flags().StringSliceVarP(&excludePatterns, "exclude", "e", []string{}, "paths or patterns to exclude from scan")

	serveCmd.Flags().StringVarP(&port, "port", "P", "8080", "port for the API server")
	serveCmd.Flags().StringVarP(&configPath, "config", "c", "", "path to rules.yaml")

	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(serveCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
