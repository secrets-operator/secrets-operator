package gcp

import (
	"fmt"
	"github.com/secrets-operator/secrets-operator/api/v1alpha1"
	"github.com/secrets-operator/secrets-operator/pkg/builders"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GcpPodTemplateSpec(store v1alpha1.SecretStore, builder *builders.PodTemplateBuilder) corev1.PodTemplateSpec {
	return builder.
		WithServiceAccount(store.Spec.Provider.GcpSecretsManager.Auth.WorkloadIdentity.ServiceAccount).
		PodTemplate
}

func GcpServiceAccount(store v1alpha1.SecretStore) corev1.ServiceAccount {
	b := builders.NewServiceAccountBuilder(corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      store.Spec.Provider.GcpSecretsManager.Auth.WorkloadIdentity.ServiceAccount,
			Namespace: store.Namespace,
		},
	})
	return b.WithAnnotations(map[string]string{
		"iam.gke.io/gcp-service-account": fmt.Sprintf("%s@%s.iam.gserviceaccount.com",
			store.Spec.Provider.GcpSecretsManager.Auth.WorkloadIdentity.GcpServiceAccount,
			store.Spec.Provider.GcpSecretsManager.ProjectId),
	}).ServiceAccount
}
