package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestPresentationLayerArchitecture(t *testing.T) {
	archtest.Package(t, presentationLayer).ShouldNotDependDirectlyOn(
		packagesLayer,
		applicationLayer,
		persistenceLayer,
	)
}
