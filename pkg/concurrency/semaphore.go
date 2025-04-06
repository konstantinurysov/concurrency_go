package concurrency

type Semaphore struct {
	semaphore chan struct{}
}

func NewSemaphore(concurrency int) Semaphore {
	return Semaphore{
		semaphore: make(chan struct{}, concurrency),
	}
}

func (s *Semaphore) Acquire() {
	if s == nil || s.semaphore == nil {
		return
	}

	s.semaphore <- struct{}{}
}

func (s *Semaphore) Release() {
	if s == nil || s.semaphore == nil {
		return
	}

	<-s.semaphore
}

func (s *Semaphore) WithAcquire(f func()) {
	if s == nil || s.semaphore == nil {
		return
	}

	s.Acquire()
	defer s.Release()
	f()
}
