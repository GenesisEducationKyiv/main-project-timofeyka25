package architecture

import (
	"github.com/matthewmcnew/archtest"
	"testing"
)

func TestDomainLayerArchitecture(t *testing.T) {
	archtest.Package(t, domainLayer).ShouldNotDependOn(
		configPackage,
		loggerPackage,
		utilsPackage,
		packagesLayer,
		applicationLayer,
		persistenceLayer,
		presentationLayer,
	)
}
