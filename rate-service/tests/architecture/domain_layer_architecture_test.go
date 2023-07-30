package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestDomainLayerArchitecture(t *testing.T) {
	archtest.Package(t, domainLayer).ShouldNotDependOn(
		configPackage,
		loggerPackage,
		utilsPackage,
		packages,
		applicationLayer,
		persistenceLayer,
		presentationLayer,
	)
}
