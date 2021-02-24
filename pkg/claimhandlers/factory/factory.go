package factory

import (
	"context"
	"fmt"

	"github.com/secrets-operator/secrets-operator/api/v1alpha1"
	"github.com/secrets-operator/secrets-operator/pkg/claimhandlers"
	"github.com/secrets-operator/secrets-operator/pkg/claimhandlers/kubernetesclaim"
)

func CreateClaimHandler(claim v1alpha1.SecretClaim, ctx context.Context) (claimhandlers.ClaimHandler, error) {
	if claim.Spec.KubernetesClaim != nil {
		return kubernetesclaim.NewHandler(claim, ctx), nil
	}
	return nil, fmt.Errorf("unable to create claim handler - unable to determine claim type")
}
