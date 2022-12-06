package reconciliation

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
	"fmt"
	"reflect"

	"cloud.example.com/annotation-operator/controllers/utils"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// todo: rename package to reconciliation

type MapperResult int

const (
	MapperResultExists MapperResult = 1 << iota
	MapperResultNotFound
	MapperResultError
	MapperResultExistsButNotAnnotationFound
	MapperFinalizerRemoved
	MapperFinalizerInstalled
	MapperFinalizerSkipped
)

// IngressMapper provides API for working with ingress
type IngressMapper struct {
	c client.Client
}

func NewIngressMapper(c client.Client) *IngressMapper {
	return &IngressMapper{
		c: c,
	}
}

func (i *IngressMapper) UpdateStatus(state *LoopState) (err error) {
	// check if object has not been deleted
	var r MapperResult
	var s *LoopState
	s, r, err = i.Get(state.NamespacedName)
	switch r {
	case MapperResultError:
		return err
	case MapperResultNotFound:
		// object was deleted
		return nil
	}
	// update the planned object
	s.Ingress.Annotations[AnnotationStatus] = state.Status.String()
	return i.c.Update(context.TODO(), s.Ingress)
}

func (i *IngressMapper) Get(selector types.NamespacedName) (rs *LoopState, result MapperResult, err error) {
	var ing = &netv1.Ingress{}
	err = i.c.Get(context.TODO(), selector, ing)
	result, err = i.getConverterResult(err, ing)
	if result == MapperResultError {
		return nil, result, err
	}
	rs, err = NewLoopState(ing)
	if err != nil {
		result = MapperResultError
	}
	return rs, result, err
}

// Equal compares given ingress annotations and Ingres.Spec. If any of ingresses doesn't exist, returns false
func (i *IngressMapper) Equal(rs1 *LoopState, rs2 *LoopState) bool {
	if rs1 == nil || rs2 == nil {
		return false
	}
	if !reflect.DeepEqual(rs1.Spec, rs2.Spec) {
		return false
	}
	if !reflect.DeepEqual(rs1.Ingress.Spec, rs2.Ingress.Spec) {
		return false
	}
	return true
}

func (i *IngressMapper) TryInjectFinalizer(rs *LoopState) (MapperResult, error) {
	if rs == nil || rs.Ingress == nil {
		return MapperResultError, fmt.Errorf("injecting finalizer from nil values")
	}
	if !utils.Contains(rs.Ingress.GetFinalizers(), Finalizer) {
		rs.Ingress.SetFinalizers(append(rs.Ingress.GetFinalizers(), Finalizer))
		err := i.c.Update(context.TODO(), rs.Ingress)
		if err != nil {
			return MapperResultError, err
		}
		return MapperFinalizerInstalled, nil
	}
	return MapperFinalizerSkipped, nil
}

func (i *IngressMapper) TryRemoveFinalizer(rs *LoopState, finalize func(*LoopState) error) (MapperResult, error) {
	if rs == nil || rs.Ingress == nil {
		return MapperResultError, fmt.Errorf("removing finalizer from nil values")
	}
	if utils.Contains(rs.Ingress.GetFinalizers(), Finalizer) {
		isMarkedToBeDeleted := rs.Ingress.GetDeletionTimestamp() != nil
		if !isMarkedToBeDeleted {
			return MapperFinalizerSkipped, nil
		}
		err := finalize(rs)
		if err != nil {
			return MapperResultError, err
		}
		rs.Ingress.SetFinalizers(utils.Remove(rs.Ingress.GetFinalizers(), Finalizer))
		err = i.c.Update(context.TODO(), rs.Ingress)
		if err != nil {
			return MapperResultError, err
		}
		return MapperFinalizerRemoved, nil
	}
	return MapperFinalizerSkipped, nil
}

func (i *IngressMapper) getConverterResult(err error, ing *netv1.Ingress) (MapperResult, error) {
	if err != nil && errors.IsNotFound(err) {
		return MapperResultNotFound, nil
	} else if err != nil {
		return MapperResultError, err
	}
	if _, found := ing.GetAnnotations()[AnnotationStrategy]; !found {
		return MapperResultExistsButNotAnnotationFound, nil
	}
	return MapperResultExists, nil
}
