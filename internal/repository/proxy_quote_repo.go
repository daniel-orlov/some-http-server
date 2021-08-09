package repository

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"some-http-server/internal/cache"
	"some-http-server/internal/service"
	"some-http-server/internal/types"
)

type ProxyQuoteRepo struct {
	repo  service.QuoteRepo
	cache cache.Store
}

func NewProxyQuoteRepo(repo service.QuoteRepo, cache cache.Store) service.QuoteRepo {
	if cache == nil {
		return repo
	}

	return &ProxyQuoteRepo{repo: repo, cache: cache}
}

func (p *ProxyQuoteRepo) Save(ctx context.Context, quote *types.FullQuoteData) (string, error) {
	return p.repo.Save(ctx, quote)
}

func (p *ProxyQuoteRepo) Read(ctx context.Context, id, accountID string) (*types.FullQuoteData, error) {
	log := ctxlogrus.Extract(ctx)

	if fqd, ok := p.getReadFromCache(ctx, id, accountID); ok {
		log.Trace("'Read' hit the cache")
		return &fqd, nil
	}
	log.Trace("'Read' missed the cache, getting data directly from the db")

	result, err := p.repo.Read(ctx, id, accountID)
	if err != nil {
		return nil, err
	}

	p.saveReadIntoCache(ctx, result, id, accountID)

	return result, nil
}

func (p *ProxyQuoteRepo) getReadFromCache(ctx context.Context, id, accountID string) (types.FullQuoteData, bool) {
	log := ctxlogrus.Extract(ctx)

	cachedFQD, ok := p.cache.Get(getReadCacheKey(id, accountID))
	if !ok {
		return types.FullQuoteData{}, false
	}

	fqd, ok := cachedFQD.(types.FullQuoteData)
	if !ok {
		log.Errorf("Got data of type %T but wanted types.FullQuoteData", cachedFQD)
		return types.FullQuoteData{}, false
	}

	return fqd, ok
}

func getReadCacheKey(id, accountID string) string {
	return fmt.Sprintf("%s:%s", accountID, id)
}

func (p *ProxyQuoteRepo) saveReadIntoCache(ctx context.Context, fqd *types.FullQuoteData, id, accountID string) {
	log := ctxlogrus.Extract(ctx)

	if !p.cache.SetWithTTL(getReadCacheKey(id, accountID), *fqd) {
		log.Tracef("saveReadIntoCache call SetWithTTL dropped by cache: %v", *fqd)
	}

	log.Tracef("Successfully saved Read in cache with ttl: %v", *fqd)
}
