package utils

type Pool struct {
}

func NewPool[T any](n int) *Pool {
	var workers = make([](chan T), n)
	for i := 0; i < n; i++ {
		workers[i] = make(chan T, 1)
	}
	return &Pool{}
}

func (p *Pool) Run() {
	
}

func (p *Pool) Do(f func()) {
	go func() {
		f()
	}()

}
