// Copyright 2020 rvazquez@redhat.com
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package webhook

import (
	"fmt"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

type parameter struct {
	name         string
	annotation   string
	defaultValua string
}

const (

	// parameter names
	paramNodeID            = "node-id"
	paramClusterID         = "cluster-id"
	paramContainerName     = "container-name"
	paramPorts             = "ports"
	paramHostPortMapings   = "host-port-mappings"
	paramImage             = "envoy-image"
	paramADSConfigMap      = "ads-configmap"
	paramConfigVolume      = "config-volume"
	paramTLSVolume         = "tls-volume"
	paramClientCertificate = "client-certificate"
	paramEnvoyExtraArgs    = "envoy-extra-args"

	// default values
	DefaultContainerName       = "envoy-sidecar"
	DefaultImage               = "envoyproxy/envoy:v1.14.1"
	DefaultBootstrapConfigMap  = "envoy-sidecar-bootstrap"
	DefaultConfigVolume        = "envoy-sidecar-bootstrap"
	DefaultTLSVolume           = "envoy-sidecar-tls"
	DefaultClientCertificate   = "envoy-sidecar-client-cert"
	DefaultEnvoyExtraArgs      = ""
	DefaultEnvoyConfigBasePath = "/etc/envoy/bootstrap"
	DefaultEnvoyConfigFileName = "config.json"
	DefaultEnvoyTLSBasePath    = "/etc/envoy/tls/client"

	marin3rAnnotationsDomain = "marin3r.3scale.net"
)

type envoySidecarConfig struct {
	name             string
	image            string
	ports            []corev1.ContainerPort
	adsConfigMap     string
	nodeID           string
	clusterID        string
	tlsVolume        string
	configVolume     string
	clientCertSecret string
	extraArgs        string
}

func getStringParam(key string, annotations map[string]string) string {

	var defaults = map[string]string{
		paramContainerName:     DefaultContainerName,
		paramImage:             DefaultImage,
		paramADSConfigMap:      DefaultBootstrapConfigMap,
		paramConfigVolume:      DefaultConfigVolume,
		paramTLSVolume:         DefaultTLSVolume,
		paramClientCertificate: DefaultClientCertificate,
		paramEnvoyExtraArgs:    DefaultEnvoyExtraArgs,
	}

	// return the value specified in the corresponding annotation, if any
	if value, ok := annotations[fmt.Sprintf("%s/%s", marin3rAnnotationsDomain, key)]; ok {
		return value
	}

	// the default for cluster-id is to be set to the node-id value, which
	// has no default but is always set by the user, otherwise we won't get
	// this far as there are previous checks that validate that node-id is set
	if key == "cluster-id" {
		return getNodeID(annotations)
	}

	// return a default value
	return defaults[key]
}

func getNodeID(annotations map[string]string) string {
	// the node-id annotation is always present, otherwise
	// mutation won't even be triggered
	return annotations[fmt.Sprintf("%s/%s", marin3rAnnotationsDomain, paramNodeID)]
}

func (esc *envoySidecarConfig) PopulateFromAnnotations(annotations map[string]string) error {

	esc.name = getStringParam(paramContainerName, annotations)
	esc.image = getStringParam(paramImage, annotations)

	ports, err := getContainerPorts(annotations)
	if err != nil {
		return err
	}
	esc.ports = ports
	esc.adsConfigMap = getStringParam(paramADSConfigMap, annotations)
	esc.nodeID = getNodeID(annotations)
	esc.clusterID = getStringParam(paramClusterID, annotations)
	esc.configVolume = getStringParam(paramConfigVolume, annotations)
	esc.tlsVolume = getStringParam(paramTLSVolume, annotations)
	esc.clientCertSecret = getStringParam(paramClientCertificate, annotations)
	esc.extraArgs = getStringParam(paramEnvoyExtraArgs, annotations)

	return nil
}

// port spec format is "name:port[:protocol],name:port[:protcol]"
func getContainerPorts(annotations map[string]string) ([]corev1.ContainerPort, error) {
	plist := []corev1.ContainerPort{}

	if ports, ok := annotations[fmt.Sprintf("%s/%s", marin3rAnnotationsDomain, paramPorts)]; ok {

		for _, containerPort := range strings.Split(ports, ",") {

			p, err := parsePortSpec(containerPort)
			if err != nil {
				return []corev1.ContainerPort{}, err
			}

			hp, err := hostPortMapping(p.Name, annotations)
			if err != nil {
				return []corev1.ContainerPort{}, err
			}
			if hp != 0 {
				p.HostPort = hp
			}
			plist = append(plist, *p)

		}

	} else {
		// no ports defined for envoy sidecar
		return []corev1.ContainerPort{}, nil
	}

	return plist, nil
}

func parsePortSpec(spec string) (*corev1.ContainerPort, error) {
	var port corev1.ContainerPort

	params := strings.Split(spec, ":")
	// Each port spec should at least contain name and port number
	if len(params) < 2 {
		return nil, fmt.Errorf("Incorrect format, the por specification format for the envoy sidecar container is 'name:port[:protocol]'")
	}

	port.Name = params[0]
	p, err := portNumber(params[1])
	if err != nil {
		return nil, err
	}
	port.ContainerPort = int32(p)

	// Check that protocol matches one of the allowed protocols
	if len(params) == 3 {
		if params[2] == "TCP" || params[2] == "UDP" || params[2] == "SCTP" {
			port.Protocol = corev1.Protocol(params[2])

		} else {
			return nil, fmt.Errorf("Unsupported port protocol '%s'", params[2])
		}
	}

	return &port, nil

}

func hostPortMapping(portName string, annotations map[string]string) (int32, error) {

	if specs, ok := annotations[fmt.Sprintf("%s/%s", marin3rAnnotationsDomain, paramHostPortMapings)]; ok {
		for _, spec := range strings.Split(specs, ",") {
			params := strings.Split(spec, ":")
			if len(params) != 2 {
				return 0, fmt.Errorf("Incorrect number of params in host-port-mapping spec '%v'", spec)
			}
			if params[0] == portName {
				p, err := portNumber(params[1])
				if err != nil {
					return 0, err
				}
				return p, nil
			}
		}
	}
	return 0, nil
}

func portNumber(sport string) (int32, error) {

	// Parse and validate the port number
	iport, err := strconv.Atoi(sport)
	if err != nil {
		return 0, fmt.Errorf("%v doesn't look like a port number, check your port specs", sport)
	}
	// Check if port is within allowed ranges. Privileged ports are not allowed
	if iport < 1024 || iport > 65535 {
		return 0, fmt.Errorf("Port number %v is not in the range 1024-65535", iport)
	}
	return int32(iport), nil
}

func (esc *envoySidecarConfig) container() corev1.Container {

	container := corev1.Container{
		Name:    esc.name,
		Image:   esc.image,
		Command: []string{"envoy"},
		Args: []string{
			"-c",
			fmt.Sprintf("%s.%s", DefaultEnvoyConfigBasePath, DefaultEnvoyConfigFileName),
			"--service-node",
			esc.nodeID,
			"--service-cluster",
			esc.clusterID,
		},
		Ports: esc.ports,
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      esc.tlsVolume,
				ReadOnly:  true,
				MountPath: DefaultEnvoyTLSBasePath,
			},
			{
				Name:      esc.configVolume,
				ReadOnly:  true,
				MountPath: DefaultEnvoyConfigBasePath,
			},
		},
	}

	if esc.extraArgs != "" {
		for _, arg := range strings.Split(esc.extraArgs, " ") {
			container.Args = append(container.Args, arg)
		}
	}

	return container
}

func (esc *envoySidecarConfig) volumes() []corev1.Volume {

	volumes := []corev1.Volume{
		{
			Name: esc.tlsVolume,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: esc.clientCertSecret,
				},
			},
		},
		{
			Name: esc.configVolume,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: esc.adsConfigMap,
					},
				},
			},
		},
	}

	return volumes
}
