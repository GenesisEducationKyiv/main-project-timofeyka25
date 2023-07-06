package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestPackagesLayerArchitecture(t *testing.T) {
	archtest.Package(t, packagesLayer).ShouldNotDependOn(
		configPackage,
		loggerPackage,
		utilsPackage,
		domainLayer,
		applicationLayer,
		persistenceLayer,
		presentationLayer,
	)
}
