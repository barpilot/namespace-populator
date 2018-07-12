package service

import (
	"bytes"
	"text/template"

	corev1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	"github.com/barpilot/namespace-populator/log"
	"github.com/barpilot/namespace-populator/util/namespace"
)

const (
	AnnotationIndicator = "namespace-populator.barpilot.io/namespace"
	AnnotationSelector  = "namespace-populator.barpilot.io/selector"
)

// Populator is a populator service.
type Populator interface {
	CreateManifests(*corev1.Namespace) error
}

// ConfigMapPopulator echoes the received object name.
type ConfigMapPopulator struct {
	logger log.Logger

	labels string
	k8sCli kubernetes.Interface
	dynCli dynamic.Interface
}

type generic map[string]interface{}

// NewConfigMapPopulator returns a new ConfigMapPopulator.
func NewConfigMapPopulator(logger log.Logger, labels string, k8sCli kubernetes.Interface, dynCli dynamic.Interface) *ConfigMapPopulator {
	return &ConfigMapPopulator{
		logger: logger,
		labels: labels,
		k8sCli: k8sCli,
		dynCli: dynCli,
	}
}

func (c *ConfigMapPopulator) CreateManifests(ns *corev1.Namespace) error {
	l, err := c.k8sCli.CoreV1().ConfigMaps(namespace.Namespace()).List(metav1.ListOptions{LabelSelector: c.labels})
	if err != nil {
		return err
	}

	for _, cm := range l.Items {

		if !validateConfigMap(&cm, ns) {
			continue
		}

		for filename, manifest := range cm.Data {
			var result bytes.Buffer
			t := template.Must(template.New(filename).Parse(manifest))
			if err := t.Execute(&result, ns); err != nil {
				return err
			}

			obj := generic{}

			if err := yaml.NewYAMLOrJSONDecoder(&result, 4096).Decode(&obj); err != nil {
				c.logger.Infof(err.Error())
				return err
			}

			unStrObj := unstructured.Unstructured{Object: obj}
			flagObject(&unStrObj, ns.Name)

			//c.logger.Infof(unStrObj.GetName())

			_, err := c.dynCli.Resource(c.groupVersionResource(unStrObj)).Namespace(unStrObj.GetNamespace()).Create(&unStrObj)
			if err != nil && !apierrors.IsAlreadyExists(err) {
				c.logger.Infof(err.Error())
				return err
			}
		}
	}
	return nil
}

func (c *ConfigMapPopulator) groupVersionResource(unStrObj unstructured.Unstructured) schema.GroupVersionResource {
	gvr, _ := meta.UnsafeGuessKindToResource(unStrObj.GroupVersionKind())
	//c.logger.Infof(gvr.String())
	return gvr
}

func validateConfigMap(cm *corev1.ConfigMap, ns *corev1.Namespace) bool {
	if selector, ok := cm.Annotations[AnnotationSelector]; ok {
		ls, err := metav1.ParseToLabelSelector(selector)
		if err != nil {
			return false
		}

		s, err := metav1.LabelSelectorAsSelector(ls)
		if err != nil {
			return false
		}
		if !s.Matches(labels.Set(ns.Labels)) {
			return false
		}
	}
	return true
}

func flagObject(obj *unstructured.Unstructured, namespace string) {
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations[AnnotationIndicator] = namespace
	obj.SetAnnotations(annotations)
}
