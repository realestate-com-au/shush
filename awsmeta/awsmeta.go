package awsmeta

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// GetMetaData ... fetch AWS meta-data.
//
func GetMetaData(path string) (contents []byte, err error) {
	url := "http://169.254.169.254/latest/meta-data/" + path

	req, _ := http.NewRequest("GET", url, nil)
	client := http.Client{
		Timeout: time.Millisecond * 100,
	}

	resp, err := client.Do(req)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("Code %d returned for url %s", resp.StatusCode, url)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	return []byte(body), err
}

// GetRegion ... get the effective EC2 region.
//
func GetRegion() string {
	path := "placement/availability-zone"

	resp, err := GetMetaData(path)
	if err != nil {
		return ""
	}

	az := string(resp)
	if len(az) < 1 {
		return ""
	}

	//returns us-west-2a, just return us-west-2
	return string(az[:len(az)-1])
}
