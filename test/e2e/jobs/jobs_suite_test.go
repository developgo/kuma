package jobs_test

import (
	"testing"

	. "github.com/onsi/ginkgo"

	"github.com/kumahq/kuma/pkg/test"
	"github.com/kumahq/kuma/test/e2e/jobs"
	"github.com/kumahq/kuma/test/framework"
)

func TestE2EJobsKubernetes(t *testing.T) {
	if framework.IsK8sClustersStarted() {
		test.RunSpecs(t, "Jobs Kubernetes Suite")
	} else {
		t.SkipNow()
	}
}

var _ = Describe("Jobs on Kubernetes", jobs.Jobs)
