package collector

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/alexruf/fritzmond/config"
	"github.com/alexruf/fritzmond/fritzbox"
)

type Collector struct {
	ctx      context.Context
	cfg      config.Config
	fritzbox fritzbox.Fritzbox
}

func New(ctx context.Context, cfg config.Config, fritzbox fritzbox.Fritzbox) Collector {
	return Collector{
		ctx:      ctx,
		cfg:      cfg,
		fritzbox: fritzbox,
	}
}

func (c Collector) Start(wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

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
	log.Println("Collecting data from FRITZ!Box")
	data, err := c.fritzbox.GetCommonLinkProperties()
	if err != nil {
		log.Println(err.Error())
	}
	if data != nil {
		log.Printf("GetCommonLinkProperties:\n%s\n", data)
	}
}
