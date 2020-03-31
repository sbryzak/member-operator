package test

import (
	"context"

	"github.com/codeready-toolchain/toolchain-common/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type NamespaceAssertion struct {
	namespace      *corev1.Namespace
	client         client.Client
	namespacedName types.NamespacedName
	t              test.T
}

func (a *NamespaceAssertion) loadNamespace() error {
	namespace := &corev1.Namespace{}
	err := a.client.Get(context.TODO(), a.namespacedName, namespace)
	a.namespace = namespace
	return err
}

func AssertThatNamespace(t test.T, name string, client client.Client) *NamespaceAssertion {
	return &NamespaceAssertion{
		client:         client,
		namespacedName: types.NamespacedName{Name: name},
		t:              t,
	}
}

func (a *NamespaceAssertion) HasNoOwnerReference() *NamespaceAssertion {
	err := a.loadNamespace()
	require.NoError(a.t, err)
	assert.Empty(a.t, a.namespace.OwnerReferences)
	return a
}

func (a *NamespaceAssertion) HasLabel(key, value string) *NamespaceAssertion {
	err := a.loadNamespace()
	require.NoError(a.t, err)
	require.Contains(a.t, a.namespace.Labels, key)
	assert.Equal(a.t, a.namespace.Labels[key], value)
	return a
}

func (a *NamespaceAssertion) HasResource(name string, obj runtime.Object) *NamespaceAssertion {
	err := a.loadNamespace()
	require.NoError(a.t, err)
	err = a.client.Get(context.TODO(), types.NamespacedName{Namespace: a.namespace.Name, Name: name}, obj)
	require.NoError(a.t, err)
	return a
}
