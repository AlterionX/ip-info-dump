package whois

import (
	"log"
	"net"

    "github.com/AlterionX/ip-info-dump/infosource/base"

	api "github.com/likexian/whois"
	parser "github.com/likexian/whois-parser"
)

// Isolate methods we want to stub for testing.
var baseAPICall = api.Whois
var baseParserCall = parser.Parse

type WhoIs struct {}

func (source WhoIs) FetchInfo(arg net.IP) <-chan base.InfoResult {
    info := make(chan base.InfoResult)
    go func () {
        raw, err := baseAPICall(string(arg))
        if err != nil {
            log.Printf("Attempt to fetch response from WhoIs api with argument %q failed due to %q.", string(arg), err.Error())
            info <- base.InfoResult {
                Info: nil,
                Err: err,
            }
        }

        parsed, err := baseParserCall(raw)
        if err != nil {
            log.Printf("Attempt to parse response %q from WhoIs api failed due to %q.", raw, err.Error())
            info <- base.InfoResult {
                Info: nil,
                Err: err,
            }
        }

        info <- base.InfoResult {
            Info: parsed,
            Err: err,
        }
    }()
    return info
}

func (source WhoIs) Name() string {
    return "whois"
}

