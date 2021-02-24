package azurekeyvault

import "github.com/secrets-operator/secrets-operator/api/v1alpha1"

type AzureKeyVaultProvider struct {
	v1alpha1.AzureKeyVaultProvider
}

func (a AzureKeyVaultProvider) Location() string {
	return a.VaultName
}
