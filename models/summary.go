package models

import (
	"github.com/lPoltergeist/rinha-backend.git/data"
	"github.com/shopspring/decimal"
)

type SummaryData struct {
	DefaultCount   int
	DefaultAmount  decimal.Decimal
	FallbackCount  int
	FallbackAmount decimal.Decimal
}

type EndpointSummary struct {
	TotalRequests int     `json:"totalRequests"`
	TotalAmount   float64 `json:"totalAmount"`
}

type ResponseSummary struct {
	Default  EndpointSummary `json:"default"`
	Fallback EndpointSummary `json:"fallback"`
}

func (s *SummaryData) Add(amount float64, isFallback bool) error {

	context := data.Context
	client := data.Client

	pipe := client.TxPipeline()

	decAmount := decimal.NewFromFloat(amount).Round(2)
	amountFloat := decAmount.InexactFloat64()

	if isFallback {
		pipe.Incr(context, "summary:fallback:count")
		pipe.IncrByFloat(context, "summary:fallback:amount", amountFloat)
	} else {
		pipe.Incr(context, "summary:default:count")
		pipe.IncrByFloat(context, "summary:default:amount", amountFloat)
	}

	_, err := pipe.Exec(context)
	return err
}

func (s *SummaryData) BuildSummary() ResponseSummary {
	context := data.Context
	client := data.Client

	var defaultCount, fallbackCount int
	var defaultAmount, fallbackAmount decimal.Decimal

	if val, err := client.Get(context, "summary:default:count").Int(); err == nil {
		defaultCount = val
	}
	if val, err := client.Get(context, "summary:default:amount").Result(); err == nil {
		defaultAmount, _ = decimal.NewFromString(val)
	}

	if val, err := client.Get(context, "summary:fallback:count").Int(); err == nil {
		fallbackCount = val
	}

	if val, err := client.Get(context, "summary:fallback:amount").Result(); err == nil {
		fallbackAmount, _ = decimal.NewFromString(val)
	}

	return ResponseSummary{
		Default: EndpointSummary{
			TotalRequests: defaultCount,
			TotalAmount:   defaultAmount.Round(2).InexactFloat64(),
		},
		Fallback: EndpointSummary{
			TotalRequests: fallbackCount,
			TotalAmount:   fallbackAmount.Round(2).InexactFloat64(),
		},
	}
}
