package query

import "github.com/AlterionX/ip-info-dump/infosource/base"

type QueryInfo struct {}

func (QueryInfo) FetchInfo(query base.Query) <-chan base.InfoResult {
	info := make(chan base.InfoResult)
	go func() {
		info <- base.InfoResult {
			Info: map[string]string {
				"resolved_ip": query.IP.String(),
				"resolved_address": query.Address,
			},
			Err: nil,
		}
	}()
	return info
}

func (QueryInfo) Name() string {
	return "query"
}
