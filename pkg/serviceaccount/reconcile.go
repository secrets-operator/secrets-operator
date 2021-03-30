package serviceaccount

import (
	"context"
	"fmt"
	secretoperatorv1alpha1 "github.com/secrets-operator/secrets-operator/api/v1alpha1"
	"github.com/secrets-operator/secrets-operator/pkg/controllerutil"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Reconcile(client kubernetes.Interface, account corev1.ServiceAccount, owner client.Object) error {

	err := secretoperatorv1alpha1.AddToScheme(scheme.Scheme)
	if err != nil {
		return err
	}

	if err = controllerutil.SetControllerReference(owner, &account, scheme.Scheme); err != nil {
		return err
	}

	create := func() error {
		_, err := client.CoreV1().ServiceAccounts(account.Namespace).Create(context.Background(), &account, metav1.CreateOptions{})
		return err
	}
	// First check if deployment exists
	_, err = client.CoreV1().ServiceAccounts(account.Namespace).Get(context.Background(), account.Name, metav1.GetOptions{})
	if err != nil && apierrors.IsNotFound(err) {
		return create()
	} else if err != nil {
		return fmt.Errorf("failed to get service account %s/%s: %w", account.Namespace, account.Name, err)
	}

	_, err = client.CoreV1().ServiceAccounts(account.Namespace).Update(context.Background(), &account, metav1.UpdateOptions{})
	return err
}
