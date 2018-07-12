package controller

import (
	"context"
	"fmt"

	"github.com/spotahome/kooper/operator/controller"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
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

// New returns a new controller.
func New(config Config, k8sCli kubernetes.Interface, dynCli dynamic.Interface, logger log.Logger) (*Controller, error) {

	ret := NewNamespaceRetrieve(k8sCli)
	populatorSrv := service.NewConfigMapPopulator(logger, config.Labels, k8sCli, dynCli)
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

func (h *handler) Add(_ context.Context, obj runtime.Object) error {
	namespace, ok := obj.(*corev1.Namespace)
	if !ok {
		return fmt.Errorf("Not a namespace")
	}
	h.populatorSrv.CreateManifests(namespace)
	return nil
}

func (h *handler) Delete(_ context.Context, _ string) error {
	return nil
}
