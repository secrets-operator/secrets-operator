package deployment

import (
	"context"
	"fmt"
	"github.com/secrets-operator/secrets-operator/api/v1alpha1"
	secretoperatorv1alpha1 "github.com/secrets-operator/secrets-operator/api/v1alpha1"
	"github.com/secrets-operator/secrets-operator/pkg/builders"
	"github.com/secrets-operator/secrets-operator/pkg/controllerutil"
	"github.com/secrets-operator/secrets-operator/pkg/secretstores/gcp"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	StoreOperatorImage = "mellardc/storeoperator:1.0"
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

// New creates a Deployment from the given params.
func New(params Params) appsv1.Deployment {
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      params.Name,
			Namespace: params.Namespace,
			Labels:    params.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: params.Selector,
			},
			Template: params.PodTemplateSpec,
			Replicas: &params.Replicas,
			Strategy: params.Strategy,
		},
	}
}

func DeploymentParams(store v1alpha1.SecretStore) Params {
	podSpec := newPodTemplateSpec(store)
	return Params{
		Name:            store.Name,
		Namespace:       store.Namespace,
		Selector:        NewLabels(store.Name),
		PodTemplateSpec: podSpec,
		Replicas:        1,
	}
}

func newPodTemplateSpec(store v1alpha1.SecretStore) corev1.PodTemplateSpec {
	builder := builders.NewPodTemplateBuilder().
		WithImage(StoreOperatorImage).
		WithLabels(NewLabels(store.Name))

	if store.Spec.Provider.GcpSecretsManager != nil {
		return gcp.GcpPodTemplateSpec(store, builder)
	}

	return builder.PodTemplate
}

func Reconcile(client kubernetes.Interface, deployment appsv1.Deployment, owner client.Object) error {

	err := secretoperatorv1alpha1.AddToScheme(scheme.Scheme)
	if err != nil {
		return err
	}

	if err := controllerutil.SetControllerReference(owner, &deployment, scheme.Scheme); err != nil {
		return err
	}

	create := func() error {
		_, err := client.AppsV1().Deployments(deployment.Namespace).Create(context.Background(), &deployment, metav1.CreateOptions{})
		return err
	}
	// First check if deployment exists
	_, err = client.AppsV1().Deployments(deployment.Namespace).Get(context.Background(), deployment.Name, metav1.GetOptions{})
	if err != nil && apierrors.IsNotFound(err) {
		return create()
	} else if err != nil {
		return fmt.Errorf("failed to get deployment %s/%s: %w", deployment.Namespace, deployment.Name, err)
	}

	_, err = client.AppsV1().Deployments(deployment.Namespace).Update(context.Background(), &deployment, metav1.UpdateOptions{})
	return err
}

// NewLabels constructs a new set of labels for a Kibana pod
func NewLabels(storeName string) map[string]string {
	return map[string]string{"secret-operator.io/store-operator": storeName}
}
