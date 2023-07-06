package architecture

import (
	"github.com/matthewmcnew/archtest"
	"testing"
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
