package source

import (
	"fmt"

	"github.com/secrets-operator/secrets-operator/api/v1alpha1"
	"github.com/secrets-operator/secrets-operator/pkg/generation"
)

func HandleProperty(propertySource v1alpha1.PropertySource) (string, error) {
	if propertySource.PropertyGenerator != nil {
		return generation.Generate(*propertySource.PropertyGenerator)
	}
	return "", fmt.Errorf("unable to determine how to source propery")
}
