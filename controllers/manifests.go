// Copyright 2025 NetCracker Technology Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllers

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"net/url"
	"strings"

	api "github.com/Netcracker/qubership-prometheus-adapter-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

const (
	promUsernameTmpl = "USERNAME"
	promPasswordTmpl = "PASSWORD"
)

var (
	//go:embed assets/*.yaml
	assets embed.FS

	PrometheusAdapterDeployment     = "assets/deployment.yaml"
	PrometheusAdapterService        = "assets/service.yaml"
	PrometheusAdapterServiceAccount = "assets/serviceaccount.yaml"
	PrometheusAdapterConfigMap      = "assets/configmap.yaml"

	CustomMetricsClusterRole        = "assets/custom-metrics-clusterrole.yaml"
	CustomMetricsClusterRoleBinding = "assets/custom-metrics-clusterrolebinding.yaml"
	TlsCertificateSecretDirCa       = "/etc/adapter/certificates/ca/"
	TlsCertificateSecretDirCert     = "/etc/adapter/certificates/cert/"
	TlsCertificateSecretDirKey      = "/etc/adapter/certificates/key/"

	PrometheusAdapterCustomMetrics = "prometheus-adapter-custom-metrics"

	MinimalReplicasCount int32 = 0
)

// MustAssetReader loads and return the asset for the given name as bytes reader.
// Panics when the asset loading would return an error.
func MustAssetReader(asset string) io.Reader {
	content, _ := assets.ReadFile(asset)
	return bytes.NewReader(content)
}

// Factory provides methods to build resources manifests for prometheus-adapter.
type Factory struct {
	prometheusAdapterCr *api.PrometheusAdapter
}

// NewFactory create an instance of the factory with loaded CR.
func NewFactory(prometheusAdapterCr *api.PrometheusAdapter) *Factory {
	return &Factory{
		prometheusAdapterCr: prometheusAdapterCr,
	}
}

