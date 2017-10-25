package core

import (
	"fmt"
	"time"

	"github.com/fabianbaier/kryptominer-api/config"
	"github.com/fabianbaier/kryptominer-api/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type updateETHWalletRequest struct {
	ActiveWorkers   int64   `json:"activeWorkers"`
	AverageHashrate float64 `json:"averageHashrate"`
	Balance         int64   `json:"balance"`
	Time            int64   `json:"time"`
}

type apiResponse struct {
	Status  string
	Message string `json:"message,omitempty"`
}

func Start(appConfig config.Config) {
	port := fmt.Sprintf(":%d", appConfig.Port)
	logrus.Infof("Starting Kryptominer API on 0.0.0.0%s", port)

	ec := NewEthermineClient(appConfig.EthWallet)
	db, err := storage.GetConnection(appConfig.DBURL)
	if err != nil {

	}

	// setup ticker for fetching Ethermine data
	go func() {
		ticker := time.NewTicker(time.Second * 91)
		for _ = range ticker.C {
			res, err := ec.FetchCurrentStats()
			if err != nil {
				logrus.Warningf("failed fetching current stats: %s", err)
			}
			w := storage.Wallet{
				ActiveWorkers:   res.Data.ActiveWorkers,
				AverageHashrate: res.Data.AverageHashrate,
				Time:            res.Data.Time,
				Balance:         res.Data.Balance,
				Address:         appConfig.EthWallet,
			}
			err = db.InsertWalletState(w)
			if err != nil {
				logrus.Warningf("failed saving wallet state to db: %s", err)
			}
		}
	}()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("store", db)
		c.Next()
	})

	r.GET("/api/v1/eth/wallets/:address", getETHWalletHandler)
	r.POST("/api/v1/eth/wallets/:address", updateETHWalletHandler)
	r.Run(port)
}

func getETHWalletHandler(c *gin.Context) {
	db, ok := c.Get("store")
	if !ok {
		return
	}
	ethAddress := c.Param("address")
	w, err := db.(*storage.Connection).GetWallet(ethAddress)
	if err != nil {
		message := apiResponse{
			Status:  "FAILED",
			Message: fmt.Sprintf("%s", err),
		}
		c.JSON(500, message)
		return
	} else {
		c.JSON(200, w)
	}
}

func updateETHWalletHandler(c *gin.Context) {
	var req updateETHWalletRequest
	err := c.BindJSON(&req)
	if err != nil {
		// bail out
		logrus.Warningf("Update Wallet Error: Parsing: %s", err)
		return
	}
	addr := c.Param("address")
	payload := storage.Wallet{
		ActiveWorkers:   req.ActiveWorkers,
		Address:         addr,
		AverageHashrate: req.AverageHashrate,
		Time:            req.Time,
		Balance:         req.Balance,
	}

	db, ok := c.Get("store")
	if !ok {
		return
	}

	err = db.(*storage.Connection).InsertWalletState(payload)
	if err != nil {
		message := apiResponse{
			Status:  "FAILED",
			Message: fmt.Sprintf("%s", err),
		}
		c.JSON(500, message)
		return
	}
	c.JSON(201, &apiResponse{
		Status: "OK",
	})
}
