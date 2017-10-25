package storage

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Wallet struct {
	ActiveWorkers   int64         `json:"activeWorkers" bson:"activeWorkers"`
	Address         string        `json:"address" bson:"address"`
	AverageHashrate float64       `json:"averageHashrate" bson:"averageHashrate"`
	Balance         int64         `json:"balance" bson:"balance"`
	Delta           int64         `json:"delta" bson:"delta"`
	Id              bson.ObjectId `json:"-" bson:"_id,omitempty"`
	Paid            int64         `json:"paid" bson:"paid"`
	Time            int64         `json:"time" bson:"time"`
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
	// initialize objectid
	//w.Id = bson.NewObjectId()

	// calculate delta
	var wOld Wallet
	err := s.Find(bson.M{"address": w.Address}).Sort("-_id").One(&wOld)
	if err != nil {
		if "not found" == fmt.Sprintf("%s", err) {
			logrus.Infof("Creating new wallet: %s", w.Address)
			w.Delta = w.Balance
			err = s.Insert(w)
			if err != nil {
				logrus.Warningf("Error when inserting initial wallet: %s", err)
				return err
			}
			return nil
		}
		logrus.Warningf("Insertion Error while fetching wallet: %s", err)
		return err
	}

	w.Paid = wOld.Paid
	// Testing for new cycle
	logrus.Infof("Old: ID: %s, Time: %d, Delta: %d, Balance: %d", wOld.Id, wOld.Time, wOld.Delta, wOld.Balance)
	if (w.Balance - wOld.Balance) < 0 {
		w.Paid += 1000000000000000000
		d := wOld.Balance - 1000000000000000000
		if d < 0 {
			d = -d
		}
		w.Delta = d + w.Balance
	} else {
		w.Delta = w.Balance - wOld.Balance
	}
	logrus.Infof("New: ID: %s, Time: %d, Delta: %d, Balance: %d", w.Id, w.Time, w.Delta, w.Balance)

	err = s.Insert(w)
	if err != nil {
		logrus.Warningf("Error when inserting to Wallet: %s", err)
		return err
	}
	return nil
}

func (c *Connection) GetWallet(address string) (Wallet, error) {
	s := c.session.DB("Wallet").C("eth")

	var w Wallet
	err := s.Find(bson.M{"address": address}).Sort("-_id").One(&w)
	if err != nil {
		logrus.Warningf("Error while fetching wallet: %s", err)
		return Wallet{}, err
	}
	return w, nil
}
