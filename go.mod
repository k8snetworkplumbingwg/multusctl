module github.com/tliron/multusctl

go 1.16

// replace github.com/tliron/kutil => /Depot/Projects/RedHat/kutil

require (
	github.com/k8snetworkplumbingwg/network-attachment-definition-client v1.1.0
	github.com/spf13/cobra v1.1.3
	github.com/tliron/kutil v0.1.22
	k8s.io/api v0.20.4
	k8s.io/apiextensions-apiserver v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v0.20.4
	k8s.io/klog/v2 v2.6.0
)
