package infosource

import (
    "log"
    "net"
    "errors"
)

type InfoResult struct {
    info map[string]interface{}
    err error
}

type InfoSource interface {
    IsCompatible(arg string) bool
    FetchInfo(arg net.IP) <-chan InfoResult
    Name() string
}

func GetAllSources() []InfoSource {
    //return []InfoSource{whois {}, geoip {}, virustotal {}}
    return []InfoSource{MockSource {}}
}

var BadArgument = errors.New("Requested IP/domain is malformed.")

// Gets the relevant info for the argument after attempting to process.
//
// Returns either the data gathered or one of the known errors in this package.
func GetInfo(arg string, sources []InfoSource) (map[string]interface{}, error) {
    ip := net.ParseIP(arg)
    if ip == nil {
        // TODO Further research should be done to determine which ip address to use of if all of them should be used.
        ips, err := net.LookupIP(arg)
        if err != nil {
            log.Printf("Looking for ip for %q failed.", arg);
            return nil, BadArgument
        }
        ip = ips[0]
    }

    outputs := make(map[string](<-chan InfoResult))
    for _, source := range sources {
        outputs[source.Name()] = source.FetchInfo(ip)
    }

    data := make(map[string]interface{})
    for name, output := range outputs {
        info, ok := <-output
        // Errors here do NOT terminate, since, despite the error, we may have other API calls that succeeded.
        // Not terminating is picked over terminating since this API is a demo, and it isn't critical that
        // its operation be successful.
        // Switching to early termination is simple, though, so if being correctly formed is more important,
        // it's pretty easy to switch.
        if !ok {
            log.Printf("Channel for source %q failed.", name);
            continue
        }
        if info.err != nil {
            log.Printf("Attempt to retrieve info from %q failed due to %q.", name, info.err.Error())
            continue
        }
        data[name] = info.info
    }

    return data, nil
}