// PrometheusAdapterDeployment builds the Deployemnt resource manifest
// and fill it with parameters from the CR.
func (f *Factory) PrometheusAdapterDeployment() (*appsv1.Deployment, error) {
	d := appsv1.Deployment{}
	err := yaml.NewYAMLOrJSONDecoder(MustAssetReader(PrometheusAdapterDeployment), 100).Decode(&d)

	if err != nil {
		return nil, err
	}

	d.SetNamespace(f.prometheusAdapterCr.GetNamespace())

	// Set container parameters
	for i := range d.Spec.Template.Spec.Containers {
		c := &d.Spec.Template.Spec.Containers[i]
		if c.Name == "prometheus-adapter" {
			c.Image = f.prometheusAdapterCr.Spec.Image
			c.Resources = f.prometheusAdapterCr.Spec.Resources
			if f.prometheusAdapterCr.Spec.Auth != nil && f.prometheusAdapterCr.Spec.Auth.BasicAuth != nil {
				c.Env = append(c.Env, corev1.EnvVar{
					Name: promUsernameTmpl,
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: f.prometheusAdapterCr.Spec.Auth.BasicAuth.Username.Name,
							},
							Key: "username",
						},
					},
				})

				c.Env = append(c.Env, corev1.EnvVar{
					Name: promPasswordTmpl,
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: f.prometheusAdapterCr.Spec.Auth.BasicAuth.Username.Name,
							},
							Key: "password",
						},
					},
				})

				promUrl, err := url.Parse(f.prometheusAdapterCr.Spec.PrometheusURL)

				if err != nil {
					return nil, err
				}

				promUrl.User = url.UserPassword("$("+promUsernameTmpl+")", "$("+promPasswordTmpl+")")
				decodedPromUrl, err := url.QueryUnescape(promUrl.String())

				if err != nil {
					return nil, err
				}

				c.Args = append(c.Args, "--prometheus-url="+decodedPromUrl)
			} else {
				c.Args = append(c.Args, "--prometheus-url="+f.prometheusAdapterCr.Spec.PrometheusURL)
			}
			c.Args = append(c.Args, "--metrics-relist-interval="+f.prometheusAdapterCr.Spec.MetricsRelistInterval)
			if f.prometheusAdapterCr.Spec.TLSEnabled && f.prometheusAdapterCr.Spec.TLSConfig != nil {
				caVolume := corev1.Volume{
					Name: f.prometheusAdapterCr.Spec.TLSConfig.CA.Name,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: f.prometheusAdapterCr.Spec.TLSConfig.CA.Name}}}
				caMount := corev1.VolumeMount{
					MountPath: TlsCertificateSecretDirCa,
					Name:      f.prometheusAdapterCr.Spec.TLSConfig.CA.Name}
				if f.prometheusAdapterCr.Spec.TLSConfig.CA != nil {
					appendCA(f, &d, c, caVolume, caMount)
				}
				if f.prometheusAdapterCr.Spec.TLSConfig.Cert != nil && f.prometheusAdapterCr.Spec.TLSConfig.KeySecret != nil {
					certVolume := corev1.Volume{
						Name: f.prometheusAdapterCr.Spec.TLSConfig.Cert.Name,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: f.prometheusAdapterCr.Spec.TLSConfig.Cert.Name}}}
					certMount := corev1.VolumeMount{
						MountPath: TlsCertificateSecretDirCert,
						Name:      f.prometheusAdapterCr.Spec.TLSConfig.Cert.Name}
					appendCert(f, &d, c, caVolume, caMount, certVolume, certMount)
					appendKey(f, &d, c, caVolume, caMount, certVolume, certMount)
				}
			}
		}
	}

	if f.prometheusAdapterCr.Spec.Replicas != nil && *f.prometheusAdapterCr.Spec.Replicas >= MinimalReplicasCount {
		d.Spec.Replicas = f.prometheusAdapterCr.Spec.Replicas
	}

	d.Spec.Template.Spec.NodeSelector = f.prometheusAdapterCr.Spec.NodeSelector

	if f.prometheusAdapterCr.Spec.SecurityContext != nil {
		d.Spec.Template.Spec.SecurityContext.RunAsUser = f.prometheusAdapterCr.Spec.SecurityContext.RunAsUser
		d.Spec.Template.Spec.SecurityContext.FSGroup = f.prometheusAdapterCr.Spec.SecurityContext.FSGroup
	}

	// Set labels
	d.Labels["app.kubernetes.io/instance"] = getInstanceLabel(d.GetName(), d.GetNamespace())
	d.Labels["app.kubernetes.io/version"] = getTagFromImage(f.prometheusAdapterCr.Spec.Image)

	if f.prometheusAdapterCr.Spec.Labels != nil {
		for k, v := range f.prometheusAdapterCr.Spec.Labels {
			d.Labels[k] = v
		}
	}

	if d.Annotations == nil && f.prometheusAdapterCr.Spec.Annotations != nil {
		d.SetAnnotations(f.prometheusAdapterCr.Spec.Annotations)
	} else {
		for k, v := range f.prometheusAdapterCr.Spec.Annotations {
			d.Annotations[k] = v
		}
	}

	// Set labels
	d.Spec.Template.Labels["app.kubernetes.io/instance"] = getInstanceLabel(d.GetName(), d.GetNamespace())
	d.Spec.Template.Labels["app.kubernetes.io/version"] = getTagFromImage(f.prometheusAdapterCr.Spec.Image)

	if f.prometheusAdapterCr.Spec.Labels != nil {
		for k, v := range f.prometheusAdapterCr.Spec.Labels {
			d.Spec.Template.Labels[k] = v
		}
	}

	if d.Spec.Template.Annotations == nil && f.prometheusAdapterCr.Spec.Annotations != nil {
		d.Spec.Template.SetAnnotations(f.prometheusAdapterCr.Spec.Annotations)
	} else {
		for k, v := range f.prometheusAdapterCr.Spec.Annotations {
			d.Spec.Template.Annotations[k] = v
		}
	}

	if d.Spec.Template.Spec.Affinity == nil && f.prometheusAdapterCr.Spec.Affinity != nil {
		d.Spec.Template.Spec.Affinity = f.prometheusAdapterCr.Spec.Affinity
	}

	if d.Spec.Template.Spec.Tolerations == nil && f.prometheusAdapterCr.Spec.Tolerations != nil {
		d.Spec.Template.Spec.Tolerations = f.prometheusAdapterCr.Spec.Tolerations
	} else {
		copy(d.Spec.Template.Spec.Tolerations, f.prometheusAdapterCr.Spec.Tolerations)
	}

	if len(strings.TrimSpace(f.prometheusAdapterCr.Spec.PriorityClassName)) > 0 {
		d.Spec.Template.Spec.PriorityClassName = f.prometheusAdapterCr.Spec.PriorityClassName
	}

	return &d, nil
}

