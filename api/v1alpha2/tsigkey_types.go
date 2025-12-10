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

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TSIGKeySpec defines the desired state of TSIGKey
type TSIGKeySpec struct {
	// Algorithm for the TSIG key, e.g., "hmac-sha256", "hmac-sha512"
	// +kubebuilder:validation:Enum:=hmac-md5;hmac-sha1;hmac-sha224;hmac-sha256;hmac-sha384;hmac-sha512
	Algorithm string `json:"algorithm"`
	// SecretRef references a Kubernetes Secret containing the TSIG key value
	// The secret must contain a key named "key" with the base64-encoded secret
	// Either SecretRef or Key must be provided, but not both
	// +optional
	SecretRef *SecretReference `json:"secretRef,omitempty"`
	// Key contains the base64-encoded TSIG key value directly in the spec
	// WARNING: Storing sensitive data directly in the CRD is insecure as it will be visible
	// in plain text in the Kubernetes API and etcd. Use SecretRef instead for production environments.
	// Either SecretRef or Key must be provided, but not both
	// +optional
	Key *string `json:"key,omitempty"`
}

// SecretReference contains information to locate a Kubernetes Secret
type SecretReference struct {
	// Name of the secret in the same namespace (for TSIGKey) or specified namespace (for ClusterTSIGKey)
	Name string `json:"name"`
	// Namespace of the secret (only used for ClusterTSIGKey, ignored for TSIGKey)
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// TSIGKeyStatus defines the observed state of TSIGKey
type TSIGKeyStatus struct {
	// ID of the TSIG key in PowerDNS
	// +optional
	ID *string `json:"id,omitempty"`
	// Name of the TSIG key in PowerDNS
	// +optional
	Name *string `json:"name,omitempty"`
	// Algorithm used for the TSIG key
	// +optional
	Algorithm *string `json:"algorithm,omitempty"`
	// SyncStatus indicates the synchronization status
	// +optional
	SyncStatus *string `json:"syncStatus,omitempty"`
	// Conditions represent the latest available observations of the TSIGKey's state
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// ObservedGeneration reflects the generation of the most recently observed TSIGKey
	// +optional
	ObservedGeneration *int64 `json:"observedGeneration,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:storageversion
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Namespaced

// +kubebuilder:printcolumn:name="Algorithm",type="string",JSONPath=".spec.algorithm"
// +kubebuilder:printcolumn:name="ID",type="string",JSONPath=".status.id"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.syncStatus"
// TSIGKey is the Schema for the tsigkeys API
type TSIGKey struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TSIGKeySpec   `json:"spec,omitempty"`
	Status TSIGKeyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TSIGKeyList contains a list of TSIGKey
type TSIGKeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TSIGKey `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TSIGKey{}, &TSIGKeyList{})
}

// IsInExpectedStatus returns true if Status.SyncStatus and Status.ObservedGeneration are, at least, at expected value
func (t *TSIGKey) IsInExpectedStatus(expectedMinimumObservedGeneration int64, expectedSyncStatus string) bool {
	return t.Status.ObservedGeneration != nil &&
		*t.Status.ObservedGeneration >= expectedMinimumObservedGeneration &&
		t.Status.SyncStatus != nil &&
		*t.Status.SyncStatus == expectedSyncStatus
}
