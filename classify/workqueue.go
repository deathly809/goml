package classify

// Everything a thread needs to survivie
type workData struct {
	i int
	f func(i int)
}

// Runs until the jobs channel is closed
func thread(jobs <-chan workData) {
	for data := range jobs {
		data.f(data.i)
	}
}

// WorkQueue is a queue of work that is run on multiple threads
type WorkQueue struct {
	numThreads int
	jobs       chan workData
	started    bool
}

// Init the work queue
func (w *WorkQueue) Init(numThreads int) {
	w.jobs = make(chan workData, 100)
	w.numThreads = numThreads
}

// Start the work queue
func (w *WorkQueue) Start() {
	if !w.started {
		for i := 0; i < w.numThreads; i++ {
			go thread(w.jobs)
		}
	}
}

// Stop the work queue
func (w *WorkQueue) Stop() {
	close(w.jobs)
}

// Enqueue places work in the queue, blocks
func (w *WorkQueue) Enqueue(start, end int, f func(int)) {
	for i := start; i <= end; i++ {
		w.jobs <- workData{
			i: i,
			f: f,
		}
	}
}
