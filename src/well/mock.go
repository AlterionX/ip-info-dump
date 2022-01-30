package infosource

import "net"

type MockSource struct {}

func (MockSource) IsCompatible(arg string) bool {
    return true
}

func (MockSource) FetchInfo(arg net.IP) <-chan InfoResult {
    info := make(chan InfoResult)
    go func () {
        ip_string, err := arg.MarshalText()
        if err != nil {
            info <- InfoResult {
                info: nil,
                err: err,
            }
        }
        info <- InfoResult {
            info: map[string]interface{}{
                "staticValue": "mock_info",
                "ip": string(ip_string),
                "rawIp:": arg,
            },
            err: nil,
        }
    }()
    return info
}

func (MockSource) Name() string {
    return "MockSource"
}

