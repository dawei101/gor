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

func (l *Limiter) Watch(f OnLimiter) {
	if f != nil {
		go func() {
			for {
				f(l.c)
			}
		}()
	}
}
