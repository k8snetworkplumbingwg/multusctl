package client

import (
	"context"

	netpkg "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/clientset/versioned"
	"github.com/tliron/kutil/logging"
	apiextensionspkg "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kubernetespkg "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var log = logging.GetLogger("multusctl.client")

type Client struct {
	config        *rest.Config
	kubernetes    *kubernetespkg.Clientset
	apiExtensions *apiextensionspkg.Clientset
	net           *netpkg.Clientset
	namespace     string
	context       context.Context
}

func NewClient(masterUrl string, kubeconfigPath string, namespace string) (*Client, error) {
	config, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
	if err != nil {
		return nil, err
	}

	kubernetes, err := kubernetespkg.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	apiExtensions, err := apiextensionspkg.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	net, err := netpkg.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		config:        config,
		kubernetes:    kubernetes,
		apiExtensions: apiExtensions,
		net:           net,
		namespace:     namespace,
		context:       context.TODO(),
	}, nil
}
