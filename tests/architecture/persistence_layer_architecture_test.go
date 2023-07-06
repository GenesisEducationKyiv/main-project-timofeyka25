package architecture

import (
	"github.com/matthewmcnew/archtest"
	"testing"
)

func TestPersistenceLayerArchitecture(t *testing.T) {
	archtest.Package(t, persistenceLayer).ShouldNotDependOn(
		presentationLayer,
	)
}

func TestPersistenceLayerHaveTests(t *testing.T) {
	archtest.Package(t, persistenceLayer).IncludeTests()
}
