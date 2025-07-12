package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var HealthStatus = make(map[string]interface{})

func StartHealthCheck() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("🔥 panic em VerifyHealth: %v\n", r)
			}
		}()

		for {
			fmt.Println("✅ chamou VerifyHealth")
			VerifyHealth()
			time.Sleep(5 * time.Second)
		}
	}()
}

func VerifyHealth() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("🔥 panic em VerifyHealth: %v\n", r)
		}
	}()
	type healthResponse struct {
		Failure bool `json:"failure"`
	}

	endpoints := map[string]string{
		"default":  "http://payment-processor-default:8080/payments/service-health",
		"fallback": "http://payment-processor-fallback:8080/payments/service-health",
	}

	for name, url := range endpoints {
		resp, err := http.Get(url)
		if err != nil {
			HealthStatus[name] = map[string]interface{}{
				"endpoint": url,
				"failure":  "request_error: " + err.Error(),
			}
			panic(err)

		}
		defer resp.Body.Close()

		var result healthResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			HealthStatus[name] = map[string]interface{}{
				"endpoint": url,
				"failure":  "decode_error: " + err.Error(),
			}
			panic(err)
		}

		HealthStatus[name] = map[string]interface{}{
			"endpoint": url,
			"failure":  result.Failure,
		}
	}
}

func ReturnHealthyEndpoint() string {

	defaultEntry, _ := HealthStatus["default"].(map[string]interface{})
	fallbackEntry, _ := HealthStatus["fallback"].(map[string]interface{})

	defaultFailure, _ := defaultEntry["failure"].(bool)
	fallbackFailure, _ := fallbackEntry["failure"].(bool)

	if defaultFailure {
		return "http://payment-processor-fallback:8080/payments"
	} else if fallbackFailure {
		return "http://payment-processor-default:8080/payments"
	} else {
		return "http://payment-processor-default:8080/payments"
	}

}
