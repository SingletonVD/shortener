package concurrency

type Semaphore struct {
	semaCh chan struct{}
}

func NewSemaphore(active int) *Semaphore {
	return &Semaphore{
		semaCh: make(chan struct{}, active),
	}
}

func (s *Semaphore) Acquire() {
	s.semaCh <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.semaCh
}
