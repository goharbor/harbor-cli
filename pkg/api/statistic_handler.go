package api

import (
	 
	"fmt"
	"net/http"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/statistic"
 
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type StatisticHandler struct {
	client *statistic.Client
}

func NewStatisticHandler() (*StatisticHandler, error) {
	_, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create statistic handler: %w", err)
	}

	return &StatisticHandler{client: client.Statistic}, nil
}

func (h *StatisticHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	params := &statistic.GetStatisticParams{
		Context: r.Context(),
	}

	stats, err := h.client.GetStatistic(params.Context, params)
	if err != nil {
		log.Errorf("failed to retrieve statistics: %v", err)
		http.Error(w, "failed to retrieve statistics", http.StatusInternalServerError)
		return
	}

	responseData, err := stats.GetPayload().MarshalBinary()
	if err != nil {
		log.Errorf("failed to marshal statistics: %v", err)
		http.Error(w, "failed to process statistics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}
