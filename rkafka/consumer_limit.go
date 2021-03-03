package rkafka

type Limiter struct {
	n int
	c chan struct{}
}

func (l *Limiter) Run(f func()) {
	l.c <- struct{}{}
	go func() {
		f()
		<-l.c
	}()
}
