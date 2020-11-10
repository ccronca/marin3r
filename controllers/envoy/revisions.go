package controllers

import (
	"context"
	"fmt"
	"hash/fnv"
	"sort"

	envoyv1alpha1 "github.com/3scale/marin3r/apis/envoy/v1alpha1"
	common "github.com/3scale/marin3r/pkg/common"
	envoy "github.com/3scale/marin3r/pkg/envoy"

	"github.com/operator-framework/operator-lib/status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	nodeIDTag    = "marin3r.3scale.net/node-id"
	versionTag   = "marin3r.3scale.net/config-version"
	envoyAPITag  = "marin3r.3scale.net/envoy-api"
	maxRevisions = 10
)

func (r *EnvoyConfigReconciler) ensureEnvoyConfigRevision(ctx context.Context,
	ec *envoyv1alpha1.EnvoyConfig, version string) error {

	// Get the list of revisions for the current version
	ecrList := &envoyv1alpha1.EnvoyConfigRevisionList{}
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: map[string]string{
			nodeIDTag:   ec.Spec.NodeID,
			versionTag:  version,
			envoyAPITag: string(ec.GetEnvoyAPIVersion()),
		},
	})

	if err != nil {
		return newCacheError(UnknownError, "ensureEnvoyConfigRevision", err.Error())
	}
	err = r.Client.List(ctx, ecrList, &client.ListOptions{LabelSelector: selector})
	if err != nil {
		if err != nil {
			return newCacheError(UnknownError, "ensureEnvoyConfigRevision", err.Error())
		}
	}

	// Got wrong number of revisions
	if len(ecrList.Items) > 1 {
		return newCacheError(UnknownError, "ensureEnvoyConfigRevision", fmt.Sprintf("more than one revision exists for config version '%s', cannot reconcile", version))
	}

	// Revision does not yet exists, create one
	if len(ecrList.Items) == 0 {
		// Create the revision for this config version
		ecr := &envoyv1alpha1.EnvoyConfigRevision{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s", ec.Spec.NodeID, version),
				Namespace: ec.ObjectMeta.Namespace,
				Labels: map[string]string{
					nodeIDTag:   ec.Spec.NodeID,
					versionTag:  version,
					envoyAPITag: string(ec.GetEnvoyAPIVersion()),
				},
			},
			Spec: envoyv1alpha1.EnvoyConfigRevisionSpec{
				NodeID:         ec.Spec.NodeID,
				EnvoyAPI:       pointer.StringPtr(string(ec.GetEnvoyAPIVersion())),
				Version:        version,
				Serialization:  ec.Spec.Serialization,
				EnvoyResources: ec.Spec.EnvoyResources,
			},
		}
		// Set the ec as the owner and controller of the revision
		if err := controllerutil.SetControllerReference(ec, ecr, r.Scheme); err != nil {
			return newCacheError(UnknownError, "ensureEnvoyConfigRevision", err.Error())
		}
		err = r.Client.Create(ctx, ecr)
		if err != nil {
			return newCacheError(UnknownError, "ensureEnvoyConfigRevision", err.Error())
		}
	}

	return nil
}

func (r *EnvoyConfigReconciler) reconcileRevisionList(ctx context.Context, ec *envoyv1alpha1.EnvoyConfig) error {

	// Get all revisions owned by this EnvoyConfig that match the envoy API version
	ecrList := &envoyv1alpha1.EnvoyConfigRevisionList{}
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: map[string]string{
			nodeIDTag:   ec.Spec.NodeID,
			envoyAPITag: string(ec.GetEnvoyAPIVersion()),
		},
	})
	if err != nil {
		return newCacheError(UnknownError, "consolidateRevisionList", err.Error())
	}
	err = r.Client.List(ctx, ecrList, &client.ListOptions{LabelSelector: selector})
	if err != nil {
		if err != nil {
			return newCacheError(UnknownError, "consolidateRevisionList", err.Error())
		}
	}

	// Sort the revisions:
	// if publication timestamp is defined, by publication timestamp
	// if publication timestamp is not defined, by creation timestamp
	sort.SliceStable(ecrList.Items, func(i, j int) bool {
		var iTime, jTime metav1.Time
		if ecrList.Items[i].Status.LastPublishedAt.IsZero() {
			iTime = ecrList.Items[i].GetCreationTimestamp()
		} else {
			iTime = ecrList.Items[i].Status.LastPublishedAt
		}

		if ecrList.Items[j].Status.LastPublishedAt.IsZero() {
			jTime = ecrList.Items[j].GetCreationTimestamp()
		} else {
			jTime = ecrList.Items[j].Status.LastPublishedAt
		}

		return iTime.Before(&jTime)
	})

	// Generate the list using the previous order
	revisionList := make([]envoyv1alpha1.ConfigRevisionRef, len(ecrList.Items))
	for idx, ecr := range ecrList.Items {
		revisionList[idx] = envoyv1alpha1.ConfigRevisionRef{
			Version: ecr.Spec.Version,
			Ref: corev1.ObjectReference{
				Kind:       ecr.Kind,
				Name:       ecr.ObjectMeta.Name,
				Namespace:  ecr.Namespace,
				UID:        ecr.UID,
				APIVersion: ecr.APIVersion,
			},
		}
	}

	// Update the revision list in the EC status
	patch := client.MergeFrom(ec.DeepCopy())
	ec.Status.ConfigRevisions = revisionList

	// Remove older revisions if max have been reached
	ec.Status.ConfigRevisions = trimRevisions(ec.Status.ConfigRevisions, maxRevisions)

	err = r.Client.Status().Patch(ctx, ec, patch)
	if err != nil {
		return newCacheError(UnknownError, "consolidateRevisionList", err.Error())
	}

	return nil
}

