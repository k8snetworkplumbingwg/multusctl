#!/bin/bash

function test () {

  #install sdn
  kubectl create -f ./assets/kube-flannel.yml
  #install multus
  ./multusctl install --wait
  #create network attachment called 'netat' 
  ./multusctl create netat --url=./assets/config-macvlan.yml

  #create macvlan bridge
  sudo ip link add data-plane-vlan link $MASTER_IFACE type macvlan mode bridge
  sudo ip link set data-plane-vlan up

  #create pod using data-plane-vlan using the new network attachment
  kubectl create -f ./assets/pod.yml
}

test $MASTER_IFACE
