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

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// +kubebuilder:printcolumn:name="Algorithm",type="string",JSONPath=".spec.algorithm"
// +kubebuilder:printcolumn:name="ID",type="string",JSONPath=".status.id"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.syncStatus"
// ClusterTSIGKey is the Schema for the clustertsigkeys API
type ClusterTSIGKey struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TSIGKeySpec   `json:"spec,omitempty"`
	Status TSIGKeyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClusterTSIGKeyList contains a list of ClusterTSIGKey
type ClusterTSIGKeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterTSIGKey `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterTSIGKey{}, &ClusterTSIGKeyList{})
}

// IsInExpectedStatus returns true if Status.SyncStatus and Status.ObservedGeneration are, at least, at expected value
func (t *ClusterTSIGKey) IsInExpectedStatus(expectedMinimumObservedGeneration int64, expectedSyncStatus string) bool {
	return t.Status.ObservedGeneration != nil &&
		*t.Status.ObservedGeneration >= expectedMinimumObservedGeneration &&
		t.Status.SyncStatus != nil &&
		*t.Status.SyncStatus == expectedSyncStatus
}
