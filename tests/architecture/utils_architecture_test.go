package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestUtilsHaveNoDependencies(t *testing.T) {
	archtest.Package(t, utilsPackage).ShouldNotDependOn(
		configPackage,
		loggerPackage,
		domainLayer,
		packagesLayer,
		applicationLayer,
		persistenceLayer,
		presentationLayer,
	)
}

func TestUtilsHaveTests(t *testing.T) {
	archtest.Package(t, utilsPackage).IncludeTests()
}
