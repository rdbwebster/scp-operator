package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(GroupVersion,
		&SCPcluster{},
		&SCPclusterList{},
		&ManagedOperator{},
		&ManagedOperatorList{},
	)

	metav1.AddToGroupVersion(scheme, GroupVersion)
	return nil
}
