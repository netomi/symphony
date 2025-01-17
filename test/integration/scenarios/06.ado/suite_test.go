package scenarios_test

import (
	"context"
	"testing"
	"time"

	_ "embed"

	"github.com/eclipse-symphony/symphony/packages/testutils/conditions"
	"github.com/eclipse-symphony/symphony/packages/testutils/expectations/kube"
	"github.com/eclipse-symphony/symphony/packages/testutils/logger"
	"github.com/eclipse-symphony/symphony/test/integration/lib/shell"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

//go:embed manifest/instance.yaml
var defaultInstanceManifest []byte

//go:embed manifest/target.yaml
var defaultTargetManifest []byte

//go:embed manifest/solution.yaml
var defaultSolutionManifest []byte

var successfullTargetExpectation = kube.Must(kube.Target("target", "default", kube.WithCondition(conditions.All(
	kube.ProvisioningSucceededCondition,
	//kube.OperationIdMatchCondition,
))))

var successfullInstanceExpectation = kube.Must(kube.Instance("instance", "default", kube.WithCondition(conditions.All(
	kube.ProvisioningSucceededCondition,
	//kube.OperationIdMatchCondition,
))))

var failedTargetExpectation = kube.Must(kube.Target("target", "default", kube.WithCondition(conditions.All(
	kube.ProvisioningFailedCondition,
	//kube.OperationIdMatchCondition,
))))

var failedInstanceExpectation = kube.Must(kube.Instance("instance", "default", kube.WithCondition(conditions.All(
	kube.ProvisioningFailedCondition,
	//kube.OperationIdMatchCondition,
))))

var absentInstanceExpectation = kube.Must(kube.AbsentInstance("instance", "default"))
var absentTargetExpectation = kube.Must(kube.AbsentTarget("target", "default"))

var _ = BeforeSuite(func(ctx context.Context) {
	// err := shell.LocalenvCmd(ctx, "mage cluster:load")
	// Expect(err).ToNot(HaveOccurred())

	shell.LocalenvCmd(ctx, "mage cluster:deploy")

	logger.SetDefaultLogger(GinkgoWriter.Printf)
}, NodeTimeout(5*time.Minute))

func TestScenarios(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scenarios Suite")
}
