package controller

import (
	"fmt"

	"github.com/spotahome/kooper/operator/controller"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"

	"github.com/barpilot/namespace-populator/log"
	"github.com/barpilot/namespace-populator/service"
)

// Controller is a controller that echoes pod events.
type Controller struct {
	controller.Controller

	config Config
	logger log.Logger
}

// New returns a new Echo controller.
func New(config Config, k8sCli kubernetes.Interface, logger log.Logger) (*Controller, error) {

	ret := NewNamespaceRetrieve(k8sCli)
	populatorSrv := service.NewConfigMapPopulator(logger, config.Configmaps[0], k8sCli)
	handler := &handler{populatorSrv: populatorSrv, logger: logger}

	ctrl := controller.NewSequential(config.ResyncPeriod, handler, ret, nil, logger)

	return &Controller{
		Controller: ctrl,
		config:     config,
		logger:     logger,
	}, nil
}

const (
	addPrefix    = "ADD"
	deletePrefix = "DELETE"
)

type handler struct {
	populatorSrv service.Populator
	logger       log.Logger
}

func (h *handler) Add(obj runtime.Object) error {
	namespace, ok := obj.(*corev1.Namespace)
	if !ok {
		return fmt.Errorf("Not a namespace")
	}
	h.logger.Infof("youhouuo a new namespace")
	h.populatorSrv.CreateManifests(namespace)
	return nil
}

func (h *handler) Delete(s string) error {
	return nil
}
