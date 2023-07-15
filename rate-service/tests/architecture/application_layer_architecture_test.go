package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestApplicationLayerArchitecture(t *testing.T) {
	archtest.Package(t, applicationLayer).ShouldNotDependOn(
		packages,
		persistenceLayer,
		presentationLayer,
	)
}

func TestApplicationLayerHaveTests(t *testing.T) {
	archtest.Package(t, applicationLayer).IncludeTests()
}
