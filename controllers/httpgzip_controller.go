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
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/vadasambar/httpgzip/api/v1alpha1"
	typev1alpha3 "istio.io/api/networking/v1alpha3"
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

	if strings.TrimSpace(string(hg.Spec.ApplyTo.Kind)) == "pod" {

	}

	ef := &clientv1alpha3.EnvoyFilter{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.istio.io/v1alpha3",
			Kind:       "EnvoyFilter",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      hg.Name,
			Namespace: hg.Namespace,
		},
		Spec: typev1alpha3.EnvoyFilter{
			WorkloadSelector: &typev1alpha3.WorkloadSelector{
				Labels: hg.Spec.ApplyTo.Selector,
			},
			ConfigPatches: []*typev1alpha3.EnvoyFilter_EnvoyConfigObjectPatch{
				&typev1alpha3.EnvoyFilter_EnvoyConfigObjectPatch{
					ApplyTo: typev1alpha3.EnvoyFilter_HTTP_FILTER,
					Match: &typev1alpha3.EnvoyFilter_EnvoyConfigObjectMatch{
						Context: typev1alpha3.EnvoyFilter_PatchContext(),
					},
					Patch: &typev1alpha3.Patch{},
				},
			},
		},
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HttpGzipReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.HttpGzip{}).
		Complete(r)
}
