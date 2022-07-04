package controllers

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	httpgzipv1alphav1 "github.com/vadasambar/httpgzip/api/v1alpha1"
	networkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var _ = Describe("HttpGzip Controller", func() {
	Context("When creating HttpGzip", func() {
		It("Should create the correct HttpGzip and EnvoyFilter for pods", func() {
			By("By filling in the right values for EnvoyFilter")
			ctx := context.Background()
			httpgzip := &httpgzipv1alphav1.HttpGzip{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "apps.vadasambar.com/v1alpha1",
					Kind:       "HttpGzip",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "httpgzip-sample-pods",
					Namespace: "default",
				},
				Spec: httpgzipv1alphav1.HttpGzipSpec{
					ApplyTo: httpgzipv1alphav1.ApplyTo{
						Kind: httpgzipv1alphav1.Pod,
						Selector: map[string]string{
							"app": "productpage",
						},
					},
				},
			}

			Expect(k8sClient.Create(ctx, httpgzip)).Should(Succeed())

			// 1. Read envoy filter test file
			// 2. Load it as an envoy filter
			// 3. Compare all the fields in the loaded file with the envoy filter fetched from the api-server
			d, err := os.ReadFile("../testfiles/envoy_pod_filter.yaml")
			Expect(err).NotTo(HaveOccurred())

			var envoyPodFilter networkingv1alpha3.EnvoyFilter
			err = yaml.Unmarshal(d, &envoyPodFilter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				var httpgzipCreated httpgzipv1alphav1.HttpGzip
				err := k8sClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: "httpgzip-sample-pods"}, &httpgzipCreated)
				if err != nil {
					fmt.Println("Error getting HttpGzip", err)
					return false
				}

				var envoyFilter networkingv1alpha3.EnvoyFilter
				err = k8sClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: "httpgzip-sample-pods"}, &envoyFilter)
				if err != nil {
					fmt.Println("Error getting EnvoyFilter", err)
					return false
				}

				envoyFilter.OwnerReferences[0].UID = ""

				result := httpgzipCreated.Spec.ApplyTo.Kind == httpgzipv1alphav1.Pod &&
					httpgzipCreated.Spec.ApplyTo.Selector["app"] == "productpage" &&
					reflect.DeepEqual(envoyPodFilter.Spec.DeepCopy(), envoyFilter.Spec.DeepCopy()) &&
					// DeepCopy is not needed for comparing OwnerReferences
					reflect.DeepEqual(envoyPodFilter.OwnerReferences, envoyFilter.OwnerReferences)

				return result

			}, time.Second*30, time.Second*2).Should(BeTrue())
		})
		It("Should create the correct HttpGzip and EnvoyFilter for Istio gateways", func() {
			ctx := context.Background()
			httpgzip := &httpgzipv1alphav1.HttpGzip{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "apps.vadasambar.com/v1alpha1",
					Kind:       "HttpGzip",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "httpgzip-sample-gateway",
					Namespace: "default",
				},
				Spec: httpgzipv1alphav1.HttpGzipSpec{
					ApplyTo: httpgzipv1alphav1.ApplyTo{
						Kind: httpgzipv1alphav1.Gateway,
						Selector: map[string]string{
							"app": "productpage",
						},
					},
				},
			}

			Expect(k8sClient.Create(ctx, httpgzip)).Should(Succeed())

			// 1. Read envoy filter test file
			// 2. Load it as an envoy filter
			// 3. Compare all the fields in the loaded file with the envoy filter fetched from the api-server
			d, err := os.ReadFile("../testfiles/envoy_gateway_filter.yaml")
			Expect(err).NotTo(HaveOccurred())

			var envoyGatewayFilter networkingv1alpha3.EnvoyFilter
			err = yaml.Unmarshal(d, &envoyGatewayFilter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() bool {
				var httpgzipCreated httpgzipv1alphav1.HttpGzip
				err := k8sClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: "httpgzip-sample-gateway"}, &httpgzipCreated)
				if err != nil {
					fmt.Println("Error getting HttpGzip", err)
					return false
				}

				var envoyFilter networkingv1alpha3.EnvoyFilter
				err = k8sClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: "httpgzip-sample-gateway"}, &envoyFilter)
				if err != nil {
					fmt.Println("Error getting EnvoyFilter", err)
					return false
				}

				envoyFilter.OwnerReferences[0].UID = ""

				result := httpgzipCreated.Spec.ApplyTo.Kind == httpgzipv1alphav1.Gateway &&
					httpgzipCreated.Spec.ApplyTo.Selector["app"] == "productpage" &&
					reflect.DeepEqual(envoyGatewayFilter.Spec.DeepCopy(), envoyFilter.Spec.DeepCopy()) &&
					// DeepCopy is not needed for comparing OwnerReferences
					reflect.DeepEqual(envoyGatewayFilter.OwnerReferences, envoyFilter.OwnerReferences)

				return result

			}, time.Second*30, time.Second*2).Should(BeTrue())
		})
	})
})
