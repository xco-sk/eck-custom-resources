package utils

import (
	ctrl "sigs.k8s.io/controller-runtime"
	"time"
)

func GetRequeueResult() ctrl.Result {
	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: time.Duration(time.Duration.Minutes(1)),
	}
}
