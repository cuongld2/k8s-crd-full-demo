package api

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type DatabaseSpec struct {
	DbName      string `json:"dbName"`
	Description string `json:"description,omitempty"`
	Total       int    `json:"total"`
	Available   int    `json:"available"`
	DbType      string `json:"dbType"`
	Tags        string `json:"tags,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Database struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec DatabaseSpec `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Database `json:"items"`
}
