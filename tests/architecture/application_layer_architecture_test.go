package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestApplicationLayerArchitecture(t *testing.T) {
	archtest.Package(t, applicationLayer).ShouldNotDependOn(
		packagesLayer,
		persistenceLayer,
		presentationLayer,
	)
}

func TestApplicationLayerHaveTests(t *testing.T) {
	archtest.Package(t, applicationLayer).IncludeTests()
}
