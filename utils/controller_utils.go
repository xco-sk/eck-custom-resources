package utils

import (
	"context"
	"fmt"
	"time"

	configv2 "github.com/xco-sk/eck-custom-resources/apis/config/v2"
	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
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

func GetUserSecret(cli client.Client, ctx context.Context, namespace string, auth *configv2.UsernamePasswordAuthentication, secret *k8sv1.Secret) error {
	if err := cli.Get(ctx, client.ObjectKey{Namespace: namespace, Name: auth.SecretName}, secret); err != nil {
		return err
	}
	return nil
}

func GetCertificateSecret(cli client.Client, ctx context.Context, namespace string, certificate *configv2.PublicCertificate, secret *k8sv1.Secret) error {
	if err := cli.Get(ctx, client.ObjectKey{Namespace: namespace, Name: certificate.SecretName}, secret); err != nil {
		return err
	}
	return nil
}

func CommonEventFilter() predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
		},
	}
}
