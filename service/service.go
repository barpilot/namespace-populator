package service

import (
	"bytes"
	"text/template"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"

	"github.com/barpilot/namespace-populator/log"
	"github.com/barpilot/namespace-populator/util/namespace"
	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	"github.com/operator-framework/operator-sdk/pkg/util/k8sutil"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// Echo is simple echo service.
type Populator interface {
	CreateManifests(*corev1.Namespace) error
}

// SimpleEcho echoes the received object name.
type ConfigMapPopulator struct {
	logger log.Logger

	configmap string
	k8sCli    kubernetes.Interface
}

type generic map[string]interface{}

// NewSimpleEcho returns a new SimpleEcho.
func NewConfigMapPopulator(logger log.Logger, cm string, k8sCli kubernetes.Interface) *ConfigMapPopulator {
	return &ConfigMapPopulator{
		logger:    logger,
		configmap: cm,
		k8sCli:    k8sCli,
	}
}

func (c *ConfigMapPopulator) CreateManifests(ns *corev1.Namespace) error {

	cm, err := c.k8sCli.CoreV1().ConfigMaps(namespace.Namespace()).Get(c.configmap, metav1.GetOptions{})
	if err != nil {
		return err
	}

	for filename, manifest := range cm.Data {
		var result bytes.Buffer
		//c.logger.Infof(manifest)
		t := template.Must(template.New(filename).Parse(manifest))
		if err := t.Execute(&result, ns); err != nil {
			return err
		}
		//c.logger.Infof(result.String())

		obj := generic{}
		//test, _, _ := runtime.Decoder.Decode($result, nil, nil)
		//c.logger.Infof("%+v", test)

		if err := yaml.NewYAMLOrJSONDecoder(&result, 4096).Decode(&obj); err != nil {
			c.logger.Infof(err.Error())
			return err
		}

		unStrObj := unstructured.Unstructured{Object: obj}
		flagObject(&unStrObj, ns.Name)

		strObj := k8sutil.RuntimeObjectFromUnstructured(&unStrObj)

		if err := action.Create(strObj); err != nil && !apierrors.IsAlreadyExists(err) {
			c.logger.Infof(err.Error())
			return err
		}
	}
	return nil
}

func flagObject(obj metav1.Object, namespace string) {
	annotations := obj.GetAnnotations()
	annotations["namespace-populator.barpilot.io/namespace"] = namespace
	obj.SetAnnotations(annotations)
}
