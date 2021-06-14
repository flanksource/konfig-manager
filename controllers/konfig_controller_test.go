package controllers

import (
	"context"
	"time"

	konfigmanagerv1 "github.com/flanksource/konfig-manager/api/v1"
	"github.com/flanksource/konfig-manager/pkg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Konfig controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		KonfigName      = "test-konfig"
		KonfigNamespace = "default"
		outputKind      = configMapKind
		outputName      = "properties"
		outputNamespace = "default"
		timeout         = time.Second * 10
		interval        = time.Millisecond * 250
	)
	Context("When Creating new Konfig Object", func() {
		It("Should create specified output type object", func() {
			By("creating a new konfig object")
			ctx := context.Background()
			configMap := &v1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					Kind:       configMapKind,
					APIVersion: coreAPIVersion,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-cm",
					Namespace: "default",
				},
				Data: map[string]string{
					"test-key": "test-value",
				},
			}
			Expect(k8sClient.Create(ctx, configMap)).Should(Succeed())

			secret := &v1.Secret{
				TypeMeta: metav1.TypeMeta{
					Kind:       secretKind,
					APIVersion: coreAPIVersion,
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "test-secret",
				},
				Data: map[string][]byte{
					"secret-key": []byte("secret-value"),
				},
			}
			Expect(k8sClient.Create(ctx, secret)).Should(Succeed())

			konfig := &konfigmanagerv1.Konfig{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Konfig",
					APIVersion: "konfigmanager.flanksource.com/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: KonfigNamespace,
					Name:      KonfigName,
				},
				Spec: konfigmanagerv1.KonfigSpec{
					Hierarchy: []pkg.Item{
						{
							Kind:      "ConfigMap",
							Namespace: "default",
							Name:      "test-cm",
						},
						{
							Kind:      "Secret",
							Name:      "test-secret",
							Namespace: "default",
						},
					},
					Output: konfigmanagerv1.Output{
						Name:      outputName,
						Namespace: outputNamespace,
						Kind:      outputKind,
					},
				},
			}
			Expect(k8sClient.Create(ctx, konfig)).Should(Succeed())

			konfigLookupKey := types.NamespacedName{Name: KonfigName, Namespace: KonfigNamespace}
			createdKonfig := &konfigmanagerv1.Konfig{}

			// We'll need to retry getting this newly created Konfig, given that creation may not immediately happen.
			Eventually(func() bool {
				err := k8sClient.Get(ctx, konfigLookupKey, createdKonfig)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			outputLookupKey := types.NamespacedName{Name: outputName, Namespace: outputNamespace}
			createdOutput := &v1.ConfigMap{}
			By("By checking the Konfig Created object has required properties")
			Eventually(func() bool {
				err := k8sClient.Get(ctx, outputLookupKey, createdOutput)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdOutput.Data["test-key"]).Should(Equal("test-value"))
			Expect(createdOutput.Data["secret-key"]).Should(Equal("secret-value"))
		})
	})
})
