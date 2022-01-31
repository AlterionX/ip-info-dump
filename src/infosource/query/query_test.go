package query

import (
	"net"
	"testing"

	"github.com/AlterionX/ip-info-dump/infosource/base"
	"github.com/stretchr/testify/assert"
)

func Test_QueryInfo_FetchInfo(t *testing.T) {
	source := QueryInfo {}

	result_channel := source.FetchInfo(base.Query {
		IP: net.IPv4(0, 0, 0, 0),
		Address: "localhost",
	})

	result, ok := <- result_channel

	assert.True(t, ok, "channel to be fine")

	info := result.Info.(map[string]string)
	err := result.Err

	assert.Nil(t, err, "no errors")
	assert.Equal(t, "0.0.0.0", info["resolved_ip"], "no errors")
	assert.Equal(t, "localhost", info["resolved_address"], "no errors")
}

func Test_QueryInfo_Name(t *testing.T) {
	source := QueryInfo {}
	assert.Equal(t, source.Name(), "query", "name to match json conventions")
}
