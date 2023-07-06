package architecture

import (
	"github.com/matthewmcnew/archtest"
	"testing"
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
