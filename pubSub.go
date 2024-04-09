package main

type Pubsub struct {
	Listners *[]chan<- string
}

func Newpubsub() Pubsub {
	return Pubsub{
		Listners: &[]chan<- string{},
	}
}

func (p Pubsub) subscribe() <-chan string {
	ch := make(chan string)
	*p.Listners = append(*p.Listners, ch)

	return ch
}

func (p Pubsub) notifyAll(msg string) {

	for _, ch := range *p.Listners {
		ch <- msg
	}

}

func (p Pubsub) close() {
	for _, ch := range *p.Listners {
		close(ch)
	}
}
