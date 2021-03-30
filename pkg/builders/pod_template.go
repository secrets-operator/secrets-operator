package builders

import corev1 "k8s.io/api/core/v1"

type PodTemplateBuilder struct {
	PodTemplate corev1.PodTemplateSpec
}

// NewPodTemplateBuilder returns an initialized PodTemplateBuilder with some defaults.
func NewPodTemplateBuilder() *PodTemplateBuilder {
	builder := &PodTemplateBuilder{}
	return builder.setDefaults()
}

// setDefaults sets up a default Container in the pod template,
// and disables service account token auto mount.
func (b *PodTemplateBuilder) setDefaults() *PodTemplateBuilder {
	b.PodTemplate.Spec.Containers = append(b.PodTemplate.Spec.Containers, corev1.Container{Name: "store-operator"})
	return b
}

func (b *PodTemplateBuilder) WithServiceAccount(serviceAccount string) *PodTemplateBuilder {
	if b.PodTemplate.Spec.ServiceAccountName == "" {
		b.PodTemplate.Spec.ServiceAccountName = serviceAccount
	}
	return b
}

func (b *PodTemplateBuilder) WithImage(image string) *PodTemplateBuilder {
	if b.PodTemplate.Spec.Containers[0].Image == "" {
		b.PodTemplate.Spec.Containers[0].Image = image
	}
	return b
}

// WithEnv appends the given env vars to the Container, unless already provided in the template.
func (b *PodTemplateBuilder) WithEnv(vars ...corev1.EnvVar) *PodTemplateBuilder {
	b.PodTemplate.Spec.Containers[0].Env = vars
	return b
}

// WithLabels sets the given labels, but does not override those that already exist.
func (b *PodTemplateBuilder) WithLabels(labels map[string]string) *PodTemplateBuilder {
	b.PodTemplate.Labels = labels
	return b
}

// WithAnnotations sets the given annotations, but does not override those that already exist.
func (b *PodTemplateBuilder) WithAnnotations(annotations map[string]string) *PodTemplateBuilder {
	b.PodTemplate.Annotations = annotations
	return b
}
