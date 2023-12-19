package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Статусы операций по cache
	keycloakCacheCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "ord",
		Subsystem: "site_client_process",
		Name:      "keycloak_cache_counter",
		Help:      "Count keycloak cache events",
	}, []string{"status"})
)

func Init(ctx context.Context, mux *http.ServeMux) *http.ServeMux {
	// Заводим метрики
	prometheus.MustRegister(
		keycloakCacheCounter,
	)

	// Роут по которому будет стучаться
	mux.Handle("/metrics", promhttp.Handler())

	return mux
}

// Записываем размер файла в бакеты
func IncKeycloakCacheEvent(status string) {
	keycloakCacheCounter.WithLabelValues(status).Inc()
}
