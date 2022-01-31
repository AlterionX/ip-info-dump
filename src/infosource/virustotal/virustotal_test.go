package virustotal

import (
	"net"
	"testing"

	"github.com/AlterionX/ip-info-dump/infosource/base"
	"github.com/stretchr/testify/assert"
)

func Test_VirusTotal_FetchInfo_real(t *testing.T) {
	source := VirusTotal {}

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

	info := result.Info.(map[string]interface{})
	err = result.Err

	assert.Nil(t, err, "no errors")
	assert.NotNil(t, info["domain"], "domain to exist")
	assert.NotNil(t, info["ip"], "ip to exist")
}

// TODO Figure out how to make FetchInfo fail for an ip address without stubbing or network manipulation.
func Test_VirusTotal_FetchInfo_real_badDomain(t *testing.T) {
	source := VirusTotal {}

	ips, err := net.LookupIP("example.com")
	assert.NotEmpty(t, ips, "example.com should always exist")
	assert.Nil(t, err, "example.com should always exist")
	ip := ips[0]

	result_channel := source.FetchInfo(base.Query {
		IP: ip,
		Address: "notadomain",
	})

	result, ok := <- result_channel

	assert.True(t, ok, "channel to be fine")

	info := result.Info
	err = result.Err

	assert.Nil(t, info, "no valid info return")
	assert.NotNil(t, err, "error to exist")
}

// TODO Add more tests to simulate a network failing.

func Test_VirusTotal_Name(t *testing.T) {
	source := VirusTotal {}
	assert.Equal(t, source.Name(), "virustotal", "name to match json conventions")
}

