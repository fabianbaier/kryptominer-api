package core

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type EthermineClient struct {
	client http.Client
	url    string
}

type EthermineCurrentStatsResponse struct {
	Status string    `json:"status"`
	Data   etherData `json:"data"`
}

type etherData struct {
	AverageHashrate float64 `json:"reportedHashrate"`
	ActiveWorkers   int64   `json:"activeWorkers"`
	LastSeen        int64   `json:"lastseen"`
	Time            int64   `json:"time"`
	Balance         int64   `json:"unpaid"`
	Worker          string  `json:"worker"`
}

func NewEthermineClient(wallet string) *EthermineClient {
	c := http.Client{
		Timeout: time.Second * 20, // Maximum of 10 secs
	}
	return &EthermineClient{
		client: c,
		url:    "https://api.ethermine.org/miner/" + wallet + "/currentStats",
	}
}

func (e *EthermineClient) FetchCurrentStats() (*EthermineCurrentStatsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, e.url, nil)
	if err != nil {
		logrus.Warningf("Request Error: %s", err)
		return nil, err
	}

	req.Header.Set("User-Agent", "kryptominer")

	res, err := e.client.Do(req)
	if err != nil {
		logrus.Warningf("Client Error: %s", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Warningf("Read Error: %s", err)
		return nil, err
	}

	ethermineStats := &EthermineCurrentStatsResponse{}
	jsonErr := json.Unmarshal(body, ethermineStats)
	if jsonErr != nil {
		logrus.Warningf("JSON Error: %s", jsonErr)
	}
	return ethermineStats, nil
}
