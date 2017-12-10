package processor

import (
	"errors"
	"fmt"
	"rabbitmqwebapp/messagebroker"
	"rabbitmqwebapp/models"
	"sync"
	"time"

	"github.com/lucsky/cuid"
)

type IProcessor interface {
	Process(phoneNumber string) (models.PortResponse, error)
	SendPortResponse(requestId string) bool
}

type processorImplementer struct {
	timeOut         time.Duration
	requestMutex    sync.Mutex
	messageBroker   messagebroker.IMessageBroker
	pendingRequests map[string]chan models.PortResponse
}

type messageListnerImplementer struct {
	processorImplementer *processorImplementer
}

func NewProcessor(timeOut time.Duration, messageBroker messagebroker.IMessageBroker) IProcessor {
	ret := &processorImplementer{
		timeOut:         timeOut,
		requestMutex:    sync.Mutex{},
		messageBroker:   messageBroker,
		pendingRequests: make(map[string]chan models.PortResponse),
	}

	return ret
}

func (p *processorImplementer) Process(phoneNumber string) (models.PortResponse, error) {
	requestId := cuid.New()
	println(fmt.Sprintf("REQUEST ID=== %s", requestId))

	responseChannel := p.registerRequest(requestId)

	timeToWait := time.NewTimer(p.timeOut)

	select {
	case response := <-responseChannel:
		return response, nil
	case <-timeToWait.C:
		errorResponse := models.PortResponse{
			RequestId:   requestId,
			PhoneNumber: phoneNumber,
			CanPort:     false,
			Message:     "Did not receive callback response in time",
		}
		return errorResponse, errors.New("Did not receive callback response in time")
	}
}

func (p *processorImplementer) SendPortResponse(requestId string) bool {
	result := models.PortResponse{
		RequestId:   requestId,
		PhoneNumber: "4153007086",
		Message:     "Phone number is Portable",
		CanPort:     true,
	}
	messageListner := messageListnerImplementer{
		processorImplementer: p,
	}
	err := p.messageBroker.DispatchMessage(requestId, result, messageListner)
	if err != nil {
		println(fmt.Sprintf("An error occurred dispatching message for requestId: %s to message broker", requestId))
		return false
	}

	return true
}

func (listner messageListnerImplementer) MessageConsumed(result models.PortResponse) {
	println(fmt.Sprintf("YAYYYYY! Message listener fired for requestId: %s", result.RequestId))
	responseChannel, hadChannel := listner.processorImplementer.pendingRequests[result.RequestId]

	if hadChannel {
		println("PUSHED MESSAGE ON TO THE CHANNEL")
		responseChannel <- result
	} else {
		println(fmt.Sprintf("CHANNEL DOES NOT EXIST FOR REQUESTID===%s", result.RequestId))
	}
}

func (p *processorImplementer) registerRequest(requestId string) chan models.PortResponse {
	p.requestMutex.Lock()
	defer p.requestMutex.Unlock()

	responseChannel := make(chan models.PortResponse)
	p.pendingRequests[requestId] = responseChannel

	return responseChannel
}
