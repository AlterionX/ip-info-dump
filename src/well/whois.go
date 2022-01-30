package infosource

import (
    "fmt"
    api "github.com/likexian/whois"
)

func basicWhois() {
    result, err := api.Whois("data")
    if err != nil {
        // Handle Error
        return
    }
    fmt.Printf("Hello %q", result)
}

type whois struct {}

