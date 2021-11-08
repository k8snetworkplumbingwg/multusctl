multusctl
=========

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Latest Release](https://img.shields.io/github/release/k8snetworkplumbingwg/multusctl.svg)](https://github.com/k8snetworkplumbingwg/multusctl/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/k8snetworkplumbingwg/multusctl)](https://goreportcard.com/report/github.com/k8snetworkplumbingwg/multusctl)

CLI tool for [Multus CNI](https://github.com/k8snetworkplumbingwg/multus-cni).

`multusctl` can be used to install and uninstall Multus. Additionally, it provides a more
user-friendly way to work with network attachment definitions, as compared to using `kubectl`
directly.

Additional features:

* CNI configs are normally in JSON, but `multusctl` allows you to use YAML (the default), JSON,
  XML, and even CBOR.
* When listing network attachement definitions it will also show you the resources that are
  attached to it.

Examples:

    multusctl install --wait
    multusctl create myattachment --url=assets/config.yaml
    multusctl list
    multusctl get myattachment --format=json

The CNI config can also be provided via stdin:

    cat assets/config.yaml | multusctl create myattachment

You can also use this tool as a `kubectl` plugin. Just rename the `multusctl` executable to, say,
`kubectl-nad`, and then you can do this:

    kubectl nad create myattachment --url=assets/config.yaml
