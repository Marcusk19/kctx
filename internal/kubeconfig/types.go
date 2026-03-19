package kubeconfig

// Config represents a minimal kubeconfig structure.
type Config struct {
	APIVersion     string      `yaml:"apiVersion"`
	Kind           string      `yaml:"kind"`
	CurrentContext string      `yaml:"current-context,omitempty"`
	Preferences    struct{}    `yaml:"preferences,omitempty"`
	Clusters       []Cluster   `yaml:"clusters,omitempty"`
	Contexts       []Context   `yaml:"contexts,omitempty"`
	Users          []User      `yaml:"users,omitempty"`
}

type Cluster struct {
	Name    string      `yaml:"name"`
	Cluster ClusterData `yaml:"cluster"`
}

type ClusterData struct {
	Server                string `yaml:"server,omitempty"`
	CertificateAuthority  string `yaml:"certificate-authority,omitempty"`
	CertificateAuthorityData string `yaml:"certificate-authority-data,omitempty"`
	InsecureSkipTLSVerify bool   `yaml:"insecure-skip-tls-verify,omitempty"`
}

type Context struct {
	Name    string      `yaml:"name"`
	Context ContextData `yaml:"context"`
}

type ContextData struct {
	Cluster   string `yaml:"cluster,omitempty"`
	Namespace string `yaml:"namespace,omitempty"`
	User      string `yaml:"user,omitempty"`
}

type User struct {
	Name string      `yaml:"name"`
	User interface{} `yaml:"user,omitempty"`
}
