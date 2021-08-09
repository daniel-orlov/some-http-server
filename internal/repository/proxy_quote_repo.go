package repository

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"some-http-server/internal/types"
)

type ProxyQuoteRepo struct {
	repo *QuoteRepository
	r    RedisConection //TODO add redis connection
}

func NewProxyQuoteRepo(repo *QuoteRepository, r RedisConection) *ProxyQuoteRepo {
	return &ProxyQuoteRepo{repo: repo, r: r}
}

func (p *ProxyQuoteRepo) Save(ctx context.Context, quote *types.FullQuoteData) (string, error) {
	return p.repo.Save(ctx, quote)
}

func (p *ProxyQuoteRepo) Read(ctx context.Context, id string) (*types.FullQuoteData, error) {
	log := ctxlogrus.Extract(ctx)
	log.Trace("'Read' hit the cache")
	// TODO implement caching
	log.Trace("'Read' fell through the cache")

	return p.repo.Read(ctx, id)
}
