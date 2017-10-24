package storage

import (
	"gopkg.in/mgo.v2"
	"sync"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)


type Wallet struct {
	ActiveWorkers   uint64 `json:"activeWorkers",bson:"activeWorkers"`
	Address			string `json:"address",bson:"address"`
	AverageHashrate float64 `json:"averageHashrate",bson:"averageHashrate"`
	Balance         float64 `json:"balance",bson:"balance"`
	Delta           float64 `json:"address",bson:"address"`
	Paid            float64 `json:"paid",bson:"paid"`
	Time            int64 `json:"time",bson:"time"`
}


type Connection struct {
	session *mgo.Session
}

var sess *mgo.Session
var lock sync.Once

func GetConnection(url string) (*Connection, error) {
	lock.Do(func() {
		s, err := mgo.Dial(url)
		if err != nil {
			logrus.Fatal(err)
		}
		sess = s
	})
	c := &Connection{
		session: sess.Copy(),
	}
	return c, nil
}

func (c *Connection) InsertWalletState(w Wallet) error {
	s := c.session.DB("Wallet").C("eth")

	// calculate delta

	err := s.Insert(w)
	if err != nil {
		logrus.Warningf("Error when inserting to Wallet: %s", err)
		return err
	}
	return nil
}

func (c *Connection) GetWallet(address string) (Wallet, error) {
	s := c.session.DB("Wallet").C("eth")

	var w Wallet
	err := s.Find(bson.M{"address": address}).Sort("-time").One(&w)
	if err != nil {
		logrus.Warningf("Error while fetching wallet: %s", err)
		return Wallet{}, err
	}
	return w, nil
}