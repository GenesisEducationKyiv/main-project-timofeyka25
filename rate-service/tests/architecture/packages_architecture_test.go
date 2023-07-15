package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestPackagesHaveNoDependencies(t *testing.T) {
	archtest.Package(t, packages).ShouldNotDependOn(
		configPackage,
		loggerPackage,
		utilsPackage,
		domainLayer,
		applicationLayer,
		persistenceLayer,
		presentationLayer,
	)
}
