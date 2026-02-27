package awsmeta

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// GetMetaData ... fetch AWS meta-data.
func GetMetaData(path string) (contents []byte, err error) {
	url := "http://169.254.169.254/latest/meta-data/" + path

	req, _ := http.NewRequest("GET", url, nil)
	if metadataToken, _ := getMetaDataToken(); metadataToken != "" {
		req.Header.Set("X-aws-ec2-metadata-token", metadataToken)
	}

	client := http.Client{
		Timeout: time.Millisecond * 100,
	}

	resp, err := client.Do(req)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("awsmeta: code %d returned for url %s", resp.StatusCode, url)
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return
	}

	return []byte(body), err
}

// GetRegion ... get the effective EC2 region.
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

// getMetaDataToken ... get a metadata token for IMDSv2
func getMetaDataToken() (token string, err error) {
	url := "http://169.254.169.254/latest/api/token"

	req, _ := http.NewRequest("PUT", url, nil)
	req.Header.Set("X-aws-ec2-metadata-token-ttl-seconds", "5")

	client := http.Client{
		Timeout: time.Millisecond * 100,
	}

	resp, err := client.Do(req)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("awsmeta: code %d returned for url %s", resp.StatusCode, url)
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return
	}

	return string(body), err
}
