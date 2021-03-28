// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package deployment

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/secrets-operator/secrets-operator/pkg/reconciler"
	"hash/fnv"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// TemplateHashLabelName is a label to annotate a Kubernetes resource
	// with the hash of its initial template before creation.
	TemplateHashLabelName = "secrets-operator/template-hash"
)

var (
	defaultRevisionHistoryLimit int32
)

// Params to specify a Deployment specification.
type Params struct {
	Name            string
	Namespace       string
	Selector        map[string]string
	Labels          map[string]string
	PodTemplateSpec corev1.PodTemplateSpec
	Replicas        int32
	Strategy        appsv1.DeploymentStrategy
}

// Int32 returns a pointer to an Int32
func Int32(v int32) *int32 { return &v }

// GetTemplateHashLabel returns the template hash label value if set, or an empty string.
func GetTemplateHashLabel(labels map[string]string) string {
	return labels[TemplateHashLabelName]
}

// SetTemplateHashLabel adds a label containing the hash of the given template into the
// given labels. This label can then be used for template comparisons.
func SetTemplateHashLabel(labels map[string]string, template interface{}) map[string]string {
	return setHashLabel(TemplateHashLabelName, labels, template)
}

func setHashLabel(labelName string, labels map[string]string, template interface{}) map[string]string {
	if labels == nil {
		labels = map[string]string{}
	}
	labels[labelName] = HashObject(template)
	return labels
}

// HashObject writes the specified object to a hash using the spew library
// which follows pointers and prints actual values of the nested objects
// ensuring the hash does not change when a pointer changes.
// The returned hash can be used for object comparisons.
//
// This is inspired by controller revisions in StatefulSets:
// https://github.com/kubernetes/kubernetes/blob/8de1569ddae62e8fab559fe6bd210a5d6100a277/pkg/controller/history/controller_history.go#L89-L101
func HashObject(object interface{}) string {
	hf := fnv.New32()
	printer := spew.ConfigState{
		Indent:         " ",
		SortKeys:       true,
		DisableMethods: true,
		SpewKeys:       true,
	}
	_, _ = printer.Fprintf(hf, "%#v", object)
	return fmt.Sprint(hf.Sum32())
}

// New creates a Deployment from the given params.
func New(params Params) appsv1.Deployment {
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      params.Name,
			Namespace: params.Namespace,
			Labels:    params.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			RevisionHistoryLimit: Int32(defaultRevisionHistoryLimit),
			Selector: &metav1.LabelSelector{
				MatchLabels: params.Selector,
			},
			Template: params.PodTemplateSpec,
			Replicas: &params.Replicas,
			Strategy: params.Strategy,
		},
	}
}

// ReconcileDeployment creates or updates the given deployment for the specified owner.
func Reconcile(
	k8sClient client.Client,
	expected appsv1.Deployment,
	owner client.Object,
) (appsv1.Deployment, error) {
	// label the deployment with a hash of itself
	expected = WithTemplateHash(expected)

	reconciled := &appsv1.Deployment{}
	err := reconciler.ReconcileResource(reconciler.Params{
		Client:     k8sClient,
		Owner:      owner,
		Expected:   &expected,
		Reconciled: reconciled,
		NeedsUpdate: func() bool {
			// compare hash of the deployment at the time it was built
			return GetTemplateHashLabel(reconciled.Labels) != GetTemplateHashLabel(expected.Labels)
		},
		UpdateReconciled: func() {
			expected.DeepCopyInto(reconciled)
		},
	})
	return *reconciled, err
}

// WithTemplateHash returns a new deployment with a hash of its template to ease comparisons.
func WithTemplateHash(d appsv1.Deployment) appsv1.Deployment {
	dCopy := *d.DeepCopy()
	dCopy.Labels = SetTemplateHashLabel(dCopy.Labels, dCopy)
	return dCopy
}
