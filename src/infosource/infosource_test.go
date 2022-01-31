package infosource

import (
	"testing"

	"github.com/AlterionX/ip-info-dump/infosource/base"
	"github.com/AlterionX/ip-info-dump/infosource/mock"
	"github.com/stretchr/testify/assert"
)

// Slightly integration.
func Test_GetAllSources(t *testing.T) {
	sources := GetAllSources()

	assert.Equal(t, len(sources), 2, "there should be three supported sources.")

	// This has no failure paths.
}

// This could replace a few tests, but the others are there for granularity.
func Test_GetInfo_mock(t *testing.T) {
	earlyDeath := "early_death"
	onlyReturn := "only_return"
	failure := "failure"
	returnedData := "survivor"
	sources := []base.InfoSource{
		mock.MockSource{
			EarlyExit: true,
			GivenName: &earlyDeath,
		},
		mock.MockSource{
			Info:      &returnedData,
			GivenName: &onlyReturn,
		},
		mock.MockSource{
			Err:       base.SourceFailure,
			GivenName: &failure,
		},
	}

	info, err := GetInfo("example.com", sources)
	assert.Nil(t, err, "invalid sources to silently fail and domain to be correct")
	assert.Len(t, info, 1, "only one surviving source")
	assert.Contains(t, info, onlyReturn, "the correct source to have survived")
	// TODO Check example.com's IP and stuff.

}

func Test_GetInfo_real(t *testing.T) {
	// TODO Actually run GetInfo against google.com
	t.Skip("Avoid running with real data unless explicitly specified.")
}
