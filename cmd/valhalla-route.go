package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Point struct {
	Latitude  float64
	Longitude float64
}

func NewPointFromString(str_latlon string) (*Point, error) {

	latlon := strings.Split(str_latlon, ",")

	if len(latlon) != 2 {
		return nil, errors.New("Invalid lat,lon string")
	}

	lat, err := strconv.ParseFloat(latlon[0], 64)

	if err != nil {
		return nil, err
	}

	lon, err := strconv.ParseFloat(latlon[1], 64)

	if err != nil {
		return nil, err
	}

	pt := Point{
		Latitude:  lat,
		Longitude: lon,
	}

	return &pt, nil
}

type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Type      string  `json:"type"`
}

type Journey struct {
	Locations []Location `json:"locations"`
	Costing   string     `json:"costing"`
}

type Valhalla struct {
	Endpoint string
	ApiKey   string
}

func (v *Valhalla) Route(from *Point, to *Point, costing string) ([]byte, error) {

	loc_from := Location{
		Latitude:  from.Latitude,
		Longitude: from.Longitude,
		Type:      "break",
	}

	loc_to := Location{
		Latitude:  to.Latitude,
		Longitude: to.Longitude,
		Type:      "break",
	}

	journey := Journey{
		Locations: []Location{loc_from, loc_to},
		Costing:   costing,
	}

	body, err := json.Marshal(journey)

	if err != nil {
	   return nil, err
	}
	
	query := url.Values{}
	query.Set("json", string(body))
	query.Set("api_key", v.ApiKey)

	u := url.URL{
		RawQuery: query.Encode(),
		Host:     v.Endpoint,
		Path:     "/route",
		Scheme:   "https",
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", u.String(), nil)

	if err != nil {
		return nil, err
	}

	rsp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		return nil, errors.New(rsp.Status)
	}

	body, err = ioutil.ReadAll(rsp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func main() {

	api_key := flag.String("api-key", "mapzen-xxxxxx", "...")
	endpoint := flag.String("endpoint", "valhalla.mapzen.com", "...")
	costing := flag.String("costing", "auto", "...")

	str_from := flag.String("from", "", "...")
	str_to := flag.String("to", "", "...")

	flag.Parse()

	from, err := NewPointFromString(*str_from)

	if err != nil {
		log.Fatal(err)
	}

	to, err := NewPointFromString(*str_to)

	if err != nil {
		log.Fatal(err)
	}

	v := Valhalla{
		Endpoint: *endpoint,
		ApiKey:   *api_key,
	}

	b, err := v.Route(from, to, *costing)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
