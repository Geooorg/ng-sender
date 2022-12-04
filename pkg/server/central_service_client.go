package server

import (
	"encoding/json"
	"io"
	"log"
	"ng-receiver/pkg/httpclient"
)

type StationsListDto struct {
	Receivers []struct {
		ID    string
		Name  string
		Hosts []struct {
			Hostname string
			Port     string
		}
		ReceiverType struct {
			Category string
			Name     string
		}
	}
}

func (s *Server) fetchStationsList() (StationsListDto, error) {

	res, err := httpclient.NewClient().Get(s.StationsEndpoint)

	if err != nil {
		log.Println("WARN: Could not fetch latest stations list")
		return StationsListDto{}, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, _ := io.ReadAll(res.Body)
	var stations StationsListDto
	err = json.Unmarshal(body, &stations)
	if err != nil {
		log.Println("WARN: Could not fetch latest stations list")
		return StationsListDto{}, err
	}

	return stations, nil
}
