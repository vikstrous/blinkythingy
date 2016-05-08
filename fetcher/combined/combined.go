package combined

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/vikstrous/blinkythingy"
	"github.com/vikstrous/blinkythingy/fetcher"
)

type combinedFetcher struct {
	fetchers []fetcher.Fetcher
}

func New(fetchers ...fetcher.Fetcher) fetcher.Fetcher {
	return &combinedFetcher{fetchers: fetchers}
}

func (cf *combinedFetcher) FetchStatuses() error {
	errs := []error{}
	for i, fetcher := range cf.fetchers {
		logrus.Debugf("fetching %d", i)
		err := fetcher.FetchStatuses()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", errs)
	}
	return nil
}

func (cf *combinedFetcher) ListStatuses() []blinkythingy.Color {
	combinedStatus := []blinkythingy.Color{}
	for _, fetcher := range cf.fetchers {
		combinedStatus = append(combinedStatus, fetcher.ListStatuses()...)
	}
	return combinedStatus
}
