package metrics

import (
	"expvar"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

type (
	Metrics struct {
		Goroutines *expvar.Int
		Requests   *expvar.Int
		Successes  *expvar.Int
		Errors     *expvar.Int
		DbQueries  *expvar.Int
		timer      *time.Ticker
	}
)

func New() *Metrics {
	metric := Metrics{
		Goroutines: expvar.NewInt("go_goroutines"),
		Requests:   expvar.NewInt("http_requests_total"),
		Successes:  expvar.NewInt("http_responses_success_total"),
		Errors:     expvar.NewInt("http_responses_error_total"),
		DbQueries:  expvar.NewInt("db_queries_total"),
		timer:      time.NewTicker(time.Second * 10),
	}

	go func() {
		for {
			<-metric.timer.C
			metric.Goroutines.Set(int64(runtime.NumGoroutine()))
		}
	}()
	return &metric
}

func (m *Metrics) IncreaseDbRequests(*gin.Context) {
	m.DbQueries.Add(1)
}
