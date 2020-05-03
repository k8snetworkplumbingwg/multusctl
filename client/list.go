package client

import (
	"encoding/json"

	net "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (self *Client) ListNetworkAttachmentDefinitions(namespace string) (*net.NetworkAttachmentDefinitionList, error) {
	return self.net.K8sCniCncfIoV1().NetworkAttachmentDefinitions(namespace).List(self.context, meta.ListOptions{})
}

func (self *Client) ListPodsForNetworkAttachmentDefinition(namespace string, name string) ([]*core.Pod, error) {
	var r []*core.Pod

	if pods, err := self.kubernetes.CoreV1().Pods(self.namespace).List(self.context, meta.ListOptions{}); err == nil {
		for _, pod := range pods.Items {
			if networks, ok := pod.Annotations["k8s.v1.cni.cncf.io/networks"]; ok {
				var data interface{}
				if err := json.Unmarshal([]byte(networks), &data); err == nil {
					if data_, ok := data.([]interface{}); ok {
						for _, data__ := range data_ {
							if data___, ok := data__.(map[string]interface{}); ok {
								if name_, ok := data___["name"]; ok {
									if name_ == name {
										r = append(r, &pod)
									}
								}
							}
						}
					}
				} else {
					return nil, err
				}
			}
		}
	} else {
		return nil, err
	}

	return r, nil
}
