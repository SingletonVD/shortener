package worker

import (
	"context"
	"time"

	"github.com/SingletonVD/shortener/internal/logger"
	"github.com/SingletonVD/shortener/internal/model"
	"go.uber.org/zap"
)

const (
	bufferSize       = 1024
	batchSize        = 100
	scheduleInterval = 10
)

type DeleteLinkRepository interface {
	DeleteBatch(ctx context.Context, links []model.DeleteLink) error
}

type DeleteWorker struct {
	queue          chan model.DeleteLink
	linkRepository DeleteLinkRepository
}

func NewDeleteWorker(linkRepository DeleteLinkRepository) *DeleteWorker {
	return &DeleteWorker{
		queue:          make(chan model.DeleteLink, bufferSize),
		linkRepository: linkRepository,
	}
}

func (worker *DeleteWorker) Enqueue(links []model.DeleteLink) {
	for _, link := range links {
		worker.queue <- link
	}
}

func (worker *DeleteWorker) Schedule(ctx context.Context) {
	ticker := time.NewTicker(scheduleInterval * time.Second)

	var linksBatch []model.DeleteLink

	for {
		select {
		case link := <-worker.queue:
			linksBatch = append(linksBatch, link)
			if len(linksBatch) == batchSize {
				linksBatch = worker.deleteBatch(ctx, linksBatch)
			}
		case <-ticker.C:
			if len(linksBatch) == 0 {
				continue
			}
			linksBatch = worker.deleteBatch(ctx, linksBatch)
		case <-ctx.Done():
			return
		}
	}
}

func (worker *DeleteWorker) deleteBatch(ctx context.Context, linksBatch []model.DeleteLink) []model.DeleteLink {
	err := worker.linkRepository.DeleteBatch(ctx, linksBatch)
	if err != nil {
		logger.Log.Debug("cannot delete links", zap.Error(err))
		return linksBatch
	}
	return nil
}
