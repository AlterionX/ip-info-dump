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
	horrible_ip := net.IPv4(0, 0, 66, 6)
	horrible_ip_string := "0.0.66.6"

	lookupStub := gostub.Stub(&reverseAddrLookup, func(ip string) ([]string, error) {
		if ip == good_ip_string {
			return []string{"good"}, nil
		}
		if ip == ok_ip_string {
			return []string{"ok"}, nil
		}
		if ip == bad_ip_string {
			return []string{"bad"}, nil
		}
		if ip == horrible_ip_string {
			return nil, base.BadArgument
		}
		return []string{"gibberish"}, nil
	})
	defer lookupStub.Reset()

	apiStub := gostub.Stub(&baseAPICall, func(addr string, ss ...string) (string, error) {
		if addr == "good" {
			return "good", nil
		}
		if addr == "ok" {
			return "bad", nil
		}
		if addr == "bad" {
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
		good_channel := source.FetchInfo(good_ip)
		info, ok := <-good_channel
		assert.True(t, ok, "channel to be fine")
		assert.Nil(t, info.Err, "no error to be returned")
		assert.Equal(t, "good", info.Info.(parser.WhoisInfo).Domain.Domain, "the information to be forwarded")
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
	{
		horrible_channel := source.FetchInfo(horrible_ip)
		info, ok := <-horrible_channel
		assert.True(t, ok, "channel to be fine")
		assert.Nil(t, info.Info, "no valid info to be returned")
		assert.ErrorIs(t, info.Err, base.BadArgument, "the correct error to be forwarded")
	}
}
