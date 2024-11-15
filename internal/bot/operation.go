package bot

import "time"

type Operation struct {
	chanID    int64
	operation string

	messageChanel chan string
}

func NewOperation(id int64, op string) *Operation {
	return &Operation{chanID: id, operation: op, messageChanel: make(chan string)}
}

func (o *Operation) Process() {

	timer := time.NewTimer(3 * time.Minute)

	for {
		select {
		case msg := <-o.messageChanel:
		case <-timer.C:
			//no message for 3 minutes
			break
		}
	}
}

func (o *Operation) DeliverMessage(msg string) {

}
