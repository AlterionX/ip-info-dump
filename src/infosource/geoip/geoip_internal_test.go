package geoip

import (
	"net"
	"testing"

	"github.com/AlterionX/ip-info-dump/infosource/base"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
)

func Test_GeoIP_FetchInfo_mock(t *testing.T) {
	good_ip := net.IPv4(1, 1, 1, 1)
	bad_ip := net.IPv4(0, 6, 6, 6)

	apiStub := gostub.Stub(&baseAPICall, func(ip net.IP) (GeoData, error) {
		if ip.Equal(good_ip) {
			return GeoData {
				Country: "good",
			}, nil
		}
		if ip.Equal(bad_ip) {
			return GeoData{}, base.BadArgument
		}
		return GeoData {
			Country: "gibberish",
		}, nil
	})
	defer apiStub.Reset()

	source := GeoIP {}

	{
		result_channel := source.FetchInfo(base.Query {
			IP: good_ip,
			Address: "something",
		})

		result, ok := <- result_channel
		assert.True(t, ok, "channel to be fine")
		info := result.Info.(GeoData)
		err := result.Err
		assert.Nil(t, err, "no errors for good ip")
		assert.Equal(t, GeoData{ Country: "good" }, info, "no transformations on return from baseAPICall")
	}
}

// TODO Create more tests for individual http failure points.
