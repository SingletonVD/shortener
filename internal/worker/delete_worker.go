package worker

import (
	"context"
	"time"

	"github.com/SingletonVD/shortener/internal/concurrency"
	"github.com/SingletonVD/shortener/internal/logger"
	"github.com/SingletonVD/shortener/internal/model"
	"go.uber.org/zap"
)

type DeleteLinkRepository interface {
	DeleteBatch(ctx context.Context, links []model.DeleteLink) error
}

type DeleteWorkerConfig struct {
	BufferSize        int
	BatchSize         int
	ScheduleInterval  int
	ConcurrentWriters int
}

type DeleteWorker struct {
	queue            chan model.DeleteLink
	writersSemaphore concurrency.Semaphore
	linkRepository   DeleteLinkRepository
	config           DeleteWorkerConfig
}

func NewDeleteWorker(linkRepository DeleteLinkRepository, config DeleteWorkerConfig) *DeleteWorker {
	return &DeleteWorker{
		queue:            make(chan model.DeleteLink, config.BufferSize),
		writersSemaphore: *concurrency.NewSemaphore(config.ConcurrentWriters),
		linkRepository:   linkRepository,
		config:           config,
	}
}

func (worker *DeleteWorker) Enqueue(links []model.DeleteLink) {
	go func() {
		worker.writersSemaphore.Acquire()
		for _, link := range links {
			worker.queue <- link
		}
		worker.writersSemaphore.Release()
	}()
}

func (worker *DeleteWorker) Schedule(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(worker.config.ScheduleInterval) * time.Second)

	var linksBatch []model.DeleteLink

	for {
		select {
		case link := <-worker.queue:
			linksBatch = append(linksBatch, link)
			if len(linksBatch) >= worker.config.BatchSize {
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
