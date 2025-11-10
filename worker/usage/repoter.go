package usage

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/pubudu2003060/go-proxy-prototype/worker/models"
)

type UsageRepoter struct {
	captainURL string
}

func NewUsageReporter(captainURL string) *UsageRepoter {
	return &UsageRepoter{
		captainURL: captainURL,
	}
}

func (r *UsageRepoter) ReportUsage(userID string, _bytes int64) {
	go func() {
		reqBody := models.UsageReport{
			UserID: userID,
			Bytes: _bytes,
		}

		jsonData,err := json.Marshal(reqBody)
		if err != nil {
			log.Printf("Failed to marshal usage report: %v", err)
			return
		}

		resp,err := http.Post(r.captainURL+"api/v1/usage","application/json",bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Failed to report usage: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Usage report failed with status: %d", resp.StatusCode)
		}
	}()
}