func (r *EnvoyConfigReconciler) deleteUnreferencedRevisions(ctx context.Context, ec *envoyv1alpha1.EnvoyConfig) error {
	// Get all revisions that belong to this ec
	ecrList := &envoyv1alpha1.EnvoyConfigRevisionList{}
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: map[string]string{
			nodeIDTag:   ec.Spec.NodeID,
			envoyAPITag: string(ec.GetEnvoyAPIVersion()),
		}})
	if err != nil {
		return newCacheError(UnknownError, "deleteUnreferencedRevisions", err.Error())
	}
	err = r.Client.List(ctx, ecrList, &client.ListOptions{LabelSelector: selector})
	if err != nil {
		return newCacheError(UnknownError, "deleteUnreferencedRevisions", err.Error())
	}

	// For each of the revisions, check if they are still referred from the ec
	for _, ecr := range ecrList.Items {
		if getRevisionIndex(ecr.Spec.Version, ec.Status.ConfigRevisions) == nil {
			// Keep going even if the deletion operation returns error, we really don't care,
			// the ecr will eventually get deleted in a future reconcile loop
			_ = r.Client.Delete(ctx, &ecr)
		}
	}

	return nil
}

// markRevisionPublished marks the revision that matches the provided version as the one
// to be set in the xds server cache:
//  - It will first set the 'RevisionPublished' condition to false in the current published revision
//  - It will set the 'RevisionPublished' condition to true in the revision that matches the given version
// This ensures that at a given point in time 0 or 1 revisions can have the 'PublishedRevision' to true, being
// 1 the case most of the time
func (r *EnvoyConfigReconciler) markRevisionPublished(ctx context.Context, nodeID, version, reason, msg string, envoyAPI envoy.APIVersion) error {

	// Get all revisions for this EnvoyConfig
	ecrList := &envoyv1alpha1.EnvoyConfigRevisionList{}
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: map[string]string{
			nodeIDTag:   nodeID,
			envoyAPITag: string(envoyAPI),
		},
	})
	if err != nil {
		return newCacheError(UnknownError, "markRevisionPublished", err.Error())
	}
	err = r.Client.List(ctx, ecrList, &client.ListOptions{LabelSelector: selector})
	if err != nil {
		return newCacheError(UnknownError, "markRevisionPublished", err.Error())
	}

	// Set 'RevisionPublished' to false for all revisions
	for _, ecr := range ecrList.Items {
		if ecr.Spec.Version != version && ecr.Status.Conditions.IsTrueFor(envoyv1alpha1.RevisionPublishedCondition) {
			patch := client.MergeFrom(ecr.DeepCopy())
			ecr.Status.Conditions.SetCondition(status.Condition{
				Type:    envoyv1alpha1.RevisionPublishedCondition,
				Status:  corev1.ConditionFalse,
				Reason:  status.ConditionReason("OtherVersionPublished"),
				Message: msg,
			})

			if err := r.Client.Status().Patch(ctx, &ecr, patch); err != nil {
				return newCacheError(UnknownError, "markRevisionPublished", err.Error())
			}
		}
	}

	// NOTE: from this point on, if something fails we end up with 0 revisions
	// marked as published. Shouldn't be a problem as the current version
	// is already being served by the xds server and should be fixed eventually
	// in another reconcile

	// Set the the revision that holds the given version with 'RevisionPublished' = True
	ecrList = &envoyv1alpha1.EnvoyConfigRevisionList{}
	selector, err = metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: map[string]string{
			nodeIDTag:   nodeID,
			versionTag:  version,
			envoyAPITag: string(envoyAPI),
		},
	})

	if err != nil {
		return newCacheError(UnknownError, "markRevisionPublished", err.Error())
	}

	err = r.Client.List(ctx, ecrList, &client.ListOptions{LabelSelector: selector})
	if err != nil {
		return newCacheError(UnknownError, "markRevisionPublished", err.Error())
	}

	if len(ecrList.Items) != 1 {
		return newCacheError(UnknownError, "markRevisionPublished", fmt.Sprintf("found unexpected number of envoyconfigrevisions matching version '%s'", version))
	}

	ecr := ecrList.Items[0]
	patch := client.MergeFrom(ecr.DeepCopy())
	ecr.Status.Conditions.SetCondition(status.Condition{
		Type:    envoyv1alpha1.RevisionPublishedCondition,
		Status:  corev1.ConditionTrue,
		Reason:  status.ConditionReason(reason),
		Message: msg,
	})

	if err := r.Client.Status().Patch(ctx, &ecr, patch); err != nil {
		return newCacheError(UnknownError, "markRevisionPublished", err.Error())
	}

	return nil
}

func trimRevisions(list []envoyv1alpha1.ConfigRevisionRef, max int) []envoyv1alpha1.ConfigRevisionRef {
	for len(list) > max {
		list = list[1:]
	}
	return list
}

func calculateRevisionHash(resources *envoyv1alpha1.EnvoyResources) string {
	resourcesHasher := fnv.New32a()
	common.DeepHashObject(resourcesHasher, resources)
	return rand.SafeEncodeString(fmt.Sprint(resourcesHasher.Sum32()))
}

func getRevisionIndex(version string, revisions []envoyv1alpha1.ConfigRevisionRef) *int {
	for idx, rev := range revisions {
		if rev.Version == version {
			return &idx
		}
	}
	return nil
}

func moveRevisionToLast(list []envoyv1alpha1.ConfigRevisionRef, idx int) []envoyv1alpha1.ConfigRevisionRef {

	return append(list[:idx], append(list[idx+1:], list[idx])...)
}
