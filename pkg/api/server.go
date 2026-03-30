package api

import (
	"net/http"
	"path/filepath"
	"strings"
	"vaultguard/pkg/scanner"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartServer(port string, configPath string) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Default())

	s, err := scanner.NewScanner(configPath)
	if err != nil {
		return err
	}

	// Serve UI assets
	r.GET("/favicon.png", func(c *gin.Context) {
		c.File(filepath.Join("web", "favicon.png"))
	})

	r.GET("/", func(c *gin.Context) {
		c.File(filepath.Join("web", "index.html"))
	})

	// Scan API
	r.GET("/api/scan", func(c *gin.Context) {
		repoPath := c.Query("path")
		if repoPath == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "path parameter is required"})
			return
		}

		deep := c.Query("deep") == "true"
		excludes := c.Query("exclude")
		
		// Temporary scanner instance to allow custom exclusions without polluting global config
		localScanner, _ := scanner.NewScanner(configPath)
		if excludes != "" {
			localScanner.Config.ExcludePaths = append(localScanner.Config.ExcludePaths, strings.Split(excludes, ",")...)
		}

		var findings []scanner.Finding
		var scanErr error

		// Robust URL detection
		isURL := strings.HasPrefix(repoPath, "http://") || 
				 strings.HasPrefix(repoPath, "https://") || 
				 strings.HasPrefix(repoPath, "git@")

		if isURL {
			findings, _, scanErr = localScanner.ScanRemote(repoPath, deep)
		} else {
			if deep {
				findings, scanErr = localScanner.ScanRepoFull(repoPath)
			} else {
				findings, scanErr = localScanner.ScanRepo(repoPath)
			}
		}

		if scanErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": scanErr.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"findings": findings,
			"count":    len(findings),
		})
	})

	// Report API
	r.POST("/api/report", func(c *gin.Context) {
		var req struct {
			RepoName string            `json:"repo_name"`
			Findings []scanner.Finding `json:"findings"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		report := s.GenerateReport(req.RepoName, req.Findings)
		c.Header("Content-Disposition", "attachment; filename=VaultGuard_Report.md")
		c.Data(http.StatusOK, "text/markdown", []byte(report))
	})

	return r.Run(":" + port)
}
