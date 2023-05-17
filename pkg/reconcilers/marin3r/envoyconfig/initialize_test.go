package reconcilers

import (
	"testing"

	marin3rv1alpha1 "github.com/3scale-ops/marin3r/apis/marin3r/v1alpha1"
	"github.com/3scale-ops/marin3r/pkg/envoy"
	"github.com/3scale-ops/marin3r/pkg/util/pointer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIsInitialized(t *testing.T) {
	tests := []struct {
		name               string
		envoyConfigFactory func() *marin3rv1alpha1.EnvoyConfig
		want               bool
	}{
		{
			"Initializes the resource",
			func() *marin3rv1alpha1.EnvoyConfig {
				return &marin3rv1alpha1.EnvoyConfig{}
			},
			false,
		},
		{
			"Returns true if already initialized",
			func() *marin3rv1alpha1.EnvoyConfig {
				return &marin3rv1alpha1.EnvoyConfig{
					Spec: marin3rv1alpha1.EnvoyConfigSpec{
						EnvoyAPI: pointer.New(envoy.APIv3),
					},
				}
			},
			true,
		},
		{
			"Returns false if finalizer detected",
			func() *marin3rv1alpha1.EnvoyConfig {
				return &marin3rv1alpha1.EnvoyConfig{
					ObjectMeta: metav1.ObjectMeta{
						Finalizers: []string{marin3rv1alpha1.EnvoyConfigRevisionFinalizer},
					},
					Spec: marin3rv1alpha1.EnvoyConfigSpec{
						EnvoyAPI: pointer.New(envoy.APIv3),
					},
				}
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInitialized(tt.envoyConfigFactory()); got != tt.want {
				t.Errorf("IsInitialized() = %v, want %v", got, tt.want)
			}
		})
	}
}
