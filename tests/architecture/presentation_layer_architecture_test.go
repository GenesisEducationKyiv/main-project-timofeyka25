package architecture

import (
	"github.com/matthewmcnew/archtest"
	"testing"
)

func TestPresentationLayerArchitecture(t *testing.T) {
	archtest.Package(t, presentationLayer).ShouldNotDependDirectlyOn(
		packagesLayer,
		applicationLayer,
		persistenceLayer,
	)
}
