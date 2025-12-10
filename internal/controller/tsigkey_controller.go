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

const (
	TSIGKeyReasonSynced                = "TSIGKeySynced"
	TSIGKeyMessageSyncSucceeded        = "TSIG key synced with PowerDNS instance"
	TSIGKeyReasonSynchronizationFailed = "SynchronizationFailed"
	TSIGKeyReasonSecretNotFound        = "SecretNotFound"
	TSIGKeyMessageSecretNotFound       = "Referenced secret not found"
)

// TSIGKeyReconciler reconciles a TSIGKey object
type TSIGKeyReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	PDNSClient PdnsClienter
}

//+kubebuilder:rbac:groups=dns.cav.enablers.ob,resources=tsigkeys,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dns.cav.enablers.ob,resources=tsigkeys/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dns.cav.enablers.ob,resources=tsigkeys/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

func (r *TSIGKeyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconcile TSIGKey", "TSIGKey.Name", req.Name)

	// Get TSIGKey
	tsigKey := &dnsv1alpha2.TSIGKey{}
	err := r.Get(ctx, req.NamespacedName, tsigKey)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Initialize variable to represent TSIGKey situation
	isModified := tsigKey.Status.ObservedGeneration != nil && *tsigKey.Status.ObservedGeneration != tsigKey.GetGeneration()
	isDeleted := !tsigKey.DeletionTimestamp.IsZero()

	// When updating a TSIGKey, if 'Status' is not changed, 'LastTransitionTime' will not be updated
	// So we delete condition to force new 'LastTransitionTime'
	original := tsigKey.DeepCopy()
	if !isDeleted && isModified {
		meta.RemoveStatusCondition(&tsigKey.Status.Conditions, "Available")
		if err := r.Status().Patch(ctx, tsigKey, client.MergeFrom(original)); err != nil {
			log.Error(err, "unable to patch TSIGKey status")
			return ctrl.Result{}, err
		}
	}

	return tsigKeyReconcile(ctx, tsigKey, isModified, isDeleted, r.Client, r.PDNSClient, log)
}

// SetupWithManager sets up the controller with the Manager.
func (r *TSIGKeyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dnsv1alpha2.TSIGKey{}).
		Complete(r)
}
