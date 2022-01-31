package geoip

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"encoding/json"

	"github.com/AlterionX/ip-info-dump/infosource/base"
)

type GeoData struct {
	// ISO is okay, since it's relatively well know.
	Continent string `json:"continent"`
	Country string `json:"country"`
	Region string `json:"region"`
	City string `json:"city"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

var baseAPICall = func(ip net.IP) (GeoData, error) {
	endpoint := fmt.Sprintf("http://ip-api.com/json/%s", ip.String());
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		log.Printf("Failed to fetch data (create request) for geoip query %q due to %q.", ip.String(), err.Error())
		return GeoData {}, err
	}

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Failed to fetch data (response) for geoip query for %q due to %q.", ip.String(), err.Error())
		return GeoData {}, err
	}
	if response.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch data for geoip query for %q failed with response code %d.", ip.String(), response.StatusCode)
		return GeoData {}, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Failed to load data for response to geoip query for %q due to %q.", ip.String(), err.Error())
		return GeoData {}, err
	}

	info := GeoData {}
	err = json.Unmarshal(data, &info)
	if err != nil {
		log.Printf("Failed to parse data %q for response to geoip query for %q due to %q.", data, ip.String(), err.Error())
		return GeoData {}, err
	}

	return GeoData {
	}, nil
}

func requestGeoIP(ip net.IP) {
}

type GeoIP struct {}

func (source GeoIP) FetchInfo(query base.Query) <-chan base.InfoResult {
	info := make(chan base.InfoResult)

	go func() {
		geoData, err := baseAPICall(query.IP)
		if err != nil {
			log.Printf("Attempt to fetch response from GeoIP api with argument %q failed due to %q.", query.IP, err.Error())
			info <- base.InfoResult{
				Info: nil,
				Err:  err,
			}
			return
		}

		info <- base.InfoResult{
			Info: geoData,
			Err:  err,
		}
	}()
	return info
}

func (source GeoIP) Name() string {
	return "geoip"
}
