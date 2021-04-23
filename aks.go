package pack

import (
	azureaks "github.com/citihub/probr-pack-aks/internal/azure/aks"
	kubeear "github.com/citihub/probr-pack-aks/internal/azure/kubernetes/encryption-at-rest"
	kubeaks "github.com/citihub/probr-pack-aks/internal/azure/kubernetes/iam"
	"github.com/citihub/probr-sdk/probeengine"
	"github.com/markbates/pkger"
)

// GetProbes returns a list of probe objects
func GetProbes() []probeengine.Probe {
	// TODO: Config make this configurable once the config handler has been refactored

	/*if config.Vars.ServicePacks.AKS.IsExcluded() {
		return nil
	}*/
	return []probeengine.Probe{
		azureaks.Probe,
		kubeaks.Probe,
		kubeear.Probe,
	}
}

func init() {
	// This line will ensure that all static files are bundled into pked.go file when using pkger cli tool
	// See: https://github.com/markbates/pkger
	pkger.Include("/internal/azure/aks/aks.feature")
	pkger.Include("/internal/azure/aks/aks.rego")
	pkger.Include("/internal/azure/kubernetes/iam/iam.feature")
	pkger.Include("/internal/azure/kubernetes/encryption-at-rest/encryption-at-rest.feature")
}