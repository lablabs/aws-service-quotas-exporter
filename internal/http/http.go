package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net"
	"net/http"
)

func NewHTTP(log *logrus.Logger, address string, registry *prometheus.Registry) (*HTTP, error) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	handler := http.NewServeMux()
	s := http.Server{
		Handler: handler,
	}
	h := HTTP{
		log: log,
		ln:  ln,
		s:   &s,
	}
	RegisterMetricEndpoint(handler, registry)
	RegisterStatusEndpoint(handler, &h)
	return &h, nil
}

type HTTP struct {
	log *logrus.Logger
	ln  net.Listener
	s   *http.Server
}

func (h *HTTP) Run(ctx context.Context) error {
	h.log.Infof("start http endpoint: %s", h.ln.Addr())
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		err := h.s.Serve(h.ln)
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("unable to serve http enpoint: %w", err)
		}
		return nil
	})
	<-ctx.Done()
	h.log.Debugf("http endpoint exiting...")
	err := h.s.Shutdown(ctx)
	if err != nil {
		return err
	}
	err = g.Wait()
	if err != nil {
		return err
	}
	h.log.Infof("http endpoint exit OK")
	return nil
}

func (h *HTTP) Status(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func RegisterMetricEndpoint(mux *http.ServeMux, registry *prometheus.Registry) {
	registry.MustRegister(collectors.NewGoCollector())
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))
}

func RegisterStatusEndpoint(mux *http.ServeMux, h *HTTP) {
	mux.HandleFunc("/liveness", h.Status)
	mux.HandleFunc("/readiness", h.Status)
}
