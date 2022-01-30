package whois

import (
	"log"
	"net"

	"github.com/AlterionX/ip-info-dump/infosource/base"

	api "github.com/likexian/whois"
	parser "github.com/likexian/whois-parser"
)

// Isolate methods we want to stub for testing.
var reverseAddrLookup = net.LookupAddr
var baseAPICall = api.Whois
var baseParserCall = parser.Parse

type WhoIs struct{}

func (source WhoIs) FetchInfo(arg net.IP) <-chan base.InfoResult {
	info := make(chan base.InfoResult)

	go func() {
		addresses, err := reverseAddrLookup(arg.String())
		if err != nil || len(addresses) == 0 {
			log.Printf("Attempt to fetch address of IP %q failed due to %q.", arg.String(), err.Error())
			info <- base.InfoResult{
				Info: nil,
				Err:  err,
			}
			return
		}
		addr := addresses[0]

		raw, err := baseAPICall(addr)
		if err != nil {
			log.Printf("Attempt to fetch response from WhoIs api with argument %q failed due to %q.", addr, err.Error())
			info <- base.InfoResult{
				Info: nil,
				Err:  err,
			}
			return
		}

		parsed, err := baseParserCall(raw)
		if err != nil {
			log.Printf("Attempt to parse response %q from WhoIs api call on address %q failed due to %q.", raw, addr, err.Error())
			info <- base.InfoResult{
				Info: nil,
				Err:  err,
			}
			return
		}

		info <- base.InfoResult{
			Info: parsed,
			Err:  err,
		}
	}()
	return info
}

func (source WhoIs) Name() string {
	return "whois"
}
