package infosource

import (
    "log"
    "net"

    "github.com/AlterionX/ip-info-dump/infosource/whois"
    "github.com/AlterionX/ip-info-dump/infosource/base"
)

func GetAllSources() []base.InfoSource {
    //return []InfoSource{whois {}, geoip {}, virustotal {}}
    return []base.InfoSource{whois.WhoIs {}}
}

func fetchIP(arg string) (net.IP, error) {
    // TODO Is passing in localhost here a potential security risk?
    ip := net.ParseIP(arg)
    if ip == nil {
        // TODO Further research should be done to determine which ip address to use of if all of them should be used.
        ips, err := net.LookupIP(arg)
        if err != nil {
            log.Printf("Looking for ip for %q failed.", arg);
            return nil, base.BadArgument
        }
        ip = ips[0]
    }
    return ip, nil
}

func checkSourceOutputChannel(name string, output <-chan base.InfoResult) (interface{}, error) {
    info, ok := <-output
    // Errors here do NOT terminate, since, despite the error, we may have other API calls that succeeded.
    // Not terminating is picked over terminating since this API is a demo, and it isn't critical that
    // its operation be successful.
    // Switching to early termination is simple, though, so if being correctly formed is more important,
    // it's pretty easy to switch.
    if !ok {
        log.Printf("Channel for source %q failed.", name);
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
    ip, err := fetchIP(arg)
    if err != nil {
        return nil, err
    }

    outputs := make(map[string](<-chan base.InfoResult))
    for _, source := range sources {
        outputs[source.Name()] = source.FetchInfo(ip)
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

