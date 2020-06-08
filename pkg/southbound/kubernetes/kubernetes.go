// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package kubernetes

import (
	"fmt"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
)

var log = logging.GetLogger("southbound", "kubernetes")

// NamespaceEnv is the environment variable for setting the k8s namespace
const NamespaceEnv = "NAMESPACE"

// ServiceNameEnv is the environment variable for setting the k8s service for ran-simulator
const ServiceNameEnv = "SERVICENAME"

// AddK8SServicePorts add a Port to the K8s service
func AddK8SServicePorts(ports []uint16) error {
	namespace := os.Getenv(NamespaceEnv)
	serviceName := os.Getenv(ServiceNameEnv)
	if serviceName == "" {
		serviceName = "ran-simulator"
	}

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Errorf("Failed to access cluster config %s", err.Error())
		return err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Errorf("Failed to create client set %s", err.Error())
		return err
	}
	thisService, err := clientset.CoreV1().Services(namespace).Get(serviceName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		log.Errorf("Service %s not found in namespace %s", serviceName, namespace)
		return nil
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Error("Error getting Service %s in namespace %s. Status: %v",
			serviceName, namespace, statusError.ErrStatus.Message)
		return err
	} else if err != nil {
		log.Errorf("Kubernetes API error %s", err.Error())
		return err
	}

	for _, p := range ports {
		newPort := v1.ServicePort{
			Name:       fmt.Sprintf("e2port%d", p),
			Protocol:   "TCP",
			Port:       int32(p),
			TargetPort: intstr.IntOrString{IntVal: int32(p)},
		}
		thisService.Spec.Ports = append(thisService.Spec.Ports, newPort)
	}
	log.Infof("Service %s - appended %d ports", serviceName, len(ports))
	_, err = clientset.CoreV1().Services(namespace).Update(thisService)
	if statusError, isStatus := err.(*errors.StatusError); isStatus {
		// The ports may already exist if the ran-simulator pod is restarting
		log.Infof("Error updating %s:%s. Status: %v %v",
			namespace, serviceName, statusError.ErrStatus.Reason, statusError.ErrStatus.Message)
		return nil
	} else if err != nil {
		log.Errorf("Kubernetes API error when replacing service %s %s", serviceName, err.Error())
		return err
	}
	return nil
}
