package reactsocket

import (
	"goreact"
	"log"
	"net/http"
	"wsocket"
)

func ConnectSocket(addr string) (*wsocket.Socket, *goreact.Signal, error) {
	cl, err := wsocket.ConnectSocket(addr)
	if err != nil {
		return nil, nil, err
	}
	return cl, incomingEvent(cl), nil
}

func AcceptSocket(w http.ResponseWriter, r *http.Request) (*wsocket.Socket, *goreact.Signal, error) {
	cl, err := wsocket.AcceptNewSocket(w, r)
	if err != nil {
		return nil, nil, err
	}
	return cl, incomingEvent(cl), nil
}

func incomingEvent(cl *wsocket.Socket) *goreact.Signal {

	return goreact.NewSignal(func(i goreact.Injector) {
	L:
		for {
			select {
			case event, ok := <-cl.GetRead():
				if ok == false {
					log.Println("completed")
					i.SendCompleted()
					break L
				}
				log.Println("Send Event")
				i.SendNext(event)
			case err, ok := <-cl.GetError():
				if ok == false {
					log.Println("Send Completed")
					i.SendCompleted()
					break L
				} else {
					log.Println("Send failed", err)
					i.SendFailed(err)
					return
				}
			}
		}
	})
}
