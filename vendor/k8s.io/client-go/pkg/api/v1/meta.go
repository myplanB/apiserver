/*
Copyright 2014 The Kubernetes Authors.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (obj *ObjectMeta) GetObjectMeta() metav1.Object { return obj }

// Namespace implements metav1.Object for any object with an ObjectMeta typed field. Allows
// fast, direct access to metadata fields for API objects.
func (meta *ObjectMeta) GetNamespace() string                { return meta.Namespace }
func (meta *ObjectMeta) SetNamespace(namespace string)       { meta.Namespace = namespace }
func (meta *ObjectMeta) GetName() string                     { return meta.Name }
func (meta *ObjectMeta) SetName(name string)                 { meta.Name = name }
func (meta *ObjectMeta) GetGenerateName() string             { return meta.GenerateName }
func (meta *ObjectMeta) SetGenerateName(generateName string) { meta.GenerateName = generateName }
func (meta *ObjectMeta) GetUID() types.UID                   { return meta.UID }
func (meta *ObjectMeta) SetUID(uid types.UID)                { meta.UID = uid }
func (meta *ObjectMeta) GetResourceVersion() string          { return meta.ResourceVersion }
func (meta *ObjectMeta) SetResourceVersion(version string)   { meta.ResourceVersion = version }
func (meta *ObjectMeta) GetSelfLink() string                 { return meta.SelfLink }
func (meta *ObjectMeta) SetSelfLink(selfLink string)         { meta.SelfLink = selfLink }
func (meta *ObjectMeta) GetCreationTimestamp() metav1.Time   { return meta.CreationTimestamp }
func (meta *ObjectMeta) SetCreationTimestamp(creationTimestamp metav1.Time) {
	meta.CreationTimestamp = creationTimestamp
}
func (meta *ObjectMeta) GetDeletionTimestamp() *metav1.Time { return meta.DeletionTimestamp }
func (meta *ObjectMeta) SetDeletionTimestamp(deletionTimestamp *metav1.Time) {
	meta.DeletionTimestamp = deletionTimestamp
}
func (meta *ObjectMeta) GetLabels() map[string]string                 { return meta.Labels }
func (meta *ObjectMeta) SetLabels(labels map[string]string)           { meta.Labels = labels }
func (meta *ObjectMeta) GetAnnotations() map[string]string            { return meta.Annotations }
func (meta *ObjectMeta) SetAnnotations(annotations map[string]string) { meta.Annotations = annotations }
func (meta *ObjectMeta) GetFinalizers() []string                      { return meta.Finalizers }
func (meta *ObjectMeta) SetFinalizers(finalizers []string)            { meta.Finalizers = finalizers }

func (meta *ObjectMeta) GetOwnerReferences() []metav1.OwnerReference {
	ret := make([]metav1.OwnerReference, len(meta.OwnerReferences))
	for i := 0; i < len(meta.OwnerReferences); i++ {
		ret[i].Kind = meta.OwnerReferences[i].Kind
		ret[i].Name = meta.OwnerReferences[i].Name
		ret[i].UID = meta.OwnerReferences[i].UID
		ret[i].APIVersion = meta.OwnerReferences[i].APIVersion
		if meta.OwnerReferences[i].Controller != nil {
			value := *meta.OwnerReferences[i].Controller
			ret[i].Controller = &value
		}
		if meta.OwnerReferences[i].BlockOwnerDeletion != nil {
			value := *meta.OwnerReferences[i].BlockOwnerDeletion
			ret[i].BlockOwnerDeletion = &value
		}
	}
	return ret
}

func (meta *ObjectMeta) SetOwnerReferences(references []metav1.OwnerReference) {
	newReferences := make([]metav1.OwnerReference, len(references))
	for i := 0; i < len(references); i++ {
		newReferences[i].Kind = references[i].Kind
		newReferences[i].Name = references[i].Name
		newReferences[i].UID = references[i].UID
		newReferences[i].APIVersion = references[i].APIVersion
		if references[i].Controller != nil {
			value := *references[i].Controller
			newReferences[i].Controller = &value
		}
		if references[i].BlockOwnerDeletion != nil {
			value := *references[i].BlockOwnerDeletion
			newReferences[i].BlockOwnerDeletion = &value
		}
	}
	meta.OwnerReferences = newReferences
}

func (meta *ObjectMeta) GetClusterName() string {
	return meta.ClusterName
}
func (meta *ObjectMeta) SetClusterName(clusterName string) {
	meta.ClusterName = clusterName
}
