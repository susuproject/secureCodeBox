// SPDX-FileCopyrightText: 2021 iteratec GmbH
//
// SPDX-License-Identifier: Apache-2.0

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CascadeSpec describes how and when cascading scans should be generated.
type CascadeSpec struct {
	// InheritLabels defines whether cascading scans should inherit labels from the parent scan
	// +optional
	// +kubebuilder:default=true
	InheritLabels bool `json:"inheritLabels"`

	// InheritAnnotations defines whether cascading scans should inherit annotations from the parent scan
	// +optional
	// +kubebuilder:default=true
	InheritAnnotations bool `json:"inheritAnnotations"`

	// InheritEnv defines whether cascading scans should inherit environment variables from the parent scan
	// +optional
	// +kubebuilder:default=false
	InheritEnv bool `json:"inheritEnv"`

	// InheritVolumes defines whether cascading scans should inherit volumes and volume mounts from the parent scan
	// +optional
	// +kubebuilder:default=false
	InheritVolumes bool `json:"inheritVolumes"`

	// InheritInitContainers defines whether cascading scans should inherit initContainers from the parent scan. Usually only useful when combined with inheritVolumes.
	// +optional
	// +kubebuilder:default=false
	InheritInitContainers bool `json:"inheritInitContainers"`

	// InheritHookSelector defines whether cascading scans should inherit hookSelector from the parent scan.
	// +optional
	// +kubebuilder:default=false
	InheritHookSelector bool `json:"inheritHookSelector"`

	// matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels
	// map is equivalent to an element of matchExpressions, whose key field is "key", the
	// operator is "In", and the values array contains only "value". The requirements are ANDed.
	// +optional
	MatchLabels map[string]string `json:"matchLabels,omitempty" protobuf:"bytes,1,rep,name=matchLabels"`
	// matchExpressions is a list of label selector requirements. The requirements are ANDed.
	// +optional
	MatchExpressions []metav1.LabelSelectorRequirement `json:"matchExpressions,omitempty" protobuf:"bytes,2,rep,name=matchExpressions"`
}

// ScanSpec defines the desired state of Scan
type ScanSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The name of the scanType which should be started.
	// +kubebuilder:validation:Required
	ScanType string `json:"scanType,omitempty"`

	// All CLI parameters to configure the scan container.
	// +kubebuilder:validation:Required
	Parameters []string `json:"parameters,omitempty"`

	// HookSelector allows to specify a LabelSelector with which the hooks are selected.
	HookSelector *metav1.LabelSelector `json:"hookSelector,omitempty"`

	// Env allows to specify environment vars for the scanner container. These will be merged will the env vars specified for the first container of the pod defined in the ScanType
	Env []corev1.EnvVar `json:"env,omitempty"`
	// Volumes allows to specify volumes for the scan container.
	Volumes []corev1.Volume `json:"volumes,omitempty"`
	// VolumeMounts allows to specify volume mounts for the scan container.
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`
	// InitContainers allows to specify init containers for the scan container, to pre-load data into them.
	InitContainers []corev1.Container `json:"initContainers,omitempty"`

	Cascades *CascadeSpec `json:"cascades,omitempty"`
}

// ScanStatus defines the observed state of Scan
type ScanStatus struct {
	State string `json:"state,omitempty"`

	// FinishedAt contains the time where the scan (including parser & hooks) has been marked as "Done"
	FinishedAt       *metav1.Time `json:"finishedAt,omitempty"`
	ErrorDescription string       `json:"errorDescription,omitempty"`

	// RawResultType determines which kind of ParseDefinition will be used to turn the raw results of the scanner into findings
	RawResultType string `json:"rawResultType,omitempty"`
	// RawResultFile Filename of the result file of the scanner. e.g. `nmap-result.xml`
	RawResultFile string `json:"rawResultFile,omitempty"`

	// FindingDownloadLink link to download the finding json file from. Valid for 7 days
	FindingDownloadLink string `json:"findingDownloadLink,omitempty"`
	// RawResultDownloadLink link to download the raw result file from. Valid for 7 days
	RawResultDownloadLink string `json:"rawResultDownloadLink,omitempty"`

	// FindingHeadLink link to send HEAD request to the finding json file. Valid for 7 days
	FindingHeadLink string `json:"findingHeadLink,omitempty"`
	// RawResultHeadLink link to send HEAD request to raw result file. Valid for 7 days
	RawResultHeadLink string `json:"rawResultHeadLink,omitempty"`

	Findings FindingStats `json:"findings,omitempty"`

	ReadAndWriteHookStatus []HookStatus `json:"readAndWriteHookStatus,omitempty"`

	OrderedHookStatuses [][]*HookStatus `json:"orderedHookStatuses,omitempty"`
}

// HookState Describes the State of a Hook on a Scan
type HookState string

const (
	Pending    HookState = "Pending"
	InProgress HookState = "InProgress"
	Completed  HookState = "Completed"
	Cancelled  HookState = "Cancelled"
	Failed     HookState = "Failed"
)

type HookStatus struct {
	HookName string    `json:"hookName"`
	State    HookState `json:"state"`
	JobName  string    `json:"jobName,omitempty"`
	Priority int       `json:"priority"`
	Type     HookType  `json:"type"`
}

// FindingStats contains the general stats about the results of the scan
type FindingStats struct {
	// Count indicates how many findings were identified in total
	Count uint64 `json:"count,omitempty"`
	// FindingSeverities indicates the count of finding with the respective severity
	FindingSeverities FindingSeverities `json:"severities,omitempty"`
	// FindingCategories indicates the count of finding broken down by their categories
	FindingCategories map[string]uint64 `json:"categories,omitempty"`
}

// FindingSeverities indicates the count of finding with the respective severity
type FindingSeverities struct {
	Informational uint64 `json:"informational,omitempty"`
	Low           uint64 `json:"low,omitempty"`
	Medium        uint64 `json:"medium,omitempty"`
	High          uint64 `json:"high,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="UID",type=string,JSONPath=`.metadata.uid`,description="K8s Resource UID",priority=1
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.scanType`,description="Scan Type"
// +kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`,description="Scan State"
// +kubebuilder:printcolumn:name="Findings",type=string,JSONPath=`.status.findings.count`,description="Total Finding Count"
// +kubebuilder:printcolumn:name="Parameters",type=string,JSONPath=`.spec.parameters`,description="Arguments passed to the Scanner",priority=1

// Scan is the Schema for the scans API
type Scan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScanSpec   `json:"spec,omitempty"`
	Status ScanStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ScanList type wrapping multiple Scans
type ScanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Scan `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Scan{}, &ScanList{})
}
