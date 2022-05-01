package utils

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"time"
)

type Event struct {
	Object  runtime.Object
	Name    string
	Reason  string
	Message string
}

type ErrorEvent struct {
	Event
	Err error
}

func GetRequeueResult() ctrl.Result {
	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: time.Duration(time.Duration.Minutes(1)),
	}
}

func RecordError(recorder record.EventRecorder, errorEvent ErrorEvent) {
	recorder.Event(errorEvent.Object, "Warning", errorEvent.Reason,
		fmt.Sprintf("%s for %s: %s", errorEvent.Message, errorEvent.Name, errorEvent.Err.Error()))
}

func RecordSuccess(recorder record.EventRecorder, event Event) {
	message := fmt.Sprintf("%s successful for %s", event.Reason, event.Name)
	if event.Message != "" {
		message = fmt.Sprintf("%s for %s", event.Message, event.Name)
	}

	recorder.Event(event.Object, "Normal", event.Reason, message)
}

func RecordEventAndReturn(res ctrl.Result, err error, recorder record.EventRecorder, event Event) (ctrl.Result, error) {

	if err != nil {
		RecordError(recorder, ErrorEvent{
			Event: event,
			Err:   err,
		})
	} else {
		RecordSuccess(recorder, event)
	}

	return res, err
}
