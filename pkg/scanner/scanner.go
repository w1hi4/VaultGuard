package scanner

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

type Rule struct {
	ID          string   `yaml:"id"`
	Description string   `yaml:"description"`
	Regex       string   `yaml:"regex"`
	Severity    string   `yaml:"severity"`
	Tags        []string `yaml:"tags"`
}

type Config struct {
	Rules   []Rule `yaml:"rules"`
	Entropy struct {
		Enabled   bool    `yaml:"enabled"`
		Threshold float64 `yaml:"threshold"`
		MinLength int     `yaml:"min_length"`
	} `yaml:"entropy"`
	ExcludePaths []string `yaml:"exclude_paths"`
}

type ScanLine struct {
	Commit  string
	File    string
	LineNum int
	Content string
}

type Finding struct {
	RuleID      string `json:"rule_id"`
	Description string `json:"description"`
	Commit      string `json:"commit"`
	File        string `json:"file"`
	LineNumber  int    `json:"line_number"`
	Match       string `json:"match"`
	Severity    string `json:"severity"`
}

type Scanner struct {
	Config   Config
	compiled []compiledRule
}

type compiledRule struct {
	Rule  Rule
	Regex *regexp.Regexp
}

func NewScanner(configPath string) (*Scanner, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// Pre-compile all regexes
	var compiled []compiledRule
	for _, rule := range config.Rules {
		re, err := regexp.Compile(rule.Regex)
		if err != nil {
			continue
		}
		compiled = append(compiled, compiledRule{Rule: rule, Regex: re})
	}

	return &Scanner{Config: config, compiled: compiled}, nil
}

func (s *Scanner) ScanRepo(repoPath string) ([]Finding, error) {
	return s.scanWithStrategy(repoPath, []string{"git", "log", "-p", "--all", "--full-history", "--diff-filter=A"})
}

// ScanRepoFull scans entire history including modifications (not just additions)
func (s *Scanner) ScanRepoFull(repoPath string) ([]Finding, error) {
	return s.scanWithStrategy(repoPath, []string{"git", "log", "-p", "--all", "--full-history"})
}

func (s *Scanner) scanWithStrategy(repoPath string, gitArgs []string) ([]Finding, error) {
	cmd := exec.Command(gitArgs[0], gitArgs[1:]...)
	cmd.Dir = repoPath
	output, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var findings []Finding
	findingsMu := sync.Mutex{}
	seen := make(map[string]bool)
	seenMu := sync.Mutex{}

	// Concurrency setup
	jobs := make(chan ScanLine, 1000)
	wg := sync.WaitGroup{}
	numWorkers := 8 // Scalable based on business logic or CPU

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				s.analyzeLine(job, &findings, &findingsMu, seen, &seenMu)
			}
		}()
	}

	lineScanner := bufio.NewScanner(output)
	buf := make([]byte, 0, 1024*1024)
	lineScanner.Buffer(buf, 10*1024*1024)

	var currentCommit string
	var currentFile string
	var lineNum int

	for lineScanner.Scan() {
		line := lineScanner.Text()

		if strings.HasPrefix(line, "commit ") && len(line) > 7 {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				currentCommit = parts[1]
			}
			continue
		}

		if strings.HasPrefix(line, "+++ b/") {
			currentFile = strings.TrimPrefix(line, "+++ b/")
			lineNum = 0
			continue
		}

		if strings.HasPrefix(line, "--- ") || s.isExcluded(currentFile) {
			continue
		}

		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			lineNum++
			jobs <- ScanLine{
				Commit:  currentCommit,
				File:    currentFile,
				LineNum: lineNum,
				Content: line[1:],
			}
		}
	}

	close(jobs)
	wg.Wait()
	cmd.Wait()
	return findings, nil
}

func (s *Scanner) analyzeLine(job ScanLine, findings *[]Finding, fMu *sync.Mutex, seen map[string]bool, sMu *sync.Mutex) {
	// Check all compiled rules
	for _, cr := range s.compiled {
		match := cr.Regex.FindString(job.Content)
		if match != "" {
			key := cr.Rule.ID + "|" + match + "|" + job.File
			sMu.Lock()
			if !seen[key] {
				seen[key] = true
				sMu.Unlock()

				fMu.Lock()
				*findings = append(*findings, Finding{
					RuleID:      cr.Rule.ID,
					Description: cr.Rule.Description,
					Commit:      job.Commit,
					File:        job.File,
					LineNumber:  job.LineNum,
					Match:       match,
					Severity:    cr.Rule.Severity,
				})
				fMu.Unlock()
			} else {
				sMu.Unlock()
			}
		}
	}

	// Entropy check
	if s.Config.Entropy.Enabled {
		tokens := strings.FieldsFunc(job.Content, func(r rune) bool {
			return r == '"' || r == '\'' || r == '=' || r == ':' || r == ' ' || r == ',' || r == '(' || r == ')' || r == '{' || r == '}'
		})
		for _, token := range tokens {
			token = strings.TrimSpace(token)
			if len(token) >= s.Config.Entropy.MinLength {
				ent := calculateEntropy(token)
				if ent >= s.Config.Entropy.Threshold {
					key := "entropy|" + token + "|" + job.File
					sMu.Lock()
					if !seen[key] {
						seen[key] = true
						sMu.Unlock()

						fMu.Lock()
						*findings = append(*findings, Finding{
							RuleID:      "high-entropy",
							Description: "High entropy string detected",
							Commit:      job.Commit,
							File:        job.File,
							LineNumber:  job.LineNum,
							Match:       token,
							Severity:    "MEDIUM",
						})
						fMu.Unlock()
					} else {
						sMu.Unlock()
					}
				}
			}
		}
	}
}

