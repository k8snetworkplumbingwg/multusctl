multusctl
=========

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Latest Release](https://img.shields.io/github/release/tliron/multusctl.svg)](https://github.com/tliron/multusctl/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/tliron/multusctl)](https://goreportcard.com/report/github.com/tliron/multusctl)

CLI client for [Multus CNI](https://github.com/k8snetworkplumbingwg/multus-cni).

This client can be used to manage network attachment definitions in a more user-friendly way than
working with `kubectl` directly. CNI configs are normally in JSON, but this tool allows you to use both
JSON and YAML notations.

Example:

    multusctl install --wait
    multusctl create myattachment --url=assets/config.yaml
    multusctl list
    multusctl get myattachment

The CNI config can also be provided via stdin:

    cat assets/config.yaml | multusctl create myattachment

Note that `multusctl list` it will show you not only the network attachment definitions but also which
pods are attached to them.

You can also use this tool as a `kubectl` plugin. Just rename the `multusctl` executable to, say,
`kubectl-nad`, and then you can do this:

    kubectl nad create myattachment --url=assets/config.yaml

Also see [this issue](https://github.com/k8snetworkplumbingwg/multus-cni/issues/488).
