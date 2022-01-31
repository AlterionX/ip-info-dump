package whois

import (
	"net"
	"testing"

	"github.com/AlterionX/ip-info-dump/infosource/base"

	parser "github.com/likexian/whois-parser"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
)

func Test_WhoIs_FetchInfo_mock(t *testing.T) {
	source := WhoIs{}
	good_ip := net.IPv4(1, 1, 1, 1)
	good_ip_string := "1.1.1.1"
	ok_ip := net.IPv4(1, 6, 6, 6)
	ok_ip_string := "1.6.6.6"
	bad_ip := net.IPv4(0, 6, 6, 6)
	bad_ip_string := "0.6.6.6"

	apiStub := gostub.Stub(&baseAPICall, func(addr string, ss ...string) (string, error) {
		if addr == good_ip_string {
			return "good", nil
		}
		if addr == ok_ip_string {
			return "bad", nil
		}
		if addr == bad_ip_string {
			return "", base.BadArgument
		}
		return "gibberish", nil
	})
	defer apiStub.Reset()

	parserStub := gostub.Stub(&baseParserCall, func(raw string) (parser.WhoisInfo, error) {
		if raw == "good" {
			return parser.WhoisInfo{
				Domain: &parser.Domain{
					Domain: "good",
				},
			}, nil
		}
		if raw == "bad" {
			return parser.WhoisInfo{}, base.SourceFailure
		}
		return parser.WhoisInfo{
			Domain: &parser.Domain{
				Domain: "gobblygook",
			},
		}, nil
	})
	defer parserStub.Reset()

	{
		good_channel := source.FetchInfo(base.Query { IP: good_ip, Address: good_ip_string })
		info, ok := <-good_channel
		assert.True(t, ok, "channel to be fine")
		assert.Nil(t, info.Err, "no error to be returned")
		assert.Equal(t, "good", info.Info.(parser.WhoisInfo).Domain.Domain, "the information to be forwarded")
	}
	{
		ok_channel := source.FetchInfo(base.Query { IP: ok_ip, Address: ok_ip_string })
		info, ok := <-ok_channel
		assert.True(t, ok, "channel to be fine")
		assert.Nil(t, info.Info, "no valid info to be returned")
		assert.ErrorIs(t, info.Err, base.SourceFailure, "the correct error to be forwarded")
	}
	{
		bad_channel := source.FetchInfo(base.Query { IP: bad_ip, Address: bad_ip_string })
		info, ok := <-bad_channel
		assert.True(t, ok, "channel to be fine")
		assert.Nil(t, info.Info, "no valid info to be returned")
		assert.ErrorIs(t, info.Err, base.BadArgument, "the correct error to be forwarded")
	}
}
