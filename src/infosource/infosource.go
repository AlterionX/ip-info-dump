package infosource

import (
	"log"
	"net"

	"github.com/AlterionX/ip-info-dump/infosource/base"
	"github.com/AlterionX/ip-info-dump/infosource/query"
	"github.com/AlterionX/ip-info-dump/infosource/whois"
)

func GetAllSources() []base.InfoSource {
	//return []InfoSource{whois {}, geoip {}, virustotal {}}
	return []base.InfoSource{
		query.QueryInfo {},
		whois.WhoIs{},
	}
}

func resolveQuery(arg string) (*base.Query, error) {
	// TODO Is passing in localhost here a potential security risk?
	var addr string
	ip := net.ParseIP(arg)
	if ip == nil {
		addr = arg
		// TODO Further research should be done to determine which ip address to use of if all of them should be used.
		ips, err := net.LookupIP(arg)
		if err != nil {
			log.Printf("Looking for ip for %q failed.", arg)
			return nil, base.BadArgument
		}
		log.Printf("Selecting first ip out of %q.", ips)
		ip = ips[0]
	} else {
		addresses, err := net.LookupAddr(ip.String())
		if err != nil || len(addresses) == 0 {
			log.Printf("Looking for address matching %q failed.", ip.String())
			return nil, base.BadArgument
		}
		log.Printf("Selecting first address out of %q.", addresses)
		addr = addresses[0]
	}

	query := base.Query {
		IP: ip,
		Address: addr,
	}
	return &query, nil
}

func checkSourceOutputChannel(name string, output <-chan base.InfoResult) (interface{}, error) {
	info, ok := <-output
	// Errors here do NOT terminate, since, despite the error, we may have other API calls that succeeded.
	// Not terminating is picked over terminating since this API is a demo, and it isn't critical that
	// its operation be successful.
	// Switching to early termination is simple, though, so if being correctly formed is more important,
	// it's pretty easy to switch.
	if !ok {
		log.Printf("Channel for source %q failed.", name)
		return nil, base.MissingReply
	}
	if info.Err != nil {
		log.Printf("Attempt to retrieve info from %q failed due to %q.", name, info.Err.Error())
		return nil, base.SourceFailure
	}

	return info.Info, nil
}

// Gets the relevant info for the argument after attempting to process.
//
// Returns either the data gathered or one of the known errors in this package.
func GetInfo(arg string, sources []base.InfoSource) (map[string]interface{}, error) {
	query_ptr, err := resolveQuery(arg)
	if err != nil {
		return nil, err
	}
	query := *query_ptr

	outputs := make(map[string](<-chan base.InfoResult))
	for _, source := range sources {
		outputs[source.Name()] = source.FetchInfo(query)
	}

	data := make(map[string]interface{})
	for name, output := range outputs {
		info, err := checkSourceOutputChannel(name, output)
		if err != nil {
			// Errors here do NOT terminate, since, despite the error, we may have other API calls that succeeded.
			// Not terminating is picked over terminating since this API is a demo, and it isn't critical that
			// its operation be successful.
			// Switching to early termination is simple, though, so if being correctly formed is more important,
			// it's pretty easy to switch.
			log.Printf("Attempting to retrieve message from channel of source %q failed due to %q. Ignoring...", name, err.Error())
			continue
		}

		data[name] = info
	}

	return data, nil
}
