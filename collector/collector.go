package collector

import (
	"context"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/alexruf/fritzmond/config"
	"github.com/alexruf/fritzmond/fritzbox"
	"github.com/alexruf/fritzmond/metrics"
	"github.com/nakabonne/tstorage"
)

type Collector struct {
	ctx      context.Context
	cfg      config.Config
	fritzbox fritzbox.Fritzbox
	storage  tstorage.Storage
}

func New(ctx context.Context, cfg config.Config, fritzbox fritzbox.Fritzbox, storage tstorage.Storage) Collector {
	return Collector{
		ctx:      ctx,
		cfg:      cfg,
		fritzbox: fritzbox,
		storage:  storage,
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
	var rows []tstorage.Row
	commonLinkProperties, err := c.fritzbox.GetCommonLinkProperties()
	if err != nil {
		log.Println(err.Error())
	} else if commonLinkProperties != nil {
		rows = append(rows,
			tstorage.Row{
				Metric: metrics.MetricPhysicalLinkStatus,
				DataPoint: tstorage.DataPoint{
					Value:     metrics.ConvertPhysicalLinkStatus(commonLinkProperties.NewPhysicalLinkStatus),
					Timestamp: timestamp,
				},
			},
			tstorage.Row{
				Metric: metrics.MetricLayer1DownstreamMaxBitRate,
				DataPoint: tstorage.DataPoint{
					Value:     float64(commonLinkProperties.NewLayer1DownstreamMaxBitRate),
					Timestamp: timestamp,
				},
			},
			tstorage.Row{
				Metric: metrics.MetricLayer1UpstreamMaxBitRate,
				DataPoint: tstorage.DataPoint{
					Value:     float64(commonLinkProperties.NewLayer1UpstreamMaxBitRate),
					Timestamp: timestamp,
				},
			})
	}

	totalBytesSent, err := c.fritzbox.GetTotalBytesSent()
	if err != nil {
		log.Println(err.Error())
	} else {
		rows = append(rows, tstorage.Row{
			Metric: metrics.MetricTotalBytesSent,
			DataPoint: tstorage.DataPoint{
				Value:     float64(totalBytesSent),
				Timestamp: timestamp,
			},
		})
	}

	totalBytesReceived, err := c.fritzbox.GetTotalBytesReceived()
	if err != nil {
		log.Println(err.Error())
	} else {
		rows = append(rows, tstorage.Row{
			Metric: metrics.MetricTotalBytesReceived,
			DataPoint: tstorage.DataPoint{
				Value:     float64(totalBytesReceived),
				Timestamp: timestamp,
			},
		})
	}

	totalPacketsSent, err := c.fritzbox.GetTotalPacketsSent()
	if err != nil {
		log.Println(err.Error())
	} else {
		rows = append(rows, tstorage.Row{
			Metric: metrics.MetricTotalPacketsSent,
			DataPoint: tstorage.DataPoint{
				Value:     float64(totalPacketsSent),
				Timestamp: timestamp,
			},
		})
	}

	totalPacketsReceived, err := c.fritzbox.GetTotalPacketsReceived()
	if err != nil {
		log.Println(err.Error())
	} else {
		rows = append(rows, tstorage.Row{
			Metric: metrics.MetricTotalPacketsReceived,
			DataPoint: tstorage.DataPoint{
				Value:     float64(totalPacketsReceived),
				Timestamp: timestamp,
			},
		})
	}

	if len(rows) > 0 {
		log.Printf("Inserting %d metrics into database.\n", len(rows))
		err = c.storage.InsertRows(rows)
		if err != nil {
			log.Printf("Error inserting rows into database: %s\n", err.Error())
		}
	}
}
