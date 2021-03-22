multusctl
=========

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Latest Release](https://img.shields.io/github/release/tliron/multusctl.svg)](https://github.com/tliron/multusctl/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/tliron/multusctl)](https://goreportcard.com/report/github.com/tliron/multusctl)

CLI client for [Multus CNI](https://github.com/k8snetworkplumbingwg/multus-cni).

This client can be used to manage network attachment definitions in a more user-friendly way than working with
`kubectl` directly. E.g.:

    multusctl create myattachment --file my-cni-config.yaml
    multusctl get myattachment

Some quality-of-life features:

* CNI configs are normally in JSON, but this tool allows you to use both JSON and YAML notations.
* `multusctl list` it will show you not only the network attachment definitions but also which pods are attached
  to them.
* Via `multusctl install` and `multusctl uninstall` this utility is all you need to start using Multus.

You can also use this tool as a `kubectl` plugin. Just rename the `multusctl` executable to, say,
`kubectl-multus`, and then you can do this:

    kubectl multus create myattachment --file my-cni-config.yaml

Also see [this issue](https://github.com/k8snetworkplumbingwg/multus-cni/issues/488).
