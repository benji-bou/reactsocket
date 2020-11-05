package reactsocket

import (
	"context"
	"log"
	"net/http"

	"github.com/benji-bou/wsocket"
	"github.com/reactivex/rxgo/v2"
)

func ConnectSocket(addr string) (*wsocket.Socket, rxgo.Observable, error) {
	cl, err := wsocket.ConnectSocket(addr)
	if err != nil {
		return nil, nil, err
	}
	return cl, incomingEvent(cl), nil
}

func AcceptSocket(w http.ResponseWriter, r *http.Request) (*wsocket.Socket, rxgo.Observable, error) {
	cl, err := wsocket.AcceptNewSocket(w, r)
	if err != nil {
		return nil, nil, err
	}
	obs := incomingEvent(cl)
	obs.Connect(context.Background())
	return cl, obs, nil
}

func incomingEvent(cl *wsocket.Socket) rxgo.Observable {

	return rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
	L:
		for {
			select {
			case event, ok := <-cl.Read():
				if ok == false {
					log.Println("completed")
					break L
				}
				next <- rxgo.Of(event)
			case err, ok := <-cl.Error():
				if ok == false {
					log.Println("Send Completed")
					break L
				} else {
					log.Println("Send failed", err)
					next <- rxgo.Error(err)
					break L
				}
			}
		}
	}}, rxgo.WithPublishStrategy(), rxgo.WithBufferedChannel(1))
}
