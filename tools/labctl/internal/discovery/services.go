package discovery

import "github.com/illmadecoder/labctl/internal/k8s"

// BuildServiceVars creates a template variable map from experiment service info.
func BuildServiceVars(exp *k8s.ExperimentInfo) map[string]string {
	vars := make(map[string]string)
	if exp == nil {
		return vars
	}
	for _, svc := range exp.Services {
		if svc.Endpoint != "" {
			vars[svc.Name] = svc.Endpoint
		}
	}
	return vars
}