func appendCA(f *Factory, d *appsv1.Deployment, c *corev1.Container, caVolume corev1.Volume, caMount corev1.VolumeMount) {
	d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, caVolume)
	c.VolumeMounts = append(c.VolumeMounts, caMount)
	c.Args = append(c.Args, fmt.Sprintf("--prometheus-ca-file=%s%s", caMount.MountPath, f.prometheusAdapterCr.Spec.TLSConfig.CA.Key))
}

func appendCert(f *Factory, d *appsv1.Deployment, c *corev1.Container, caVolume corev1.Volume, caMount corev1.VolumeMount, certVolume corev1.Volume, certMount corev1.VolumeMount) {
	if certVolume.Secret.SecretName != caVolume.Secret.SecretName {
		d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, certVolume)
		c.VolumeMounts = append(c.VolumeMounts, certMount)
		c.Args = append(c.Args, fmt.Sprintf("--prometheus-client-tls-cert-file=%s%s", certMount.MountPath, f.prometheusAdapterCr.Spec.TLSConfig.Cert.Key))
	} else {
		c.Args = append(c.Args, fmt.Sprintf("--prometheus-client-tls-cert-file=%s%s", caMount.MountPath, f.prometheusAdapterCr.Spec.TLSConfig.Cert.Key))
	}
}

func appendKey(f *Factory, d *appsv1.Deployment, c *corev1.Container, caVolume corev1.Volume, caMount corev1.VolumeMount, certVolume corev1.Volume, certMount corev1.VolumeMount) {
	keyVolume := corev1.Volume{
		Name: f.prometheusAdapterCr.Spec.TLSConfig.KeySecret.Name,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: f.prometheusAdapterCr.Spec.TLSConfig.KeySecret.Name}}}
	keyMount := corev1.VolumeMount{
		MountPath: TlsCertificateSecretDirKey,
		Name:      f.prometheusAdapterCr.Spec.TLSConfig.KeySecret.Name}
	if keyVolume.Secret.SecretName != caVolume.Secret.SecretName && keyVolume.Secret.SecretName != certVolume.Secret.SecretName {
		d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, keyVolume)
		c.VolumeMounts = append(c.VolumeMounts, keyMount)
		c.Args = append(c.Args, fmt.Sprintf("--prometheus-client-tls-key-file=%s%s", keyMount.MountPath, f.prometheusAdapterCr.Spec.TLSConfig.KeySecret.Key))
	} else if keyVolume.Secret.SecretName == certVolume.Secret.SecretName && keyVolume.Secret.SecretName != caVolume.Secret.SecretName {
		c.Args = append(c.Args, fmt.Sprintf("--prometheus-client-tls-key-file=%s%s", certMount.MountPath, f.prometheusAdapterCr.Spec.TLSConfig.KeySecret.Key))
	} else {
		c.Args = append(c.Args, fmt.Sprintf("--prometheus-client-tls-key-file=%s%s", caMount.MountPath, f.prometheusAdapterCr.Spec.TLSConfig.KeySecret.Key))
	}
}

