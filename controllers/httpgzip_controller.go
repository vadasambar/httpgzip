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

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	_struct "github.com/golang/protobuf/ptypes/struct"
	appsv1alpha1 "github.com/vadasambar/httpgzip/api/v1alpha1"
	structpb "google.golang.org/protobuf/types/known/structpb"
	typesv1alpha3 "istio.io/api/networking/v1alpha3"
	clientv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HttpGzipReconciler reconciles a HttpGzip object
type HttpGzipReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.vadasambar.com,resources=httpgzips,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.vadasambar.com,resources=httpgzips/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.vadasambar.com,resources=httpgzips/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HttpGzip object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *HttpGzipReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	var hg appsv1alpha1.HttpGzip
	err := r.Client.Get(ctx, req.NamespacedName, &hg)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Log.Info("unable to fetch httpgzip", "name", req.Name, "namespace", req.Namespace)
			return ctrl.Result{}, nil
		}
		log.Log.Error(err, "unable to fetch httpgzip", "name", req.Name, "namespace", req.Namespace)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var ef clientv1alpha3.EnvoyFilter
	err = r.Client.Get(ctx, req.NamespacedName, &ef)
	newEf := newEnvoyFilter(hg)

	if err == nil {
		err = r.Client.Update(ctx, newEf, &client.UpdateOptions{})
		if err != nil {
			log.Log.Info("unable to update envoyfilter resource", "name", req.Name, "namespace", req.Namespace)
			return ctrl.Result{Requeue: true}, err
		}

	} else {
		if !apierrors.IsNotFound(err) {
			log.Log.Info("unable to fetch envoyfilter not found", "name", req.Name, "namespace", req.Namespace)
			return ctrl.Result{}, err
		}

		err = r.Client.Create(ctx, newEf, &client.CreateOptions{})
		if err != nil {
			log.Log.Info("unable to create envoyfilter resource", "name", req.Name, "namespace", req.Namespace)
			return ctrl.Result{Requeue: true}, err
		}
	}

	return ctrl.Result{}, nil
}

func newEnvoyFilter(hg appsv1alpha1.HttpGzip) *clientv1alpha3.EnvoyFilter {
	context := typesv1alpha3.EnvoyFilter_SIDECAR_INBOUND
	if hg.Spec.ApplyTo.Kind == appsv1alpha1.Gateway {
		context = typesv1alpha3.EnvoyFilter_GATEWAY
	}

	return &clientv1alpha3.EnvoyFilter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.istio.io/v1alpha3",
			Kind:       "EnvoyFilter",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      hg.Name,
			Namespace: hg.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				metav1.OwnerReference{
					// TODO: fn to get ownerreference based on a resource
				},
			},
		},
		Spec: typesv1alpha3.EnvoyFilter{
			WorkloadSelector: &typesv1alpha3.WorkloadSelector{
				Labels: hg.Spec.ApplyTo.Selector,
			},
			ConfigPatches: []*typesv1alpha3.EnvoyFilter_EnvoyConfigObjectPatch{
				&typesv1alpha3.EnvoyFilter_EnvoyConfigObjectPatch{
					ApplyTo: typesv1alpha3.EnvoyFilter_HTTP_FILTER,
					Match: &typesv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch{
						Context: context,
						ObjectTypes: &typesv1alpha3.EnvoyFilter_EnvoyConfigObjectMatch_Listener{
							Listener: &typesv1alpha3.EnvoyFilter_ListenerMatch{
								FilterChain: &typesv1alpha3.EnvoyFilter_ListenerMatch_FilterChainMatch{

									Filter: &typesv1alpha3.EnvoyFilter_ListenerMatch_FilterMatch{
										Name: "envoy.filters.network.http_connection_manager",
										SubFilter: &typesv1alpha3.EnvoyFilter_ListenerMatch_SubFilterMatch{
											Name: "envoy.filters.http.router",
										},
									},
								},
							},
						},
					},
					Patch: &typesv1alpha3.EnvoyFilter_Patch{
						Operation: typesv1alpha3.EnvoyFilter_Patch_INSERT_BEFORE,
						Value: &_struct.Struct{
							Fields: map[string]*structpb.Value{
								"name": structpb.NewStringValue("envoy.filters.http.compressor"),
								"typed_config": structpb.NewStructValue(&structpb.Struct{
									Fields: map[string]*structpb.Value{
										"@type": structpb.NewStringValue("type.googleapis.com/envoy.extensions.filters.http.compressor.v3.Compressor"),
										"compressor_library": structpb.NewStructValue(&structpb.Struct{
											Fields: map[string]*structpb.Value{
												"name": structpb.NewStringValue("text_optimized"),
												"typed_config": structpb.NewStructValue(
													&structpb.Struct{
														Fields: map[string]*structpb.Value{
															"compression_strategy": structpb.NewStringValue("DEFAULT_STRATEGY"),
															"@type":                structpb.NewStringValue("type.googleapis.com/envoy.extensions.compression.gzip.compressor.v3.Gzip"),
														},
													},
												),
											},
										}),
										"remove_accept_encoding_header": structpb.NewBoolValue(true),
									},
								}),
							},
						},
					},
				},
			},
		}}
}

// SetupWithManager sets up the controller with the Manager.
func (r *HttpGzipReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.HttpGzip{}).
		Complete(r)
}
