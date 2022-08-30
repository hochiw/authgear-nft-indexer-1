// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package handler

import (
	"github.com/authgear/authgear-nft-indexer/pkg/ratelimit"
	"github.com/authgear/authgear-server/pkg/lib/infra/redis/appredis"
	ratelimit2 "github.com/authgear/authgear-server/pkg/lib/ratelimit"
	"github.com/authgear/authgear-server/pkg/util/clock"
	"github.com/authgear/authgear-server/pkg/util/log"
)

// Injectors from wire.go:

func NewRateLimiterFactory(lf *log.Factory, redis *appredis.Handle) ratelimit.Factory {
	logger := ratelimit2.NewLogger(lf)
	clock := _wireSystemClockValue
	factory := ratelimit.Factory{
		Logger: logger,
		Redis:  redis,
		Clock:  clock,
	}
	return factory
}

var (
	_wireSystemClockValue = clock.NewSystemClock()
)
