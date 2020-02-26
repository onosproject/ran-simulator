// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package manager

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// RansimulatorDb is the name of the database in Influx
const RansimulatorDb = "ransimulator"
const writeResource = "write"
const queryResource = "query"

func checkInfluxDbAvailable(influxDbAddr string) bool {
	dbURL, err := url.ParseRequestURI(fmt.Sprintf("http://%s", influxDbAddr))
	if err != nil {
		log.Error("to parse Influx DB address %s. %s", influxDbAddr, err.Error())
		return false
	}
	dbURL.Path = queryResource

	data := url.Values{}
	data.Set("q", "SHOW DATABASES")

	req, err := http.NewRequest("POST", dbURL.String(), strings.NewReader(data.Encode()))
	if err != nil {
		log.Error("Unable to form Influx DB request with %s. %s", influxDbAddr, err.Error())
		return false
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "ran-simulator")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("Error on Influx DB request with %s. %s", influxDbAddr, err.Error())
		return false
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("Error on Influx DB result with %s. %s", influxDbAddr, err.Error())
		return false
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Errorf("Expecting status code %d. Got %d. %s", 200, res.StatusCode, body)
		return false
	}
	if strings.Contains(string(body), RansimulatorDb) {
		return true
	}
	return false
}

// WriteHoMetricToInfluxDb writes out to the Influx dB - can be called from a go routine
func WriteHoMetricToInfluxDb(influxDbAddr string, ueName string, latency int64, time int64) {
	dbURL, err := url.ParseRequestURI(fmt.Sprintf("http://%s", influxDbAddr))
	if err != nil {
		log.Error("to parse Influx DB address %s. %s", influxDbAddr, err.Error())
	}
	dbURL.Path = writeResource
	dbURL.RawQuery = "db=" + RansimulatorDb
	data := fmt.Sprintf("hometrics,ue=%s value=%d %d", ueName, latency, time)
	log.Infof("Writing %s to %s", dbURL.String(), data)
	req, err := http.NewRequest("POST", dbURL.String(), strings.NewReader(data))
	if err != nil {
		log.Error("Unable to form Influx DB request with %s. %s", influxDbAddr, err.Error())
	}

	req.Header.Add("User-Agent", "ran-simulator")
	req.Header.Add("Cache-Control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("Error on Influx DB request with %s. %s", influxDbAddr, err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("Error on Influx DB result with %s. %s", influxDbAddr, err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode != 204 {
		log.Errorf("Expecting status code %d. Got %d. %s", 204, res.StatusCode, body)
	}
}
