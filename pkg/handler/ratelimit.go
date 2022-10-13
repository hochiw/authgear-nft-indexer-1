package handler

import (
	"fmt"

	"github.com/authgear/authgear-server/pkg/lib/ratelimit"
	"github.com/authgear/authgear-server/pkg/util/duration"
)

func AntiSpamProbeCollectionRequestBucket(appID string) ratelimit.Bucket {
	return ratelimit.Bucket{
		Key:         fmt.Sprintf("probe-collection-request-%s", appID),
		Size:        60,
		ResetPeriod: duration.PerHour,
	}
}
