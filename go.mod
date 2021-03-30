module github.com/secrets-operator/secrets-operator

go 1.16

require (
	cloud.google.com/go v0.75.0 // indirect
	github.com/go-logr/logr v0.3.0
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/pkg/errors v0.9.1
	github.com/sethvargo/go-password v0.2.0
	golang.org/x/oauth2 v0.0.0-20210113205817-d3ed898aa8a3 // indirect
	k8s.io/api v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v0.20.4
	k8s.io/utils v0.0.0-20210111153108-fddb29f9d009
	sigs.k8s.io/controller-runtime v0.8.1
)
