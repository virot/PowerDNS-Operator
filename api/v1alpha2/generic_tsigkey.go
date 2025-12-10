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

//nolint:dupl
package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// +kubebuilder:object:root=false
// +kubebuilder:object:generate:false
// +k8s:deepcopy-gen:interfaces=nil
// +k8s:deepcopy-gen=nil

// GenericTSIGKey is a common interface for interacting with ClusterTSIGKey
// or a namespaced TSIGKey.
type GenericTSIGKey interface {
	runtime.Object
	metav1.Object

	GetObjectMeta() *metav1.ObjectMeta
	GetTypeMeta() *metav1.TypeMeta

	GetSpec() *TSIGKeySpec
	GetStatus() TSIGKeyStatus
	SetStatus(status TSIGKeyStatus)
	Copy() GenericTSIGKey
}

// +kubebuilder:object:root:false
// +kubebuilder:object:generate:false
var _ GenericTSIGKey = &TSIGKey{}

func (t *TSIGKey) GetObjectMeta() *metav1.ObjectMeta {
	return &t.ObjectMeta
}

func (t *TSIGKey) GetTypeMeta() *metav1.TypeMeta {
	return &t.TypeMeta
}

func (t *TSIGKey) GetSpec() *TSIGKeySpec {
	return &t.Spec
}

func (t *TSIGKey) GetStatus() TSIGKeyStatus {
	return t.Status
}

func (t *TSIGKey) SetStatus(status TSIGKeyStatus) {
	t.Status = status
}

func (t *TSIGKey) Copy() GenericTSIGKey {
	return t.DeepCopy()
}

// +kubebuilder:object:root:false
// +kubebuilder:object:generate:false
var _ GenericTSIGKey = &ClusterTSIGKey{}

func (t *ClusterTSIGKey) GetObjectMeta() *metav1.ObjectMeta {
	return &t.ObjectMeta
}

func (t *ClusterTSIGKey) GetTypeMeta() *metav1.TypeMeta {
	return &t.TypeMeta
}

func (t *ClusterTSIGKey) GetSpec() *TSIGKeySpec {
	return &t.Spec
}

func (t *ClusterTSIGKey) GetStatus() TSIGKeyStatus {
	return t.Status
}

func (t *ClusterTSIGKey) SetStatus(status TSIGKeyStatus) {
	t.Status = status
}

func (t *ClusterTSIGKey) Copy() GenericTSIGKey {
	return t.DeepCopy()
}
