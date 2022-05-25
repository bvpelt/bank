package domain

import "time"

type DbpoolStat struct {
	AcquireCount         int64         `json:"acquirecount"`
	AcquireDuration      time.Duration `json:"acquireduration"`
	AcquireConns         int32         `json:"acquireconns"`
	CanceledAcquireCount int64         `json:"canceledacquirecount"`
	ConstructingConns    int32         `json:"constructingconns"`
	EmptyAcquireCount    int64         `json:"emptyacquirecount"`
	IdleConns            int32         `json:"idleconns"`
	MaxConns             int32         `json:"maxconns"`
	TotalConns           int32         `json:"totalconns"`
}
