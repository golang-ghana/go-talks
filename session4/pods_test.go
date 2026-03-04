package e2e

import (
	"context"
	"log"
	"testing"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestNginxPods(t *testing.T) {
	EXPECTED_PODS := 2

	feat := features.New("Nginx Deployed Pods").
		WithLabel("type", "API").
		Assess("Count number of default namespace pods", func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			var pods v1.PodList
			if err := c.Client().Resources(namespace).List(ctx, &pods); err != nil {
				t.Error(err)
			}
			// t.Logf("Got pods %v in namespace", len(pods.Items))
			if len(pods.Items) != EXPECTED_PODS {
				t.Fatalf("Expected %v pods in %s namespace but got %v", EXPECTED_PODS, namespace, len(pods.Items))
			}
			return ctx
		}).Feature()

	testenv.Test(t, feat)
}

func TestNginxPodsRunning(t *testing.T) {
	runningPodsFeature := features.New("Nginx Pods running").
		Assess("Deployed pods are running", func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			var pods v1.PodList

			err := c.Client().Resources(namespace).List(ctx, &pods)
			if err != nil {
				t.Error(err)
			}

			for _, pod := range pods.Items {
				log.Printf("pod %s", pod.Name)
				if pod.Status.Phase != "Running" {
					t.Fatalf("Pod %s is not running", pod.Name)
				}
			}

			return ctx
		}).Feature()

	testenv.Test(t, runningPodsFeature)
}
