package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	//"github.com/jarcoal/httpmock"
	//. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"testing"
)

const givenReceiversJson = `{
  "receivers": [
    {
      "id": "DE-HH-GS-19210",
      "name": "Hamburg 19210 - DE-HH-GS-019210",
      "hosts": [
        {
          "hostname": "host1.tld",
          "port": "19210"
        },
        {
          "hostname": "host2.tld",
          "port": "1111"
        }
      ],
      "receiverType": {
        "name": "Stationen",
        "category": "STATION"
      }
     },

	  {
		  "id": "DE-HH-GS-SI",
		  "name": "Hamburg Sirene 1",
		  "receiverType": {
			"name": "Sirene",
			"category": "SIRENE"
		  }
       }
   ]
  }`

//var _ = BeforeSuite(func() {
//	// block all HTTP requests
//	httpmock.Activate()
//})
//
//var _ = BeforeEach(func() {
//	// remove any mocks
//	httpmock.Reset()
//})
//
//var _ = AfterSuite(func() {
//	httpmock.DeactivateAndReset()
//})

func Test_JsonIsParsed(t *testing.T) {
	println(givenReceiversJson)
	var stations StationsListDto
	err := json.Unmarshal([]byte(givenReceiversJson), &stations)

	assert.Nil(t, err)
	assert.Equal(t, len(stations.Receivers), 2)
	assert.Equal(t, stations.Receivers[0].ID, "DE-HH-GS-19210")
	assert.Equal(t, stations.Receivers[0].Name, "Hamburg 19210 - DE-HH-GS-019210")
	assert.Equal(t, stations.Receivers[0].ReceiverType.Name, "Stationen")
	assert.Equal(t, stations.Receivers[0].ReceiverType.Category, "STATION")

	assert.Equal(t, len(stations.Receivers[0].Hosts), 2)
	assert.Equal(t, stations.Receivers[0].Hosts[0].Hostname, "host1.tld")
	assert.Equal(t, stations.Receivers[0].Hosts[0].Port, "19210")
	assert.Equal(t, stations.Receivers[0].Hosts[1].Hostname, "host2.tld")
	assert.Equal(t, stations.Receivers[0].Hosts[1].Port, "1111")

}

func Test_FetchStationsList_FiltersStations(t *testing.T) {
	// given
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, givenReceiversJson)
	}))
	defer svr.Close()

	s := Server{StationsEndpoint: svr.URL}

	// when
	stations, err := s.fetchStationsList()

	// then
	assert.Nil(t, err)
	assert.Equal(t, len(stations.Receivers), 2)
	assert.Equal(t, stations.Receivers[0].ID, "DE-HH-GS-19210")
	assert.Equal(t, stations.Receivers[0].Name, "Hamburg 19210 - DE-HH-GS-019210")
	assert.Equal(t, stations.Receivers[0].ReceiverType.Name, "Stationen")
	assert.Equal(t, stations.Receivers[0].ReceiverType.Category, "STATION")
}
