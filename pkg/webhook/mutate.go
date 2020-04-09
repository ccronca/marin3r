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

	"go.uber.org/zap"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	podResource = metav1.GroupVersionResource{Version: "v1", Resource: "pods"}
)

func MutatePod(req *admissionv1.AdmissionRequest, logger *zap.SugaredLogger) ([]patchOperation, error) {

	// This handler should only get called on Pod objects as per the MutatingWebhookConfiguration in the YAML file.
	// However, if (for whatever reason) this gets invoked on an object of a different kind, issue a log message but
	// let the object request pass through otherwise.
	if req.Resource != podResource {
		logger.Warnf("expect resource to be %s", podResource)
		return nil, nil
	}

	// Parse the Pod object.
	raw := req.Object.Raw
	pod := corev1.Pod{}
	if _, _, err := universalDeserializer.Decode(raw, nil, &pod); err != nil {
		return nil, fmt.Errorf("could not deserialize pod object: %v", err)
	}

	logger.Infof("AdmissionReview for Kind=%v, Namespace=%v Name=%v (%v) UID=%v patchOperation=%v UserInfo=%v",
		req.Kind, req.Namespace, req.Name, pod.Name, req.UID, req.Operation, req.UserInfo)

	if _, ok := pod.GetAnnotations()[fmt.Sprintf("%s/%s", marin3rAnnotationsDomain, paramNodeID)]; !ok {
		logger.Infof("skipping mutation for %s/%s due to missing '%s' annotation", pod.Namespace,
			pod.Name, fmt.Sprintf("%s/%s", marin3rAnnotationsDomain, paramNodeID))
		return nil, nil
	}

	// Init the list of patches
	var patches []patchOperation

	// Get the patches for the envoy sidecar container
	config := envoySidecarConfig{}
	err := config.PopulateFromAnnotations(pod.GetAnnotations())
	if err != nil {
		return []patchOperation{}, err
	}

	patches = append(patches, patchOperation{
		// "/-" refers to the end of an array in jsonPatch
		Path:  "/spec/containers/-",
		Op:    "add",
		Value: config.container(),
	})

	volumes := config.volumes()
	for _, volume := range volumes {
		patches = append(patches, patchOperation{
			// "/-" refers to the end of an array in jsonPatch
			Path:  "/spec/volumes/-",
			Op:    "add",
			Value: volume,
		})
	}

	return patches, nil
}