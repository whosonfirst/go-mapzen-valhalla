package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"	
	"github.com/whosonfirst/go-whosonfirst-uri"
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

func NewPointFromWOFId(id int64) (*Point, error) {

	url, err := uri.Id2AbsPath("https://whosonfirst.mapzen.com/data/", id)

	if err != nil {
		return nil, err
	}

	rsp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()

	f, err := feature.LoadFeatureFromReader(rsp.Body)

	if err != nil {
		return nil, err
	}

	c, err := whosonfirst.Centroid(f)

	if err != nil {
		return nil, err
	}

	coord := c.Coord()

	pt := Point{
		Latitude:  coord.Y,
		Longitude: coord.X,
	}

	return &pt, nil
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

	api_key := flag.String("api-key", "mapzen-xxxxxx", "A valid Mapzen API key.")
	endpoint := flag.String("endpoint", "valhalla.mapzen.com", "A valid Valhalla API endpoint.")
	costing := flag.String("costing", "auto", "A valid Valhalla costing.")

	str_from := flag.String("from", "", "Starting latitude,longitude position.")
	str_to := flag.String("to", "", "Destination latitude,longitude position.")

	from_wofid := flag.Int64("from-wofid", 0, "Starting Who's On First ID.")
	to_wofid := flag.Int64("to-wofid", 0, "Destination Who's On First ID.")

	flag.Parse()

	var from *Point
	var to *Point

	if *from_wofid != 0 {

		f, err := NewPointFromWOFId(*from_wofid)

		if err != nil {
			log.Fatal(err)
		}

		from = f

	} else {

		f, err := NewPointFromString(*str_from)

		if err != nil {
			log.Fatal(err)
		}

		from = f
	}

	if *to_wofid != 0 {

		t, err := NewPointFromWOFId(*to_wofid)

		if err != nil {
			log.Fatal(err)
		}

		to = t

	} else {

		t, err := NewPointFromString(*str_to)

		if err != nil {
			log.Fatal(err)
		}

		to = t
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
