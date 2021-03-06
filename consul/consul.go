package consul

import (
	"bytes"
	"sync"
	"time"

	"github.com/fvdveen/mu2-config"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type provider struct {
	client *api.Client
	key    string
	typ    string
	qOpts  *api.QueryOptions

	ch    chan *config.Config
	chsMu sync.Mutex

	errChan       chan error
	quitChan      chan interface{}
	quitWatchChan chan interface{}
}

// NewProvider creates a new provider with a consul backend
func NewProvider(c *api.Client, key string, t string, qopts *api.QueryOptions, l func() *zap.SugaredLogger) (config.Provider, error) {
	p := &provider{
		client:        c,
		key:           key,
		typ:           t,
		qOpts:         qopts,
		errChan:       make(chan error),
		quitChan:      make(chan interface{}),
		quitWatchChan: make(chan interface{}),
		ch:            make(chan *config.Config),
	}

	go func(p *provider) {
		l().Named("provider").Named("consul").Debugf("Starting...")
		var last string
		t := time.Tick(time.Second * 5)

		v := viper.New()
		v.SetConfigType(p.typ)

		var c *config.Config

		for {
			select {
			case <-p.quitWatchChan:
				l().Named("provider").Named("consul").Debugf("Stopping...")
				close(p.ch)

				close(p.errChan)
				return
			case <-t:
				kv, _, err := p.client.KV().Get(p.key, p.qOpts)
				if err != nil {
					l().Named("provider").Named("consul").Errorf("Get config %s: %v", p.key, err)
					continue
				} else if kv == nil {
					l().Named("provider").Named("consul").Warnf("Key %s does not exist", p.key)
					continue
				}

				if string(kv.Value) == last {
					continue
				}

				last = string(kv.Value)

				if err := v.ReadConfig(bytes.NewBuffer(kv.Value)); err != nil {
					l().Named("provider").Named("consul").Errorf("Read in config: %v", err)
					continue
				}

				c = config.New()

				if err := v.Unmarshal(c); err != nil {
					l().Named("provider").Named("consul").Errorf("Unmarshal into config: %v", err)
					continue
				}

				p.ch <- c
			}
		}
	}(p)

	go func(p *provider) {
		<-p.quitChan
		p.quitWatchChan <- 0
	}(p)

	return p, nil
}

// Watch watches consul for an update at the given key and then sends the updated value over the channel
func (p *provider) Watch() <-chan *config.Config {
	return p.ch
}

// Close closes the provider
func (p *provider) Close() error {
	p.quitChan <- 0
	return nil
}
