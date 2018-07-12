package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	applogger "github.com/spotahome/kooper/log"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/barpilot/namespace-populator/controller"
	"github.com/barpilot/namespace-populator/log"
)

// Main is the main program.
type Main struct {
	flags  *Flags
	config controller.Config
	logger log.Logger
}

// New returns the main application.
func New(logger log.Logger) *Main {
	f := NewFlags()
	return &Main{
		flags:  f,
		config: f.ControllerConfig(),
		logger: logger,
	}
}

// Run runs the app.
func (m *Main) Run(stopC <-chan struct{}) error {
	m.logger.Infof("initializing namespace-populator")

	// Get kubernetes rest client.
	k8sCli, dynCli, err := m.getKubernetesClient()
	if err != nil {
		return err
	}

	// Create the controller and run
	ctrl, err := controller.New(m.config, k8sCli, dynCli, m.logger)
	if err != nil {
		return err
	}

	return ctrl.Run(stopC)
}

func (m *Main) getKubernetesClient() (kubernetes.Interface, dynamic.Interface, error) {
	var err error
	var cfg *rest.Config

	// If devel mode then use configuration flag path.
	if m.flags.Development {
		cfg, err = clientcmd.BuildConfigFromFlags("", m.flags.KubeConfig)
		if err != nil {
			return nil, nil, fmt.Errorf("could not load configuration: %s", err)
		}
	} else {
		cfg, err = rest.InClusterConfig()
		if err != nil {
			return nil, nil, fmt.Errorf("error loading kubernetes configuration inside cluster, check app is running outside kubernetes cluster or run in development mode: %s", err)
		}
	}

	dynCli, err := dynamic.NewForConfig(cfg)
	kubeCli, err := kubernetes.NewForConfig(cfg)

	return kubeCli, dynCli, nil
}

func main() {
	logger := &applogger.Std{}

	stopC := make(chan struct{})
	finishC := make(chan error)
	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, syscall.SIGTERM, syscall.SIGINT)
	m := New(logger)

	// Run in background the controller.
	go func() {
		finishC <- m.Run(stopC)
	}()

	select {
	case err := <-finishC:
		if err != nil {
			fmt.Fprintf(os.Stderr, "error running controller: %s", err)
			os.Exit(1)
		}
	case <-signalC:
		logger.Infof("Signal captured, exiting...")
	}

}
