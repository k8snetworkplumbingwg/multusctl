package client

import (
	"encoding/json"

	net "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (self *Client) CreateNetworkAttachmentDefinition(name string, config string) (*net.NetworkAttachmentDefinition, error) {
	networkAttachmentDefinition := &net.NetworkAttachmentDefinition{
		ObjectMeta: meta.ObjectMeta{
			Name: name,
		},
		Spec: net.NetworkAttachmentDefinitionSpec{
			Config: config,
		},
	}

	if networkAttachmentDefinition, err := self.net.K8sCniCncfIoV1().NetworkAttachmentDefinitions(self.namespace).Create(self.context, networkAttachmentDefinition, meta.CreateOptions{}); err == nil {
		return networkAttachmentDefinition, nil
	} else {
		return nil, err
	}
}

func (self *Client) GetNetworkAttachmentDefinition(name string) (*net.NetworkAttachmentDefinition, error) {
	return self.net.K8sCniCncfIoV1().NetworkAttachmentDefinitions(self.namespace).Get(self.context, name, meta.GetOptions{})
}

func (self *Client) DeleteNetworkAttachmentDefinition(name string) error {
	return self.net.K8sCniCncfIoV1().NetworkAttachmentDefinitions(self.namespace).Delete(self.context, name, meta.DeleteOptions{})
}

func (self *Client) ListNetworkAttachmentDefinitions() (*net.NetworkAttachmentDefinitionList, error) {
	return self.net.K8sCniCncfIoV1().NetworkAttachmentDefinitions(self.namespace).List(self.context, meta.ListOptions{})
}

func (self *Client) ListPodsForNetworkAttachmentDefinition(name string) ([]*core.Pod, error) {
	var r []*core.Pod

	if pods, err := self.kubernetes.CoreV1().Pods(self.namespace).List(self.context, meta.ListOptions{}); err == nil {
		for _, pod := range pods.Items {
			if networks, ok := pod.Annotations[net.NetworkStatusAnnot]; ok {
				var data interface{}
				if err := json.Unmarshal([]byte(networks), &data); err == nil {
					if data_, ok := data.([]interface{}); ok {
						for _, data__ := range data_ {
							if data___, ok := data__.(map[string]interface{}); ok {
								if name_, ok := data___["name"]; ok {
									if name_ == name {
										r = append(r, pod.DeepCopy())
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
