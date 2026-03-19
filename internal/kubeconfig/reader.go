package kubeconfig

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ReadConfig reads and parses a kubeconfig file.
func ReadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SourceKubeconfig returns the path to the source kubeconfig.
// It checks KCTX_SOURCE_KUBECONFIG, KCTX_ORIGINAL_KUBECONFIG, then defaults to ~/.kube/config.
func SourceKubeconfig() string {
	if v := os.Getenv("KCTX_SOURCE_KUBECONFIG"); v != "" {
		return v
	}
	if v := os.Getenv("KCTX_ORIGINAL_KUBECONFIG"); v != "" {
		// May be colon-separated; return as-is for KUBECONFIG merging
		return v
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".kube", "config")
}

// ListContexts returns all context names from the source kubeconfig(s).
func ListContexts(sourcePath string) ([]string, string, error) {
	paths := splitKubeconfig(sourcePath)
	var allContexts []string
	seen := map[string]bool{}
	currentContext := ""

	for _, p := range paths {
		cfg, err := ReadConfig(p)
		if err != nil {
			continue
		}
		if currentContext == "" && cfg.CurrentContext != "" {
			currentContext = cfg.CurrentContext
		}
		for _, ctx := range cfg.Contexts {
			if !seen[ctx.Name] {
				seen[ctx.Name] = true
				allContexts = append(allContexts, ctx.Name)
			}
		}
	}
	return allContexts, currentContext, nil
}

// GetContextInfo returns the ContextData for a named context from the source kubeconfig(s).
func GetContextInfo(sourcePath, contextName string) *ContextData {
	paths := splitKubeconfig(sourcePath)
	for _, p := range paths {
		cfg, err := ReadConfig(p)
		if err != nil {
			continue
		}
		for _, ctx := range cfg.Contexts {
			if ctx.Name == contextName {
				return &ctx.Context
			}
		}
	}
	return nil
}

func splitKubeconfig(path string) []string {
	parts := strings.Split(path, ":")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
