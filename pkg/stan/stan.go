package stan

import (
	"context"
	"time"

	"github.com/carlescere/goback"
	stan "github.com/nats-io/stan.go"
	"github.com/qilin/go-core/logger"
	"github.com/qilin/go-core/provider"
)

type Stan struct {
	ctx context.Context
	cfg Config
	provider.LMT
}

type Config struct {
	Debug         bool `fallback:"shared.debug"`
	Backoff       goback.SimpleBackoff
	StanClusterID string
	StanClientID  string
	StanOptions   stan.Options
}

type ConnCfg struct {
	Backoff      goback.SimpleBackoff
	StanClientID string
}

type ConnHandler func(ctx context.Context, sc stan.Conn, log logger.Logger) error

// BackoffCopy
func (s *Stan) BackoffCopy() goback.SimpleBackoff {
	b := s.cfg.Backoff
	cb := &b
	return *cb
}

// SimpleConnect
func (s *Stan) SimpleConnect(ctx context.Context, cfg ConnCfg, log logger.Logger) (stan.Conn, error) {
	// TODO: Pool connection
	return stan.Connect(s.cfg.StanClusterID, s.cfg.StanClientID+cfg.StanClientID,
		stan.NatsURL(s.cfg.StanOptions.NatsURL),
		stan.SetConnectionLostHandler(func(_ stan.Conn, e error) {
			log.Error("connect to NATS Streaming server lost: %v", logger.Args(e))
		}),
	)
}

// SmartConnect
func (s *Stan) SmartConnect(ctx context.Context, cfg ConnCfg, log logger.Logger, handler ConnHandler) {
	// TODO: Pool connection
	log = log.WithFields(logger.Fields{"service": Prefix})
	b := cfg.Backoff
	cb := &b
	for {
		ctxStan, cancel := context.WithCancel(context.Background())
		sc, e := stan.Connect(s.cfg.StanClusterID, s.cfg.StanClientID+cfg.StanClientID,
			stan.NatsURL(s.cfg.StanOptions.NatsURL),
			stan.SetConnectionLostHandler(func(_ stan.Conn, e error) {
				log.Error("connect to NATS Streaming server lost: %v", logger.Args(e))
				cancel()
			}),
		)
		if e != nil {
			log.Error("connect to NATS Streaming server failed: %v", logger.Args(e))
			// Next attempt
			d, e := cb.NextAttempt()
			if e != nil {
				log.Error("backoff error: %v", logger.Args(e))
			}
			if d < 0 {
				d = 0
			}
			time.Sleep(d)
			continue
		}
		cb.Reset()
		if e = handler(ctx, sc, log); e != nil {
			log.Error("handler error: %v", logger.Args(e))
			goto nextAttempt
		}
		log.Debug("connected to NATS Streaming server, waiting of signal")
		// graceful shutdown
		select {
		case <-ctx.Done():
		case <-ctxStan.Done():
		case <-s.ctx.Done():
		}
	nextAttempt:
		if e := sc.Close(); e != nil {
			log.Error("close connection to NATS Streaming server failed: %v", logger.Args(e))
		}
		if ctx.Err() != nil || s.ctx.Err() != nil {
			log.Debug("close connection to NATS Streaming server")
			return
		}
		log.Debug("retry to reconnect to NATS Streaming server")
		d, e := cb.NextAttempt()
		if e != nil {
			log.Error("backoff error: %v", logger.Args(e))
		}
		// fix bug with negative time
		if d < 0 {
			cb.Reset()
			d = 0
		}
		time.Sleep(d)
		continue
	}
}

// New
func New(ctx context.Context, set provider.AwareSet, cfg *Config) *Stan {
	set.Logger = set.Logger.WithFields(logger.Fields{"service": Prefix})
	return &Stan{
		ctx: ctx,
		cfg: *cfg,
		LMT: &set,
	}
}
