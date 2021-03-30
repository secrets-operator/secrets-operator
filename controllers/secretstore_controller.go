/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/secrets-operator/secrets-operator/pkg/clients/kube"
	"github.com/secrets-operator/secrets-operator/pkg/deployment"
	"github.com/secrets-operator/secrets-operator/pkg/secretstores/gcp"
	"github.com/secrets-operator/secrets-operator/pkg/serviceaccount"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	secretoperatorv1alpha1 "github.com/secrets-operator/secrets-operator/api/v1alpha1"
)

// SecretStoreReconciler reconciles a SecretStore object
type SecretStoreReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=secret-operator.io,resources=secretstores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=secret-operator.io,resources=secretstores/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list;create;delete

func (r *SecretStoreReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("secretstore", req.NamespacedName)

	var store secretoperatorv1alpha1.SecretStore
	if err := r.Get(ctx, req.NamespacedName, &store); err != nil {
		log.Error(err, "unable to fetch CronJob")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// This reconciler needs to provision a Secret Store operator for every SecretStore CRD it encounters.
	// For Google this will mean creating a deployment running under a specific k8s service account.
	// For Azure this will mean creating a deployment with specific pod annotations for use with aad-pod-identity
	// For AWS this will mean using IRSA service account annotations

	kubeClient, err := kube.CreateClientSet()
	if err != nil {
		return ctrl.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
	}

	if store.Spec.Provider.GcpSecretsManager != nil && store.Spec.Provider.GcpSecretsManager.Auth.WorkloadIdentity != nil {
		expectedServiceAccount := gcp.GcpServiceAccount(store)
		err = serviceaccount.Reconcile(kubeClient, expectedServiceAccount, &store)
		if err != nil {
			return ctrl.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
		}
	}

	deploymentParams := deployment.DeploymentParams(store)
	expectedDeployment := deployment.New(deploymentParams)
	err = deployment.Reconcile(kubeClient, expectedDeployment, &store)
	if err != nil {
		return ctrl.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
	}

	return ctrl.Result{}, nil
}

func (r *SecretStoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&secretoperatorv1alpha1.SecretStore{}).
		Complete(r)
}
