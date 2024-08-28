package workerpool

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"sync"
	"time"

	"github.com/pingooio/stdx/log/slogx"
	"github.com/pingooio/stdx/queue"
)

type JobHandler[I any] func(ctx context.Context, input I) (err error)

type internalJobHandler = func(ctx context.Context, payload []byte) (err error)

type WorkerPool struct {
	queue          queue.Queue
	concurrencyMax uint32
	jobHandlers    map[string]internalJobHandler
	logger         *slog.Logger
}

type Options struct {
	// default: 200
	ConcurrencyMax uint32
	Logger         *slog.Logger
}

func NewPool(queue queue.Queue, options *Options) (worker *WorkerPool, err error) {
	opts := Options{
		ConcurrencyMax: 200,
		Logger:         slog.New(slogx.NewDiscardHandler()),
	}

	if options.ConcurrencyMax != 0 {
		if options.ConcurrencyMax > math.MaxInt32 {
			err = fmt.Errorf("workerpool: concurrencyMax can't be > %d", math.MaxInt32)
			return
		}

		opts.ConcurrencyMax = options.ConcurrencyMax
	}

	if options.Logger != nil {
		opts.Logger = options.Logger
	}

	worker = &WorkerPool{
		queue:       queue,
		jobHandlers: make(map[string]internalJobHandler),

		concurrencyMax: opts.ConcurrencyMax,
		logger:         opts.Logger,
	}
	return
}

func AddHandler[T any](workerPool *WorkerPool, jobType string, handler JobHandler[T]) {
	if _, exists := workerPool.jobHandlers[jobType]; exists {
		panic(fmt.Sprintf("workerpool: job handler already exists for %s", jobType))
	}

	workerPool.jobHandlers[jobType] = func(ctx context.Context, payload []byte) (err error) {
		var input T

		err = json.Unmarshal(payload, &input)
		// jsonDecoder.DisallowUnknownFields()
		if err != nil {
			err = fmt.Errorf("workerpool: error decoding job data: %w", err)
			return
		}
		return handler(ctx, input)
	}
}

func (workerPool *WorkerPool) Start(ctx context.Context) {
	jobsChan := make(chan queue.Job, workerPool.concurrencyMax)
	var wg sync.WaitGroup

	wg.Add(int(workerPool.concurrencyMax))

	// Start the background workers
	for i := uint32(0); i < workerPool.concurrencyMax; i += 1 {
		go func(ctx context.Context, jobs <-chan queue.Job) {
			defer wg.Done()
			for job := range jobs {
				workerPool.handleJob(ctx, job)
			}
		}(ctx, jobsChan)
	}

	workerPool.logger.Info("workerpool: Starting", slog.Uint64("concurrencyMax", uint64(workerPool.concurrencyMax)))

	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			workerPool.logger.Info("workerpool: Shutting down")
			close(jobsChan)
			// workerPool.queue.Stop(ctx)
			wg.Wait()
			return
		case <-ticker.C:
			// we "sleep" a little bit to avoid consumming too much CPU
		}

		jobs, err := workerPool.queue.Pull(ctx, uint64(workerPool.concurrencyMax))
		if err != nil {
			workerPool.logger.Error("workerpool: error pulling jobs from queue", slog.String("err", err.Error()))
			time.Sleep(100 * time.Millisecond)
			continue
		}
		_ = jobs

		for _, job := range jobs {
			jobsChan <- job
		}
	}
}

func (workerPool *WorkerPool) handleJob(ctx context.Context, job queue.Job) {
	var err error

	jobHandler, jobHandlerExists := workerPool.jobHandlers[job.Type]
	if !jobHandlerExists {
		workerPool.logger.Error("workerpool: job handler not found", slog.String("job.type", job.Type),
			slog.String("job.id", job.ID.String()))
		goto failjob
	}

	err = jobHandler(ctx, job.RawData)
	if err != nil {
		goto failjob
	}
	// else {
	// 	logger.Info("workerpool: job successfully executd", slog.Duration("duration", duration), slog.Group("job",
	// 		slog.String("type", job.Type), slog.String("id", job.ID.String()),
	// 	))
	// }

	err = workerPool.queue.DeleteJob(ctx, job.ID)
	if err != nil {
		workerPool.logger.Error("workerpool: error deleting job", slog.String("job.id", job.ID.String()),
			slog.String("err", err.Error()))
	}
	return

failjob:
	workerPool.logger.Error("workerpool: job failed", slog.Group("job",
		slog.String("job.id", job.ID.String()), slog.String("type", job.Type),
	), slogx.Err(err))
	err = workerPool.queue.FailJob(ctx, job)
	if err != nil {
		workerPool.logger.Error("workerpool: error marking job as failed", slog.String("job.id", job.ID.String()),
			slogx.Err(err))
		return
	}
}
