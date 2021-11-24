package collector

import (
	"context"
	"log"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/alexruf/fritzmond/config"
	"github.com/alexruf/fritzmond/fritzbox"
	"github.com/alexruf/fritzmond/metrics"
	bolt "go.etcd.io/bbolt"
)

type Collector struct {
	ctx      context.Context
	cfg      config.Config
	fritzbox fritzbox.Fritzbox
	db       *bolt.DB
}

func New(ctx context.Context, cfg config.Config, fritzbox fritzbox.Fritzbox, db *bolt.DB) Collector {
	return Collector{
		ctx:      ctx,
		cfg:      cfg,
		fritzbox: fritzbox,
		db:       db,
	}
}

func (c Collector) Start(wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	c.run()
	ticker := time.NewTicker(time.Duration(c.cfg.Interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			c.run()
		case <-c.ctx.Done():
			log.Println("Collector: Stopped")
			return
		}
	}
}

func (c Collector) run() {
	fritzUrl, err := url.Parse(c.cfg.Url)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Printf("Start collecting data from FRITZ!Box (%s)\n", fritzUrl.String())

	timestamp := time.Now().Unix()
	commonLinkProperties, err := c.fritzbox.GetCommonLinkProperties()
	if err != nil {
		log.Println(err.Error())
	} else if commonLinkProperties != nil {
		if err := c.db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(metrics.MetricPhysicalLinkStatus))
			if err != nil {
				return err
			}
			err = b.Put([]byte(strconv.FormatInt(timestamp, 10)), []byte(commonLinkProperties.NewPhysicalLinkStatus))
			if err != nil {
				return err
			}

			b, err = tx.CreateBucketIfNotExists([]byte(metrics.MetricLayer1DownstreamMaxBitRate))
			if err != nil {
				return err
			}
			err = b.Put([]byte(strconv.FormatInt(timestamp, 10)), []byte(strconv.FormatUint(uint64(commonLinkProperties.NewLayer1DownstreamMaxBitRate), 10)))
			if err != nil {
				return err
			}

			b, err = tx.CreateBucketIfNotExists([]byte(metrics.MetricLayer1UpstreamMaxBitRate))
			if err != nil {
				return err
			}
			err = b.Put([]byte(strconv.FormatInt(timestamp, 10)), []byte(strconv.FormatUint(uint64(commonLinkProperties.NewLayer1UpstreamMaxBitRate), 10)))
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			log.Printf("Error updating database: %s\n", err.Error())
		}
	}

	totalBytesSent, err := c.fritzbox.GetTotalBytesSent()
	if err != nil {
		log.Println(err.Error())
	} else {
		if err := c.db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(metrics.MetricTotalBytesSent))
			if err != nil {
				return err
			}
			err = b.Put([]byte(strconv.FormatInt(timestamp, 10)), []byte(strconv.FormatUint(uint64(totalBytesSent), 10)))
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			log.Printf("Error updating database: %s\n", err.Error())
		}
	}

	totalBytesReceived, err := c.fritzbox.GetTotalBytesReceived()
	if err != nil {
		log.Println(err.Error())
	} else {
		if err := c.db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(metrics.MetricTotalBytesReceived))
			if err != nil {
				return err
			}
			err = b.Put([]byte(strconv.FormatInt(timestamp, 10)), []byte(strconv.FormatUint(uint64(totalBytesReceived), 10)))
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			log.Printf("Error updating database: %s\n", err.Error())
		}
	}

	totalPacketsSent, err := c.fritzbox.GetTotalPacketsSent()
	if err != nil {
		log.Println(err.Error())
	} else {
		if err := c.db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(metrics.MetricTotalPacketsSent))
			if err != nil {
				return err
			}
			err = b.Put([]byte(strconv.FormatInt(timestamp, 10)), []byte(strconv.FormatUint(uint64(totalPacketsSent), 10)))
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			log.Printf("Error updating database: %s\n", err.Error())
		}
	}

	totalPacketsReceived, err := c.fritzbox.GetTotalPacketsReceived()
	if err != nil {
		log.Println(err.Error())
	} else {
		if err := c.db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(metrics.MetricTotalPacketsReceived))
			if err != nil {
				return err
			}
			err = b.Put([]byte(strconv.FormatInt(timestamp, 10)), []byte(strconv.FormatUint(uint64(totalPacketsReceived), 10)))
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			log.Printf("Error updating database: %s\n", err.Error())
		}
	}
}
