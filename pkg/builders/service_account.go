package builders

import corev1 "k8s.io/api/core/v1"

type ServiceAccountBuilder struct {
	ServiceAccount corev1.ServiceAccount
}

// NewServiceAccountBuilder returns an initialized ServiceAccountBuilder with some defaults.
func NewServiceAccountBuilder(account corev1.ServiceAccount) *ServiceAccountBuilder {
	builder := &ServiceAccountBuilder{account}
	return builder
}

func (b *ServiceAccountBuilder) WithAnnotations(annotations map[string]string) *ServiceAccountBuilder {
	if b.ServiceAccount.Annotations == nil {
		b.ServiceAccount.Annotations = annotations
	}
	return b
}
