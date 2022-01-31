package infosource

import (
	"net"
	"testing"

	"github.com/AlterionX/ip-info-dump/infosource/base"

	"github.com/stretchr/testify/assert"
)

func Test_fetchIP(t *testing.T) {
	{
		query, err := resolveQuery("::")
		assert.Nil(t, err, "no error is expected")
		assert.Equal(t, net.ParseIP("::"), query.IP, "ip is equivalent to manually parsing it")
	}

	{
		query, err := resolveQuery("example.com")
		assert.Nil(t, err, "no error is expected")
		ips, err := net.LookupIP("example.com")
		assert.Nil(t, err, "no error is expected")
		assert.Contains(t, ips, query.IP, "ip returned is in the set of ips from lookup")
		assert.Equal(t, "example.com", query.Address, "ip returned is in the set of ips from lookup")
	}

	{
		query, err := resolveQuery("thisisnotawebsite")
		assert.ErrorIs(t, err, base.BadArgument, "the argument is invalide")
		assert.Nil(t, query, "no response when argument is invalid")
	}
}

func Test_checkSourceOutputChannel(t *testing.T) {
	name := "testing_channel"
	testcases := map[string]struct {
		input        base.InfoResult
		earlyClose   bool
		expectedData interface{}
		expectedErr  error
	}{
		"earlyClose":   {earlyClose: true},
		"dataReturned": {input: base.InfoResult{}, earlyClose: false},
		"errReturend":  {input: base.InfoResult{Err: base.SourceFailure}, earlyClose: false},
	}
	for testcase, test := range testcases {
		t.Logf("Running test case %q", testcase)

		output := make(chan base.InfoResult)
		if test.earlyClose {
			close(output)
		} else {
			go func() {
				output <- test.input
			}()
		}

		result, err := checkSourceOutputChannel(name, output)
		if test.expectedData == nil {
			assert.Nil(t, result)
		} else {
			assert.Equal(t, test.expectedData, result, "function to return input")
		}
		if test.expectedErr == nil {
		} else {
			assert.ErrorIs(t, err, test.expectedErr, "function to return equivalent error")
		}
	}
}
