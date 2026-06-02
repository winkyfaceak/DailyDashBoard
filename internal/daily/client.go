package daily

import (
	"net/http"
	"time"
)

const apiUserAgent = "DailyDashBoard/1.0"

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}
