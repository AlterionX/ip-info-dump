package whois

import (
	"net"
	"testing"

	"github.com/AlterionX/ip-info-dump/infosource/base"

	parser "github.com/likexian/whois-parser"
	"github.com/stretchr/testify/assert"
)

// TODO Figure out how to simulate network failures. But this is good enough for now.
func Test_WhoIs_FetchInfo_real(t *testing.T) {
	source := WhoIs{}

	{
		// NOTE This is google's dns IP. Might need to change this if that ever changes.
		google_ip := net.IPv4(8, 8, 8, 8)

		result_channel := source.FetchInfo(base.Query {
			IP: google_ip,
			Address: "dns.google",
		})

		info, ok := <-result_channel
		assert.True(t, ok, "channel to be fine")
		assert.Nil(t, info.Err, "request for google's whois to have worked")

		assert.Equal(t, "dns.google", info.Info.(parser.WhoisInfo).Domain.Domain, "ip to be google's ip")
		// TODO Possibly analyze some more?
	}

	{
		stackoverflow_ip := net.IPv4(151, 101, 193, 69)

		result_channel := source.FetchInfo(base.Query {
			IP: stackoverflow_ip,
			Address: "stackoverflow.com",
		})

		info, ok := <-result_channel
		assert.True(t, ok, "channel to be fine")
		assert.Nil(t, info.Err, "request for google's whois to have worked")

		assert.Equal(t, "stackoverflow.com", info.Info.(parser.WhoisInfo).Domain.Domain, "ip to be google's ip")
		// TODO Possibly analyze some more?
	}
}

func Test_WhoIs_Name(t *testing.T) {
	source := WhoIs{}
	assert.Equal(t, source.Name(), "whois", "name to match json conventions")
}
