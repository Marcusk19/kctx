package kubeconfig

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// WriteSessionConfig writes a minimal kubeconfig overlay to the given path.
func WriteSessionConfig(dir string, currentContext string) error {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	cfg := Config{
		APIVersion:     "v1",
		Kind:           "Config",
		CurrentContext: currentContext,
	}
	return writeConfig(filepath.Join(dir, "config"), &cfg)
}

// WriteSessionConfigWithNamespace writes a kubeconfig overlay that includes
// a context entry to shadow the source context's namespace.
func WriteSessionConfigWithNamespace(dir string, currentContext string, ctxData *ContextData, namespace string) error {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	cfg := Config{
		APIVersion:     "v1",
		Kind:           "Config",
		CurrentContext: currentContext,
		Contexts: []Context{
			{
				Name: currentContext,
				Context: ContextData{
					Cluster:   ctxData.Cluster,
					Namespace: namespace,
					User:      ctxData.User,
				},
			},
		},
	}
	return writeConfig(filepath.Join(dir, "config"), &cfg)
}

func writeConfig(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
