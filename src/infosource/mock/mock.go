package mock

// TODO There was an attempt to make this a test only package, but the go compiler kept complaining about MockSource being not found, so it's now a build package.

import "github.com/AlterionX/ip-info-dump/infosource/base"

type MockSource struct {
	EarlyExit bool
	Info      interface{}
	Err       error
	GivenName *string
}

func (source MockSource) FetchInfo(query base.Query) <-chan base.InfoResult {
	info := make(chan base.InfoResult)
	go func() {
		if source.EarlyExit {
			close(info)
			return
		}
		if source.Info != nil {
			info <- base.InfoResult{
				Info: map[string]interface{}{
					"static_value":  "mock_info",
					"raw_ip:":       query.IP,
					"expected_info": source.Info,
				},
				Err: nil,
			}
			return
		}
		if source.Err != nil {
			info <- base.InfoResult{
				Info: nil,
				Err:  source.Err,
			}
			return
		}
	}()
	return info
}

func (source MockSource) Name() string {
	if source.GivenName == nil {
		return "mock_source"
	} else {
		return *source.GivenName
	}
}
