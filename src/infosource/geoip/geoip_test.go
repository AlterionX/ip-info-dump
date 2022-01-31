package geoip

import (
	"net"
	"testing"

	"github.com/AlterionX/ip-info-dump/infosource/base"
	"github.com/stretchr/testify/assert"
)

// NOTE This only really tests the happy path. Perhaps we should test the more exotic failures as well?
func Test_GeoIP_FetchInfo_real(t *testing.T) {
	source := GeoIP {}

	ips, err := net.LookupIP("example.com")
	assert.NotEmpty(t, ips, "example.com should always exist")
	assert.Nil(t, err, "example.com should always exist")
	ip := ips[0]

	result_channel := source.FetchInfo(base.Query {
		IP: ip,
		Address: "example.com",
	})

	result, ok := <- result_channel

	assert.True(t, ok, "channel to be fine")

	info := result.Info.(GeoData)
	err = result.Err

	assert.Nil(t, err, "no errors")
	// Can't really say exactly what these values are... but they shouldn't be empty.
	assert.NotEmpty(t, info.Country, "country to not be empty")
	assert.NotEmpty(t, info.City, "country to not be empty")
}

func Test_GeoIP_Name(t *testing.T) {
	source := GeoIP {}
	assert.Equal(t, source.Name(), "geoip", "name to match json conventions")
}
