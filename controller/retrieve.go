package controller

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// NamespaceRetrieve knows how to retrieve pods.
type NamespaceRetrieve struct {
	//namespace string
	client kubernetes.Interface
}

// NewNamespaceRetrieve returns a new namespace retriever.
func NewNamespaceRetrieve(client kubernetes.Interface) *NamespaceRetrieve {
	return &NamespaceRetrieve{
		//namespace: namespace,
		client: client,
	}
}

// GetListerWatcher knows how to return a listerWatcher of a namespace.
func (p *NamespaceRetrieve) GetListerWatcher() cache.ListerWatcher {

	return &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return p.client.CoreV1().Namespaces().List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return p.client.CoreV1().Namespaces().Watch(options)
		},
	}
}

// GetObject returns the empty namespace.
func (p *NamespaceRetrieve) GetObject() runtime.Object {
	return &corev1.Namespace{}
}
