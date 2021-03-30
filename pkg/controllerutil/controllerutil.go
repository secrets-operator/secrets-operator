package controllerutil

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

func SetControllerReference(owner metav1.Object, controlled metav1.Object, scheme *runtime.Scheme) error {
	// Validate the owner.
	ro, ok := owner.(runtime.Object)
	if !ok {
		return fmt.Errorf("%T is not a runtime.Object, cannot call SetControllerReference", owner)
	}
	if err := validateOwner(owner, controlled); err != nil {
		return err
	}

	// Create a new controller ref.
	gvk, err := apiutil.GVKForObject(ro, scheme)
	if err != nil {
		return err
	}
	ref := metav1.OwnerReference{
		APIVersion:         gvk.GroupVersion().String(),
		Kind:               gvk.Kind,
		Name:               owner.GetName(),
		UID:                owner.GetUID(),
		BlockOwnerDeletion: pointer.BoolPtr(true),
		Controller:         pointer.BoolPtr(true),
	}

	// Return early with an error if the object is already controlled.
	if existing := metav1.GetControllerOf(controlled); existing != nil && !referSameObject(*existing, ref) {
		return fmt.Errorf("object %s is already controlled by %s", controlled.GetName(), owner.GetName())
	}

	// Update owner references and return.
	upsertOwnerRef(ref, controlled)
	return nil
}

func validateOwner(owner, object metav1.Object) error {
	ownerNs := owner.GetNamespace()
	if ownerNs != "" {
		objNs := object.GetNamespace()
		if objNs == "" {
			return fmt.Errorf("cluster-scoped resource must not have a namespace-scoped owner, owner's namespace %s", ownerNs)
		}
		if ownerNs != objNs {
			return fmt.Errorf("cross-namespace owner references are disallowed, owner's namespace %s, obj's namespace %s", owner.GetNamespace(), object.GetNamespace())
		}
	}
	return nil
}

func upsertOwnerRef(ref metav1.OwnerReference, object metav1.Object) {
	owners := object.GetOwnerReferences()
	idx := indexOwnerRef(owners, ref)
	if idx == -1 {
		owners = append(owners, ref)
	} else {
		owners[idx] = ref
	}
	object.SetOwnerReferences(owners)
}

// indexOwnerRef returns the index of the owner reference in the slice if found, or -1.
func indexOwnerRef(ownerReferences []metav1.OwnerReference, ref metav1.OwnerReference) int {
	for index, r := range ownerReferences {
		if referSameObject(r, ref) {
			return index
		}
	}
	return -1
}

// Returns true if a and b point to the same object
func referSameObject(a, b metav1.OwnerReference) bool {
	aGV, err := schema.ParseGroupVersion(a.APIVersion)
	if err != nil {
		return false
	}

	bGV, err := schema.ParseGroupVersion(b.APIVersion)
	if err != nil {
		return false
	}

	return aGV.Group == bGV.Group && a.Kind == b.Kind && a.Name == b.Name
}
