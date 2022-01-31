package virustotal

import (
	"log"
	"os"

	"github.com/AlterionX/ip-info-dump/infosource/base"
	vt "github.com/VirusTotal/vt-go"
)

type VirusTotal struct {}

// This is a variable so that it's easier to mock later.
// TODO Also, consider not failing if only one of domain or ip fails.
var baseAPICall = func(query base.Query) (interface{}, error) {
	apiKey := os.Getenv("IPDUMP_VT_KEY")
	client := vt.NewClient(apiKey)

	domain_endpoint := vt.URL("domains/%s", query.Address)
	ip_endpoint := vt.URL("ip_addresses/%s", query.IP.String())

	domain_object, domain_err := client.GetObject(domain_endpoint)
	if domain_err != nil {
		log.Printf("Client failed to get domain object for domain %q due to %q.", query.Address, domain_err.Error())
		return nil, domain_err
	}

	ip_object, ip_err := client.GetObject(ip_endpoint)
	if ip_err != nil {
		log.Printf("Client failed to get ip object for domain %q due to %q.", query.Address, ip_err.Error())
		return nil, ip_err
	}

	return map[string]interface{}{
		"domain": domain_object,
		"ip": ip_object,
	}, nil
}

func (VirusTotal) FetchInfo(query base.Query) <-chan base.InfoResult {
	info := make(chan base.InfoResult)
	go func() {
		result, err := baseAPICall(query)
		if err != nil {
			info <- base.InfoResult {
				Info: nil,
				Err: err,
			}
			return
		}

		info <- base.InfoResult {
			Info: result,
			Err: nil,
		}
	}()
	return info
}

func (VirusTotal) Name() string {
	return "virustotal"
}
