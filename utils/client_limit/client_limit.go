package client_limit

type ClientLimit struct {
	queue chan struct{}
}

func New(limit int) *ClientLimit {
	return &ClientLimit{
		queue: make(chan struct{}, limit),
	}
}

func (clientLimit *ClientLimit) Occupy() {
	clientLimit.queue <- struct{}{}
}

func (clientLimit *ClientLimit) Free() {
	<- clientLimit.queue
}
