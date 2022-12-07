/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

/*
Copyright 2022 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/

import (
	"context"

	"cloud.example.com/annotation-operator/controllers/reconciliation"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/codes"

	"cloud.example.com/annotation-operator/controllers/depresolver"
	"cloud.example.com/annotation-operator/controllers/providers/dns"
	"cloud.example.com/annotation-operator/controllers/providers/metrics"
	"go.opentelemetry.io/otel/trace"

	"cloud.example.com/annotation-operator/controllers/utils"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AnnoReconciler reconciles a Anno object
type AnnoReconciler struct {
	client.Client
	Scheme           *runtime.Scheme
	Config           *depresolver.Config
	DepResolver      depresolver.GslbResolver
	DNSProvider      dns.Provider
	Tracer           trace.Tracer
	IngressMapper    reconciliation.Mapper
	ReconcilerResult *utils.ReconcileResultHandler
	Log              *zerolog.Logger
	Metrics          *metrics.PrometheusMetrics
}

func (r *AnnoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx, span := r.Tracer.Start(ctx, "Reconcile")
	defer span.End()

	// == handle request
	if req.NamespacedName.Name == "" || req.NamespacedName.Namespace == "" {
		return r.ReconcilerResult.Requeue()
	}

	rs, rr, err := r.IngressMapper.Get(req.NamespacedName)
	switch rr {
	case reconciliation.MapperResultNotFound, reconciliation.MapperResultExistsButNotAnnotationFound:
		r.Log.Info().
			Str("Namespace", req.NamespacedName.Namespace).
			Str("Ingress", req.NamespacedName.Name).
			Str("Annotation", reconciliation.AnnotationStrategy).
			Msg("Ingress or annotation not found. Stop...")
		return r.ReconcilerResult.Stop()
	case reconciliation.MapperResultError:
		r.Metrics.IncrementError(req.NamespacedName)
		r.Log.Err(err).
			Str("Namespace", req.NamespacedName.Namespace).
			Str("Ingress", req.NamespacedName.Name).
			Msg("reading Ingress error")
		return r.ReconcilerResult.Requeue()
	}

	// == handle finalizers
	if r.DNSProvider.RequireFinalizer() {
		_, fSpan := r.Tracer.Start(ctx, "Handle finalizer")
		result, err := r.handleFinalizer(rs)
		switch result {
		case reconciliation.MapperFinalizerSkipped:
			fSpan.End()
		case reconciliation.MapperFinalizerInstalled:
			r.Log.Info().
				Str("finalizer", reconciliation.Finalizer).
				Msg("Injected finalizer")
			fSpan.End()
		case reconciliation.MapperFinalizerRemoved:
			r.Log.Info().
				Str("finalizer", reconciliation.Finalizer).
				Msg("Remove injected finalizer")
			fSpan.End()
			return r.ReconcilerResult.Stop()
		case reconciliation.MapperResultError:
			r.Log.Warn().
				Str("finalizer", reconciliation.Finalizer).
				AnErr("error", err).
				Msg("Injecting finalizer error")
			fSpan.RecordError(err)
			fSpan.SetStatus(codes.Error, err.Error())
			fSpan.End()
			r.Metrics.IncrementError(rs.NamespacedName)
			return r.ReconcilerResult.RequeueError(err)
		}
	}

	r.Log.Info().
		Str("EdgeDNSZone", r.Config.DNSZone).
		Msg("* Starting Reconciliation")

	// == external-dns dnsendpoints CRs ==
	dnsEndpoint, err := r.gslbDNSEndpoint(rs)
	if err != nil {
		r.Metrics.IncrementError(rs.NamespacedName)
		return r.ReconcilerResult.RequeueError(err)
	}

	_, s := r.Tracer.Start(ctx, "SaveDNSEndpoint")
	err = r.DNSProvider.SaveDNSEndpoint(rs, dnsEndpoint)
	if err != nil {
		r.Metrics.IncrementError(rs.NamespacedName)
		return r.ReconcilerResult.RequeueError(err)
	}
	s.End()

	// == handle delegated zone in Edge DNS
	_, szd := r.Tracer.Start(ctx, "CreateZoneDelegationForExternalDNS")
	err = r.DNSProvider.CreateZoneDelegationForExternalDNS(rs)
	if err != nil {
		r.Log.Err(err).Msg("Unable to create zone delegation")
		r.Metrics.IncrementError(rs.NamespacedName)
		return r.ReconcilerResult.Requeue()
	}
	szd.End()

	// == Status =
	err = r.updateStatus(rs, dnsEndpoint)
	if err != nil {
		r.Metrics.IncrementError(rs.NamespacedName)
		return r.ReconcilerResult.RequeueError(err)
	}
	// == Finish ==========
	// Everything went fine, requeue after RECONCILE_REQUEUE_SECONDS
	r.Metrics.IncrementReconciliation(rs.NamespacedName)
	return r.ReconcilerResult.Requeue()
}

func (r *AnnoReconciler) handleFinalizer(rs *reconciliation.LoopState) (reconciliation.MapperResult, error) {

	// Inject finalizer if doesn't exists
	result, err := r.IngressMapper.TryInjectFinalizer(rs)
	if result.IsIn(reconciliation.MapperFinalizerInstalled, reconciliation.MapperResultError) {
		return result, err
	}

	// Try remove if isMarkedToBeDeleted
	result, err = r.IngressMapper.TryRemoveFinalizer(rs, r.DNSProvider.Finalize)
	if result.IsIn(reconciliation.MapperFinalizerRemoved, reconciliation.MapperResultError) {
		return result, err
	}

	return reconciliation.MapperFinalizerSkipped, nil
}
