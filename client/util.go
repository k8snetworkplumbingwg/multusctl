package client

import (
	"errors"
	"fmt"
	"time"

	"github.com/tliron/kutil/kubernetes"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	waitpkg "k8s.io/apimachinery/pkg/util/wait"
)

var timeout = 60 * time.Second

func (self *Client) getPods(appName string) (*core.PodList, error) {
	labels_ := labels.Set(map[string]string{
		"app": appName,
	})
	selector := labels_.AsSelector().String()

	if pods, err := self.kubernetes.CoreV1().Pods(self.namespace).List(self.context, meta.ListOptions{LabelSelector: selector}); err == nil {
		if len(pods.Items) > 0 {
			return pods, nil
		} else {
			return nil, fmt.Errorf("no pods for app=\"%s\" in namespace \"%s\"", appName, self.namespace)
		}
	} else {
		return nil, err
	}
}

func (self *Client) waitForDaemonSet(appName string, podAppName string) (*apps.DaemonSet, error) {
	log.Infof("waiting for daemon set for %s", appName)

	var daemonSet *apps.DaemonSet
	err := waitpkg.PollImmediate(time.Second, timeout, func() (bool, error) {
		var err error
		if daemonSet, err = self.kubernetes.AppsV1().DaemonSets(self.namespace).Get(self.context, appName, meta.GetOptions{}); err == nil {
			if daemonSet.Status.NumberReady > 0 {
				return true, nil
			}
			return false, nil
		} else {
			return false, err
		}
	})

	if err == nil {
		log.Infof("daemon set available for %s", appName)
		if err := self.waitForAPod(podAppName, daemonSet); err == nil {
			return daemonSet, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (self *Client) waitForAPod(appName string, daemonSet *apps.DaemonSet) error {
	log.Infof("waiting for a pod for %s", appName)

	return waitpkg.PollImmediate(time.Second, timeout, func() (bool, error) {
		if pods, err := self.getPods(appName); err == nil {
			for _, pod := range pods.Items {
				if self.isPodOwnedBy(&pod, daemonSet) {
					for _, condition := range pod.Status.Conditions {
						switch condition.Type {
						case core.ContainersReady:
							if condition.Status == core.ConditionTrue {
								log.Infof("pod ready for %s: %s", appName, pod.Name)
								return true, nil
							}
						}
					}
				}
			}
			return false, nil
		} else {
			return false, err
		}
	})
}

func (self *Client) isPodOwnedBy(pod *core.Pod, daemonSet *apps.DaemonSet) bool {
	for _, owner := range pod.OwnerReferences {
		if (owner.APIVersion == "apps/v1") && (owner.Kind == "DaemonSet") && (owner.UID == daemonSet.UID) {
			return true
		}
	}
	return false
}

func (self *Client) getSourceRegistryHost(registryHost string) (string, error) {
	if registryHost == "internal" {
		if registryHost, err := kubernetes.GetInternalRegistryHost(self.context, self.kubernetes); err == nil {
			return registryHost, nil
		} else {
			return "", fmt.Errorf("could not discover internal registry: %s", err.Error())
		}
	}

	if registryHost != "" {
		return registryHost, nil
	} else {
		return "", errors.New("must provide \"--registry\"")
	}
}
