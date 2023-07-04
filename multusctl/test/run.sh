#!/bin/bash
set -e 

t() {
  local PARENT_IFACE="$1"
  #install sdn
  kubectl apply -f ./assets/kube-flannel.yml

  #create macvlan bridge
  sudo ip link delete data-plane-vlan
  sudo ip link add data-plane-vlan link $PARENT_IFACE type macvlan mode bridge
  sudo ip link set data-plane-vlan up

  #install multus
  ../multusctl uninstall
  ../multusctl install --wait
  #create network attachment called 'netat' 
  ../multusctl create netat --url=assets/config-macvlan.yaml

  #create pod using data-plane-vlan using the new network attachment
  kubectl apply -f ./assets/pod.yml
  
  #veryfi pod has multiple interfaces
  #kubectl exec -it samplepod -- ip a s
}

#HOST_IFACE is the NIC name of the parent interface that will be added to k8s using multus

t "enx24f5a2f13bd3"
