package e2e

import (
	"os"
	"testing"

	"sigs.k8s.io/e2e-framework/klient/conf"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
	"sigs.k8s.io/e2e-framework/support/kind"
)

var testenv env.Environment
var namespace = "default"

func TestMain(m *testing.M) {
	testenv = env.New()
	if os.Getenv("REAL_CLUSTER") == "true" {
		path := conf.ResolveKubeConfigFile()
		cfg := envconf.NewWithKubeConfig(path)
		testenv = env.NewWithConfig(cfg)
	} else {
		kindClusterName := envconf.RandomName("kind-with-config", 16)

		testenv.Setup(
			envfuncs.CreateCluster(kind.NewProvider(), kindClusterName),
			envfuncs.CreateNamespace(namespace),
		)

		testenv.Finish(
			envfuncs.DeleteNamespace(namespace),
			envfuncs.DestroyCluster(kindClusterName),
		)
	}

	os.Exit(testenv.Run(m))
}
