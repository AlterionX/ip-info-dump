package whois

import (
	"log"

	"github.com/AlterionX/ip-info-dump/infosource/base"

	api "github.com/likexian/whois"
	parser "github.com/likexian/whois-parser"
)

// Isolate methods we want to stub for testing.
var baseAPICall = api.Whois
var baseParserCall = parser.Parse

type WhoIs struct{}

func (source WhoIs) FetchInfo(query base.Query) <-chan base.InfoResult {
	info := make(chan base.InfoResult)

	go func() {
		raw, err := baseAPICall(query.Address)
		if err != nil {
			log.Printf("Attempt to fetch response from WhoIs api with argument %q failed due to %q.", query.Address, err.Error())
			info <- base.InfoResult{
				Info: nil,
				Err:  err,
			}
			return
		}

		parsed, err := baseParserCall(raw)
		if err != nil {
			log.Printf("Attempt to parse response %q from WhoIs api call on address %q failed due to %q.", raw, query.Address, err.Error())
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
