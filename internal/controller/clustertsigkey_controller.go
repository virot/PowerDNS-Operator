/*
 * Software Name : PowerDNS-Operator
 *
 * SPDX-FileCopyrightText: Copyright (c) PowerDNS-Operator contributors
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange Business Services SA
 * SPDX-License-Identifier: Apache-2.0
 *
 * This software is distributed under the Apache 2.0 License,
 * see the "LICENSE" file for more details
 */

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	dnsv1alpha2 "github.com/powerdns-operator/powerdns-operator/api/v1alpha2"
)

// ClusterTSIGKeyReconciler reconciles a ClusterTSIGKey object
type ClusterTSIGKeyReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	PDNSClient PdnsClienter
}

//+kubebuilder:rbac:groups=dns.cav.enablers.ob,resources=clustertsigkeys,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dns.cav.enablers.ob,resources=clustertsigkeys/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dns.cav.enablers.ob,resources=clustertsigkeys/finalizers,verbs=update

func (r *ClusterTSIGKeyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconcile ClusterTSIGKey", "ClusterTSIGKey.Name", req.Name)

	// Get ClusterTSIGKey
	tsigKey := &dnsv1alpha2.ClusterTSIGKey{}
	err := r.Get(ctx, req.NamespacedName, tsigKey)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Initialize variable to represent ClusterTSIGKey situation
	isModified := tsigKey.Status.ObservedGeneration != nil && *tsigKey.Status.ObservedGeneration != tsigKey.GetGeneration()
	isDeleted := !tsigKey.DeletionTimestamp.IsZero()

	// When updating a ClusterTSIGKey, if 'Status' is not changed, 'LastTransitionTime' will not be updated
	// So we delete condition to force new 'LastTransitionTime'
	original := tsigKey.DeepCopy()
	if !isDeleted && isModified {
		meta.RemoveStatusCondition(&tsigKey.Status.Conditions, "Available")
		if err := r.Status().Patch(ctx, tsigKey, client.MergeFrom(original)); err != nil {
			log.Error(err, "unable to patch ClusterTSIGKey status")
			return ctrl.Result{}, err
		}
	}

	return tsigKeyReconcile(ctx, tsigKey, isModified, isDeleted, r.Client, r.PDNSClient, log)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterTSIGKeyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dnsv1alpha2.ClusterTSIGKey{}).
		Complete(r)
}
