package whois

import (
	"net"
	"testing"

	parser "github.com/likexian/whois-parser"

	"github.com/stretchr/testify/assert"
)

// TODO Figure out how to simulate network failures. But this is good enough for now.
func Test_WhoIs_FetchInfo_real(t *testing.T) {
	t.Skip("avoid the network for now")
	source := WhoIs{}
	// NOTE This is google's IP. Might need to change this if google's IP ever changes.
	google_ip := net.IPv4(142, 250, 217, 78)
	result_channel := source.FetchInfo(google_ip)

	info, ok := <- result_channel
	assert.True(t, ok, "channel to be fine")
	assert.Nil(t, info.Err, "request for google's whois to have worked")

	assert.Equal(t, "google.com", info.Info.(parser.WhoisInfo).Domain.Domain, "ip to be google's ip")
	// TODO Possibly analyze some more?
}

func Test_WhoIs_Name(t *testing.T) {
	source := WhoIs{}
	assert.Equal(t, source.Name(), "whois", "name to match json conventions")
}

