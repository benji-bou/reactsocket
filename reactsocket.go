package reactsocket

import (
	"context"
	"net/http"

	"github.com/benji-bou/wsocket"
	"github.com/reactivex/rxgo/v2"
)

func ConnectSocket(addr string, header http.Header) (*wsocket.Socket, rxgo.Observable, error) {
	cl, err := wsocket.ConnectSocket(addr, header)
	if err != nil {
		return nil, nil, err
	}
	obs := incomingEvent(cl)
	obs.Connect(context.Background())
	return cl, obs, nil
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
					break L
				}
				next <- rxgo.Of(event)
			case err, ok := <-cl.Error():
				if ok == false {
					break L
				} else {
					next <- rxgo.Error(err)
					break L
				}

			}
		}
	}}, rxgo.WithPublishStrategy(), rxgo.WithObservationStrategy(rxgo.Eager))
}
