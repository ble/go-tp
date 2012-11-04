package drawing

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type drawRequest struct {
	DrawPart
	response chan error
}

type getRequest chan []byte

type DrawingHandle interface {
	Draw(DrawPart) error
	Read() ([]byte, error)
	Close() error
}

type drawingAgent struct {
	drawChannel chan drawRequest
	getChannel  chan getRequest
}

func (d drawingAgent) Draw(p DrawPart) error {
	request := drawRequest{p, make(chan error)}
	d.drawChannel <- request
	timeout := time.After(time.Second)
	select {
	case err := <-request.response:
		return err
	case _ = <-timeout:
		return errors.New("timeout waiting for reply")
	}
	return errors.New("fallthrough")
}

func (d drawingAgent) Read() ([]byte, error) {
	request := make(getRequest)
	d.getChannel <- request
	timeout := time.After(time.Second)
	select {
	case data, ok := <-request:
		if ok {
			return data, nil
		} else {
			return data, errors.New("response channel closed")
		}
	case _ = <-timeout:
		return nil, errors.New("timeout waiting for reply")
	}
	return nil, errors.New("fallthrough")
}

func (d drawingAgent) Close() error {
	close(d.drawChannel)
	close(d.getChannel)
	return nil
}

func (d drawingAgent) runAgent(drawing []DrawPart) {
	anyClosed := false
	for !anyClosed {
		select {
		case draw, isValue := <-d.drawChannel:
			if !isValue {
				anyClosed = true
				continue
			}
			//TODO pre-processing on the draw part
			//relevant to existing drawing here...
			drawing = append(drawing, draw.DrawPart)
			fmt.Println("draw requested")
			draw.response <- nil
		case get, isValue := <-d.getChannel:
			if !isValue {
				anyClosed = true
				continue
			}
			fmt.Println("get requested")
			bytes, err := json.Marshal(drawing)
			if err == nil {
				get <- bytes
			} else {
				close(get)
			}
		}
	}
	fmt.Println("done")
}

func NewDrawingHandle() DrawingHandle {
	drawing := make([]DrawPart, 0, 0)
	draws := make(chan drawRequest)
	gets := make(chan getRequest)

	agent := drawingAgent{draws, gets}
	go agent.runAgent(drawing)
	return agent
}
