package client

import (
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// https://github.com/k8snetworkplumbingwg/multus-cni/blob/master/images/multus-daemonset.yml

func (self *Client) Install(runtime string, sourceRegistryHost string, wait bool) error {
	var err error

	if sourceRegistryHost, err = self.getSourceRegistryHost(sourceRegistryHost); err != nil {
		return err
	}

	if _, err = self.createCustomResourceDefinition(); err != nil {
		return err
	}

	var serviceAccount *core.ServiceAccount
	if serviceAccount, err = self.createServiceAccount(); err != nil {
		return err
	}

	var role *rbac.ClusterRole
	if role, err = self.createClusterRole(); err != nil {
		return err
	}
	if _, err = self.createClusterRoleBinding(role, serviceAccount); err != nil {
		return err
	}

	var configMap *core.ConfigMap
	if configMap, err = self.createConfigMap(); err != nil {
		return err
	}

	var daemonSet *apps.DaemonSet
	if daemonSet, err = self.createDaemonSet(sourceRegistryHost, runtime, configMap, serviceAccount); err != nil {
		return err
	}

	if wait {
		if _, err := self.waitForDaemonSet(daemonSet.Name, "multus"); err != nil {
			return err
		}
	}

	return nil
}

func (self *Client) Uninstall() {
	if err := self.kubernetes.AppsV1().DaemonSets(self.namespace).Delete(self.context, "kube-multus-ds", meta.DeleteOptions{}); err != nil {
		log.Warningf("%s", err)
	}
	if err := self.kubernetes.CoreV1().ConfigMaps(self.namespace).Delete(self.context, "multus-cni-config", meta.DeleteOptions{}); err != nil {
		log.Warningf("%s", err)
	}
	if err := self.kubernetes.RbacV1().ClusterRoleBindings().Delete(self.context, "multus", meta.DeleteOptions{}); err != nil {
		log.Warningf("%s", err)
	}
	if err := self.kubernetes.RbacV1().ClusterRoles().Delete(self.context, "multus", meta.DeleteOptions{}); err != nil {
		log.Warningf("%s", err)
	}
	if err := self.kubernetes.CoreV1().ServiceAccounts(self.namespace).Delete(self.context, "multus", meta.DeleteOptions{}); err != nil {
		log.Warningf("%s", err)
	}
	if err := self.apiExtensions.ApiextensionsV1().CustomResourceDefinitions().Delete(self.context, "network-attachment-definitions.k8s.cni.cncf.io", meta.DeleteOptions{}); err != nil {
		log.Warningf("%s", err)
	}
}

func (self *Client) createCustomResourceDefinition() (*apiextensions.CustomResourceDefinition, error) {
	name := "network-attachment-definitions.k8s.cni.cncf.io"

	customResourceDefinition := &apiextensions.CustomResourceDefinition{
		ObjectMeta: meta.ObjectMeta{
			Name: name,
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{
			Group: "k8s.cni.cncf.io",
			Names: apiextensions.CustomResourceDefinitionNames{
				Singular: "network-attachment-definition",
				Plural:   "network-attachment-definitions",
				Kind:     "NetworkAttachmentDefinition",
				ListKind: "NetworkAttachmentDefinitionList",
				ShortNames: []string{
					"net-attach-def",
				},
				Categories: []string{
					"all", // will appear in "kubectl get all"
				},
			},
			Scope: apiextensions.NamespaceScoped,
			Versions: []apiextensions.CustomResourceDefinitionVersion{
				{
					Name:    "v1",
					Served:  true,
					Storage: true, // one and only one version must be marked with storage=true
					Schema: &apiextensions.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
							Description: "NetworkAttachmentDefinition is a CRD schema specified by the Network Plumbing Working Group to express the intent for attaching pods to one or more logical or physical networks. More information available at: https://github.com/k8snetworkplumbingwg/multi-net-spec",
							Type:        "object",
							Properties: map[string]apiextensions.JSONSchemaProps{
								"apiVersion": {
									Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
									Type:        "string",
								},
								"kind": {
									Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
									Type:        "string",
								},
								"metadata": {
									Type: "object",
								},
								"spec": {
									Description: "NetworkAttachmentDefinition spec defines the desired state of a network attachment",
									Type:        "object",
									Properties: map[string]apiextensions.JSONSchemaProps{
										"config": {
											Description: "NetworkAttachmentDefinition config is a JSON-formatted CNI configuration",
											Type:        "string",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if customResourceDefinition, err := self.apiExtensions.ApiextensionsV1().CustomResourceDefinitions().Create(self.context, customResourceDefinition, meta.CreateOptions{}); err == nil {
		return customResourceDefinition, nil
	} else if errors.IsAlreadyExists(err) {
		return self.apiExtensions.ApiextensionsV1().CustomResourceDefinitions().Get(self.context, name, meta.GetOptions{})
	} else {
		return nil, err
	}
}

func (self *Client) createServiceAccount() (*core.ServiceAccount, error) {
	name := "multus"

	serviceAccount := &core.ServiceAccount{
		ObjectMeta: meta.ObjectMeta{
			Name: name,
		},
	}

	if serviceAccount, err := self.kubernetes.CoreV1().ServiceAccounts(self.namespace).Create(self.context, serviceAccount, meta.CreateOptions{}); err == nil {
		return serviceAccount, nil
	} else if errors.IsAlreadyExists(err) {
		return self.kubernetes.CoreV1().ServiceAccounts(self.namespace).Get(self.context, name, meta.GetOptions{})
	} else {
		return nil, err
	}
}

func (self *Client) createClusterRole() (*rbac.ClusterRole, error) {
	name := "multus"

	clusterRole := &rbac.ClusterRole{
		ObjectMeta: meta.ObjectMeta{
			Name: name,
		},
		Rules: []rbac.PolicyRule{
			{
				APIGroups: []string{"k8s.cni.cncf.io"},
				Resources: []string{rbac.ResourceAll},
				Verbs:     []string{rbac.VerbAll},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"pods", "pods/status"},
				Verbs:     []string{"get", "update"},
			},
			{
				APIGroups: []string{"", "events.k8s.io"},
				Resources: []string{"events"},
				Verbs:     []string{"create", "patch", "update"},
			},
		},
	}

	if clusterRole, err := self.kubernetes.RbacV1().ClusterRoles().Create(self.context, clusterRole, meta.CreateOptions{}); err == nil {
		return clusterRole, nil
	} else if errors.IsAlreadyExists(err) {
		return self.kubernetes.RbacV1().ClusterRoles().Get(self.context, name, meta.GetOptions{})
	} else {
		return nil, err
	}
}

func (self *Client) createClusterRoleBinding(role *rbac.ClusterRole, serviceAccount *core.ServiceAccount) (*rbac.ClusterRoleBinding, error) {
	name := "multus"

	clusterRoleBinding := &rbac.ClusterRoleBinding{
		ObjectMeta: meta.ObjectMeta{
			Name: name,
		},
		Subjects: []rbac.Subject{
			{
				Kind:      rbac.ServiceAccountKind, // serviceAccount.Kind is empty
				Name:      serviceAccount.Name,
				Namespace: self.namespace, // required
			},
		},
		RoleRef: rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "ClusterRole",
			Name:     role.Name,
		},
	}

	if clusterRoleBinding, err := self.kubernetes.RbacV1().ClusterRoleBindings().Create(self.context, clusterRoleBinding, meta.CreateOptions{}); err == nil {
		return clusterRoleBinding, nil
	} else if errors.IsAlreadyExists(err) {
		return self.kubernetes.RbacV1().ClusterRoleBindings().Get(self.context, name, meta.GetOptions{})
	} else {
		return nil, err
	}
}

func (self *Client) createConfigMap() (*core.ConfigMap, error) {
	name := "multus-cni-config"

	configMap := &core.ConfigMap{
		ObjectMeta: meta.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"tier": "node",
				"app":  "multus",
			},
		},
		Data: map[string]string{
			"cni-conf.json": `
{
  "name": "multus-cni-network",
  "type": "multus",
  "capabilities": {
    "portMappings": true
  },
  "delegates": [
    {
      "cniVersion": "0.3.1",
      "name": "default-cni-network",
      "plugins": [
        {
          "type": "flannel",
          "name": "flannel.1",
            "delegate": {
              "isDefaultGateway": true,
              "hairpinMode": true
            }
          },
          {
            "type": "portmap",
            "capabilities": {
              "portMappings": true
            }
          }
      ]
    }
  ],
  "kubeconfig": "/etc/cni/net.d/multus.d/multus.kubeconfig"
}
`,
		},
	}

	if configMap, err := self.kubernetes.CoreV1().ConfigMaps(self.namespace).Create(self.context, configMap, meta.CreateOptions{}); err == nil {
		return configMap, nil
	} else if errors.IsAlreadyExists(err) {
		return self.kubernetes.CoreV1().ConfigMaps(self.namespace).Get(self.context, name, meta.GetOptions{})
	} else {
		return nil, err
	}
}

func (self *Client) createDaemonSet(sourceRegistryHost string, runtime string, configMap *core.ConfigMap, serviceAccount *core.ServiceAccount) (*apps.DaemonSet, error) {
	name := "kube-multus-ds"

	true_ := true
	var terminationGracePeriodSeconds int64 = 10
	daemonSet := &apps.DaemonSet{
		ObjectMeta: meta.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"tier": "node",
				"app":  "multus",
				"name": "multus",
			},
		},
		Spec: apps.DaemonSetSpec{
			Selector: &meta.LabelSelector{
				MatchLabels: map[string]string{
					"name": "multus",
				},
			},
			UpdateStrategy: apps.DaemonSetUpdateStrategy{
				Type: apps.RollingUpdateDaemonSetStrategyType,
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: map[string]string{
						"tier": "node",
						"app":  "multus",
						"name": "multus",
					},
				},
				Spec: core.PodSpec{
					HostNetwork: true,
					Tolerations: []core.Toleration{
						{
							Operator: core.TolerationOpExists,
							Effect:   core.TaintEffectNoSchedule,
						},
					},
					ServiceAccountName: serviceAccount.Name,
					Containers: []core.Container{
						{
							Name:  "kube-multus",
							Image: sourceRegistryHost + "/k8snetworkplumbingwg/multus-cni:stable",
							Command: []string{
								"/entrypoint.sh",
							},
							Resources: core.ResourceRequirements{
								Requests: core.ResourceList{
									core.ResourceCPU:    resource.MustParse("100m"),
									core.ResourceMemory: resource.MustParse("50Mi"),
								},
								Limits: core.ResourceList{
									core.ResourceCPU:    resource.MustParse("100m"),
									core.ResourceMemory: resource.MustParse("50Mi"),
								},
							},
							SecurityContext: &core.SecurityContext{
								Privileged: &true_,
							},
						},
					},
					TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
				},
			},
		},
	}

	switch runtime {
	case "crio":
		daemonSet.Spec.Template.Spec.Containers[0].SecurityContext.Capabilities = &core.Capabilities{
			Add: []core.Capability{"SYS_ADMIN"},
		}

		daemonSet.Spec.Template.Spec.Containers[0].Args = []string{
			"--cni-version=0.3.1",
			"--cni-bin-dir=/host/usr/libexec/cni",
			"--multus-conf-file=auto",
			"--restart-crio=true",
		}

		daemonSet.Spec.Template.Spec.Containers[0].VolumeMounts = []core.VolumeMount{
			{
				Name:      "run",
				MountPath: "/run",
			},
			{
				Name:      "cni",
				MountPath: "/host/etc/cni/net.d",
			},
			{
				Name:      "cnibin",
				MountPath: "/host/usr/libexec/cni",
			},
			{
				Name:      "multus-cfg",
				MountPath: "/tmp/multus-conf",
			},
		}

		daemonSet.Spec.Template.Spec.Volumes = []core.Volume{
			{
				Name: "run",
				VolumeSource: core.VolumeSource{
					HostPath: &core.HostPathVolumeSource{
						Path: "/run",
					},
				},
			},
			{
				Name: "cni",
				VolumeSource: core.VolumeSource{
					HostPath: &core.HostPathVolumeSource{
						Path: "/etc/cni/net.d",
					},
				},
			},
			{
				Name: "cnibin",
				VolumeSource: core.VolumeSource{
					HostPath: &core.HostPathVolumeSource{
						Path: "/usr/libexec/cni",
					},
				},
			},
			{
				Name: "multus-cfg",
				VolumeSource: core.VolumeSource{
					ConfigMap: &core.ConfigMapVolumeSource{
						LocalObjectReference: core.LocalObjectReference{
							Name: configMap.Name,
						},
						Items: []core.KeyToPath{
							{
								Key:  "cni-conf.json",
								Path: "70-multus.conf",
							},
						},
					},
				},
			},
		}

	default:
		daemonSet.Spec.Template.Spec.Containers[0].Args = []string{
			"--multus-conf-file=auto",
			"--cni-version=0.3.1",
		}

		daemonSet.Spec.Template.Spec.Containers[0].VolumeMounts = []core.VolumeMount{
			{
				Name:      "cni",
				MountPath: "/host/etc/cni/net.d",
			},
			{
				Name:      "cnibin",
				MountPath: "/host/opt/cni/bin",
			},
			{
				Name:      "multus-cfg",
				MountPath: "/tmp/multus-conf",
			},
		}

		daemonSet.Spec.Template.Spec.Volumes = []core.Volume{
			{
				Name: "cni",
				VolumeSource: core.VolumeSource{
					HostPath: &core.HostPathVolumeSource{
						Path: "/etc/cni/net.d",
					},
				},
			},
			{
				Name: "cnibin",
				VolumeSource: core.VolumeSource{
					HostPath: &core.HostPathVolumeSource{
						Path: "/opt/cni/bin",
					},
				},
			},
			{
				Name: "multus-cfg",
				VolumeSource: core.VolumeSource{
					ConfigMap: &core.ConfigMapVolumeSource{
						LocalObjectReference: core.LocalObjectReference{
							Name: configMap.Name,
						},
						Items: []core.KeyToPath{
							{
								Key:  "cni-conf.json",
								Path: "70-multus.conf",
							},
						},
					},
				},
			},
		}
	}

	if daemonSet, err := self.kubernetes.AppsV1().DaemonSets(self.namespace).Create(self.context, daemonSet, meta.CreateOptions{}); err == nil {
		return daemonSet, nil
	} else if errors.IsAlreadyExists(err) {
		return self.kubernetes.AppsV1().DaemonSets(self.namespace).Get(self.context, name, meta.GetOptions{})
	} else {
		return nil, err
	}
}
