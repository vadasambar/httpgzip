package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	httpgzipv1alphav1 "github.com/vadasambar/httpgzip/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("HttpGzip Controller", func() {
	Context("When creating HttpGzip", func() {
		It("Should create the correct EnvoyFilter for pods", func() {
			By("By filling in the right values for EnvoyFilter")
			ctx := context.Background()
			httpgzip := &httpgzipv1alphav1.HttpGzip{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "apps.vadasambar.com/v1alpha1",
					Kind:       "HttpGzip",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "httpgzip-sample",
				},
				Spec: httpgzipv1alphav1.HttpGzipSpec{
					ApplyTo: httpgzipv1alphav1.ApplyTo{
						Kind: httpgzipv1alphav1.Pod,
						Selector: map[string]string{
							"app": "bookshop",
						},
					},
				},
			}

			Expect(k8sClient.Create(ctx, httpgzip)).Should(Succeed())

			Eventually(func() bool {
				var httpgzipCreated httpgzipv1alphav1.HttpGzip
				err := k8sClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: "httpgzip-sample"}, &httpgzipCreated)
				Expect(err).NotTo(HaveOccurred())

				// CONTINUE HERE

				result := httpgzipCreated.Spec.ApplyTo.Kind == httpgzipv1alphav1.Pod && httpgzipCreated.Spec.ApplyTo.Selector["app"] == "bookshop"

				return result

			}, time.Second*30, time.Second*2).Should(BeTrue())
		})
	})
})
