multusctl
=========

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Latest Release](https://img.shields.io/github/release/k8snetworkplumbingwg/multusctl.svg)](https://github.com/k8snetworkplumbingwg/multusctl/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/k8snetworkplumbingwg/multusctl)](https://goreportcard.com/report/github.com/k8snetworkplumbingwg/multusctl)

CLI tool for [Multus CNI](https://github.com/k8snetworkplumbingwg/multus-cni).

`multusctl` provides a more user-friendly way to work with network attachment definitions as
compared to using `kubectl` directly. Additionally, it can be used to install and uninstall Multus.

In particular is solves the problem of escaping/unescaping the JSON CNI configuration embedded
in the network attachment definition custom resource.

Additional features:

* CNI configurations are normally in JSON, but `multusctl` also allows you to use YAML (the
  default), XML, and even CBOR.
* Annotation awareness: when listing network attachement definitions it will also show you
  the resources that are attached to it.

Examples
--------

    multusctl install --wait
    multusctl create myattachment --url=assets/config.yaml
    multusctl list
    multusctl get myattachment --format=json

The CNI configuration can also be provided via stdin:

    cat assets/config.yaml | multusctl create myattachment

Installation
------------

One-liner to download the latest Linux AMD64 version and install it in `/usr/bin/` (requires
curl, grep, sed, and tar):

    VERSION=$(curl --silent https://api.github.com/repos/k8snetworkplumbingwg/multusctl/releases/latest | grep '"tag_name":' | sed --regexp-extended 's/.*"([^"]+)".*/\1/') curl --silent --location https://github.com/k8snetworkplumbingwg/multusctl/releases/download/$VERSION/multusctl_${VERSION:1}_linux_amd64.tar.gz | tar --directory=/usr/bin --extract --gzip multusctl

One-liner to install bash completions for the current user (for when you press TAB in bash):

  mkdir --parents ~/.local/share/bash-completion/completions/ && multusctl bash > ~/.local/share/bash-completion/completions/multusctl && exec bash

You can also use this tool as a `kubectl` plugin. Just rename the `multusctl` executable to, say,
`kubectl-multus`:

    mv /usr/bin/multusctl /usr/bin/kubectl-multus

and then you can do this:

    kubectl multus create myattachment --url=assets/config.yaml
