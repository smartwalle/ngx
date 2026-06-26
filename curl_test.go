package ngx_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/smartwalle/ngx"
)

func TestRequest_CURL_OpenMeteoForecast(t *testing.T) {
	var req = ngx.NewRequest(http.MethodGet, "https://api.open-meteo.com/v1/forecast")
	req.Query.Set("latitude", "31.2304")
	req.Query.Set("longitude", "121.4737")
	req.Query.Set("current_weather", "true")
	req.Query.Set("timezone", "Asia/Shanghai")

	var got, err = req.CURL()
	if err != nil {
		t.Fatalf("CURL returned error: %v", err)
	}

	t.Log(got)

	var wants = []string{
		"--request 'GET' ",
		"'https://api.open-meteo.com/v1/forecast?current_weather=true&latitude=31.2304&longitude=121.4737&timezone=Asia%2FShanghai' ",
	}
	for _, want := range wants {
		if !strings.Contains(got, want) {
			t.Fatalf("curl command missing %q, got: %s", want, got)
		}
	}
}