// PrometheusAdapterService builds the Service resource manifest
// and fill it with parameters from the CR.
func (f *Factory) PrometheusAdapterService() (*corev1.Service, error) {
	s := corev1.Service{}
	err := yaml.NewYAMLOrJSONDecoder(MustAssetReader(PrometheusAdapterService), 100).Decode(&s)

	if err != nil {
		return nil, err
	}

	s.SetNamespace(f.prometheusAdapterCr.GetNamespace())

	return &s, nil
}

// PrometheusAdapterServiceAccount builds the ServiceAccount resource manifest
// and fill it with parameters from the CR.
func (f *Factory) PrometheusAdapterServiceAccount() (*corev1.ServiceAccount, error) {
	sa := corev1.ServiceAccount{}
	err := yaml.NewYAMLOrJSONDecoder(MustAssetReader(PrometheusAdapterServiceAccount), 100).Decode(&sa)

	if err != nil {
		return nil, err
	}

	sa.SetNamespace(f.prometheusAdapterCr.GetNamespace())

	return &sa, nil
}

// CustomMetricsClusterRole builds the ClusterRole resource manifest for custom metrics
// and fill it with parameters from the CR.
func (f *Factory) CustomMetricsClusterRole() (*rbacv1.ClusterRole, error) {
	cr := rbacv1.ClusterRole{}
	err := yaml.NewYAMLOrJSONDecoder(MustAssetReader(CustomMetricsClusterRole), 100).Decode(&cr)

	if err != nil {
		return nil, err
	}

	cr.SetNamespace(f.prometheusAdapterCr.GetNamespace())
	cr.SetName(cr.GetNamespace() + "-" + PrometheusAdapterCustomMetrics)
	if f.prometheusAdapterCr.Spec.EnableResourceMetrics || f.prometheusAdapterCr.Spec.EnableCustomMetrics {
		cr.Rules = append(cr.Rules, rbacv1.PolicyRule{
			Verbs:     []string{"get", "list", "watch"},
			APIGroups: []string{""},
			Resources: []string{"nodes"},
		}, rbacv1.PolicyRule{
			Verbs:     []string{"get", "list", "watch"},
			APIGroups: []string{"metrics.k8s.io"},
			Resources: []string{"*"},
		}, rbacv1.PolicyRule{
			Verbs:     []string{"get", "list", "watch"},
			APIGroups: []string{"custom.metrics.k8s.io"},
			Resources: []string{"*"},
		})
	}
	return &cr, nil
}

// CustomMetricsConfigMap builds the ConfigMap resource manifest for custom metrics
// and fill it with parameters from the CR.
func (f *Factory) CustomMetricsConfigMap() (*corev1.ConfigMap, error) {
	cm := corev1.ConfigMap{}
	err := yaml.NewYAMLOrJSONDecoder(MustAssetReader(PrometheusAdapterConfigMap), 100).Decode(&cm)

	if err != nil {
		return nil, err
	}

	cm.SetNamespace(f.prometheusAdapterCr.GetNamespace())

	return &cm, nil
}

// CustomMetricsClusterRoleBinding builds the ClusterRoleBinding resource manifest for custom metrics
// and fill it with parameters from the CR.
func (f *Factory) CustomMetricsClusterRoleBinding() (*rbacv1.ClusterRoleBinding, error) {
	crb := rbacv1.ClusterRoleBinding{}
	err := yaml.NewYAMLOrJSONDecoder(MustAssetReader(CustomMetricsClusterRoleBinding), 100).Decode(&crb)

	if err != nil {
		return nil, err
	}

	crb.SetNamespace(f.prometheusAdapterCr.GetNamespace())
	crb.SetName(crb.GetNamespace() + "-" + PrometheusAdapterCustomMetrics)

	crb.RoleRef.Name = crb.GetNamespace() + "-" + PrometheusAdapterCustomMetrics
	for i := range crb.Subjects {
		s := &crb.Subjects[i]
		s.Namespace = f.prometheusAdapterCr.GetNamespace()
	}

	return &crb, nil
}