func (s *Scanner) isExcluded(path string) bool {
	for _, p := range s.Config.ExcludePaths {
		if strings.Contains(path, p) {
			return true
		}
	}
	return false
}

// ScanRemote clones a repository to a temporary folder, scans it, and cleans up.
func (s *Scanner) ScanRemote(repoURL string, deep bool) ([]Finding, int, error) {
	tempDir, err := os.MkdirTemp("", "vaultguard-*")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Clone the repository
	cmd := exec.Command("git", "clone", "--depth", "1", repoURL, tempDir)
	if deep {
		// For deep scan, we need full history, so don't use --depth 1
		cmd = exec.Command("git", "clone", repoURL, tempDir)
	}
	
	if err := cmd.Run(); err != nil {
		return nil, 0, fmt.Errorf("failed to clone repository: %v", err)
	}

	var findings []Finding
	if deep {
		findings, err = s.ScanRepoFull(tempDir)
	} else {
		findings, err = s.ScanRepo(tempDir)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("scan failed: %v", err)
	}

	return findings, len(findings), nil
}

// GenerateReport creates a detailed Markdown report of the findings.
func (s *Scanner) GenerateReport(repoName string, findings []Finding) string {
	var sb strings.Builder
	sb.WriteString("# 🛡️ VaultGuard Security Audit Report\n\n")
	sb.WriteString(fmt.Sprintf("**Repository:** `%s`  \n", repoName))
	sb.WriteString(fmt.Sprintf("**Date:** %s  \n", time.Now().Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("**Total Findings:** %d  \n\n", len(findings)))

	if len(findings) == 0 {
		sb.WriteString("## ✅ Verdict: Secure\n")
		sb.WriteString("No sensitive information or hardcoded secrets were detected in the analyzed history.\n")
		return sb.String()
	}

	// Severity breakdown
	stats := make(map[string]int)
	for _, f := range findings {
		stats[strings.ToUpper(f.Severity)]++
	}

	sb.WriteString("## 📊 Finding Breakdown\n\n")
	sb.WriteString(fmt.Sprintf("- 🔴 **CRITICAL:** %d\n", stats["CRITICAL"]))
	sb.WriteString(fmt.Sprintf("- 🟠 **HIGH:** %d\n", stats["HIGH"]))
	sb.WriteString(fmt.Sprintf("- 🟡 **MEDIUM:** %d\n", stats["MEDIUM"]))
	sb.WriteString(fmt.Sprintf("- 🔵 **LOW:** %d\n\n", stats["LOW"]))

	sb.WriteString("--- \n\n")
	sb.WriteString("## 🔍 Detailed Findings\n\n")

	for i, f := range findings {
		sb.WriteString(fmt.Sprintf("### %d. [%s] %s\n", i+1, strings.ToUpper(f.Severity), f.Description))
		sb.WriteString(fmt.Sprintf("- **Rule ID:** `%s`\n", f.RuleID))
		sb.WriteString(fmt.Sprintf("- **File:** `%s:%d`\n", f.File, f.LineNumber))
		sb.WriteString(fmt.Sprintf("- **Match:** `%s`\n", f.Match))
		sb.WriteString(fmt.Sprintf("- **Commit:** `%s`\n\n", f.Commit))
		sb.WriteString("---\n\n")
	}

	sb.WriteString("\n*This report was generated automatically by VaultGuard - The local-first secret scanner.*\n")
	return sb.String()
}

func calculateEntropy(data string) float64 {
	if data == "" {
		return 0
	}
	frequencies := make(map[rune]float64)
	for _, char := range data {
		frequencies[char]++
	}
	var entropy float64
	length := float64(len(data))
	for _, count := range frequencies {
		p := count / length
		entropy -= p * math.Log2(p)
	}
	return entropy
}
