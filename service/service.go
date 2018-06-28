package service

import (
	"bytes"
	"text/template"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"

	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/barpilot/namespace-populator/log"
)

// Echo is simple echo service.
type Populator interface {
	CreateManifests(*corev1.Namespace) error
}

// SimpleEcho echoes the received object name.
type ConfigMapPopulator struct {
	logger log.Logger

	configmap *corev1.ConfigMap
	k8sCli    kubernetes.Interface
}

// NewSimpleEcho returns a new SimpleEcho.
func NewConfigMapPopulator(logger log.Logger, cm *corev1.ConfigMap, k8sCli kubernetes.Interface) *ConfigMapPopulator {
	return &ConfigMapPopulator{
		logger:    logger,
		configmap: cm,
		k8sCli:    k8sCli,
	}
}

func (c *ConfigMapPopulator) CreateManifests(namespace *corev1.Namespace) error {

	return nil
}

func (c *ConfigMapPopulator) getmanifests(namespace *corev1.Namespace) error {

	for filename, manifest := range c.configmap.Data {
		var result bytes.Buffer
		t := template.Must(template.New(filename).Parse(manifest))
		if err := t.Execute(&result, namespace); err != nil {
			return err
		}
		var obj runtime.Object
		if err := yaml.NewYAMLToJSONDecoder(&result).Decode(obj); err != nil {
			return err
		}
		addAnnotation(obj)

		if err := action.Create(obj); err != nil && !apierrors.IsAlreadyExists(err) {
			return err
		}
	}
	return nil
}

func addAnnotation(obj metav1.Object) {
	annotations := obj.GetAnnotations()
	annotations["toto"] = "tata"
	obj.SetAnnotations(annotations)
}
