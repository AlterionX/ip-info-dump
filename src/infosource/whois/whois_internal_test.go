package whois

import (
	"net"
	"testing"

    "github.com/AlterionX/ip-info-dump/infosource/base"

	"github.com/stretchr/testify/assert"
	"github.com/prashantv/gostub"
	parser "github.com/likexian/whois-parser"
)

func Test_WhoIs_FetchInfo_mock(t *testing.T) {
	source := WhoIs{}
	good_ip := net.IPv4(1, 1, 1, 1)
	ok_ip := net.IPv4(1, 6, 6, 6)
	bad_ip := net.IPv4(0, 6, 6, 6)

	apiStub := gostub.Stub(&baseAPICall, func(ip string, ss ...string) (string, error) {
		if ip == string(good_ip) {
			return "good", nil
		}
		if ip == string(ok_ip) {
			return "bad", nil
		}
		if ip == string(bad_ip) {
			return "", base.BadArgument
		}
		return "gibberish", nil
	})
	defer apiStub.Reset()

	parserStub := gostub.Stub(&baseParserCall, func(raw string) (parser.WhoisInfo, error) {
		if raw == "good" {
			return parser.WhoisInfo {
				Domain: &parser.Domain {
					Domain: "good",
				},
			}, nil
		}
		if raw == "bad" {
			return parser.WhoisInfo {}, base.SourceFailure
		}
		return parser.WhoisInfo {
			Domain: &parser.Domain {
				Domain: "gobblygook",
			},
		}, nil
	})
	defer parserStub.Reset()

	{
		good_channel := source.FetchInfo(good_ip)
		info, ok := <-good_channel
		assert.True(t, ok, "channel to be fine")
		assert.Nil(t, info.Err, "no valid info to be returned")
		assert.Equal(t, info.Info.(parser.WhoisInfo).Domain.Domain, "good", "the information to be forwarded")
	}
	{
		ok_channel := source.FetchInfo(ok_ip)
		info, ok := <-ok_channel
		assert.True(t, ok, "channel to be fine")
		assert.Nil(t, info.Info, "no valid info to be returned")
		assert.ErrorIs(t, info.Err, base.SourceFailure, "the correct error to be forwarded")
	}
	{
		bad_channel := source.FetchInfo(bad_ip)
		info, ok := <-bad_channel
		assert.True(t, ok, "channel to be fine")
		assert.Nil(t, info.Info, "no valid info to be returned")
		assert.ErrorIs(t, info.Err, base.BadArgument, "the correct error to be forwarded")
	}
}

