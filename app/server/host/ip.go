package host

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

var Ip string

func LoadIp() error {
	if os.Getenv("GOENV") == "development" {
		Ip = "localhost"
		return nil
	}

	if os.Getenv("IS_AWS_ECS") != "" {
		var err error
		Ip, err = getAwsIp()

		if err != nil {
			return fmt.Errorf("error getting AWS ECS IP: %v", err)
		}
	} else if os.Getenv("IP") != "" {
		Ip = os.Getenv("IP")
		return nil
	}

	return nil
}

type ecsMetadata struct {
	Networks []struct {
		IPv4Addresses []string `json:"IPv4Addresses"`
	} `json:"Networks"`
}

var awsIp string

func getAwsIp() (string, error) {
	const ecsMetadataURL = "http://169.254.170.2/v2/metadata"

	resp, err := http.Get(ecsMetadataURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var metadata ecsMetadata
	err = json.Unmarshal(body, &metadata)
	if err != nil {
		return "", err
	}

	if len(metadata.Networks) == 0 || len(metadata.Networks[0].IPv4Addresses) == 0 {
		return "", errors.New("no IP address found in ECS metadata")
	}

	awsIp = metadata.Networks[0].IPv4Addresses[0]

	return awsIp, nil
}