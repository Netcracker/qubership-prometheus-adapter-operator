<!-- markdownlint-disable first-line-heading -->
<!-- markdownlint-disable line-length -->
<!-- markdownlint-disable reference-links-images -->
> This page is automatically generated with `gen-crd-api-reference-docs`.
<p>Packages:</p>
<ul>
<li>
<a href="#monitoring.qubership.org%2fv1alpha1">monitoring.qubership.org/v1alpha1</a>
</li>
</ul>
<h2 id="monitoring.qubership.org/v1alpha1">monitoring.qubership.org/v1alpha1</h2>
Resource Types:
<ul></ul>
<h3 id="monitoring.qubership.org/v1alpha1.Auth">Auth
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.qubership.org/v1alpha1.PrometheusAdapterSpec">PrometheusAdapterSpec</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>basicAuth</code><br/>
<em>
<a href="https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/api.md#monitoring.coreos.com/v1.BasicAuth">
Prometheus operator v1 .BasicAuth
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.qubership.org/v1alpha1.CustomMetricRuleConfig">CustomMetricRuleConfig
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.qubership.org/v1alpha1.CustomScaleMetricRuleSpec">CustomScaleMetricRuleSpec</a>)
</p>
<div>
<p>CustomMetricRuleConfig defines the metric exposing rule from Prometheus.
This structure is similar to the DiscoveryRule from github.com/directxman12/k8s-prometheus-adapter/pkg/config
but we can not use the original structure because it is not compliant with kube-builder&rsquo;s CRD generator.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>seriesQuery</code><br/>
<em>
string
</em>
</td>
<td>
<p>SeriesQuery specifies which metrics this rule should consider via a Prometheus query
series selector query.</p>
</td>
</tr>
<tr>
<td>
<code>seriesFilters</code><br/>
<em>
<a href="#monitoring.qubership.org/v1alpha1.RegexFilter">
[]RegexFilter
</a>
</em>
</td>
<td>
<p>SeriesFilters specifies additional regular expressions to be applied on
the series names returned from the query. This is useful for constraints
that can&rsquo;t be represented in the SeriesQuery (e.g. series matching <code>container_.+</code>
not matching <code>container_.+_total</code>. A filter will be automatically appended to
match the form specified in Name.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="#monitoring.qubership.org/v1alpha1.ResourceMapping">
ResourceMapping
</a>
</em>
</td>
<td>
<p>Resources specifies how associated Kubernetes resources should be discovered for
the given metrics.</p>
</td>
</tr>
<tr>
<td>
<code>name</code><br/>
<em>
<a href="#monitoring.qubership.org/v1alpha1.NameMapping">
NameMapping
</a>
</em>
</td>
<td>
<p>Name specifies how the metric name should be transformed between custom metric
API resources, and Prometheus metric names.</p>
</td>
</tr>
<tr>
<td>
<code>metricsQuery</code><br/>
<em>
string
</em>
</td>
<td>
<p>MetricsQuery specifies modifications to the metrics query, such as converting
cumulative metrics to rate metrics. It is a template where <code>.LabelMatchers</code> is
a the comma-separated base label matchers and <code>.Series</code> is the series name, and
<code>.GroupBy</code> is the comma-separated expected group-by label names. The delimeters
are <code>&lt;&lt;</code> and <code>&gt;&gt;</code>.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.qubership.org/v1alpha1.CustomScaleMetricRule">CustomScaleMetricRule
</h3>
<div>
<p>CustomScaleMetricRule is the Schema for the customscalemetricrules API</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#monitoring.qubership.org/v1alpha1.CustomScaleMetricRuleSpec">
CustomScaleMetricRuleSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>rules</code><br/>
<em>
<a href="#monitoring.qubership.org/v1alpha1.CustomMetricRuleConfig">
[]CustomMetricRuleConfig
</a>
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#monitoring.qubership.org/v1alpha1.CustomScaleMetricRuleStatus">
CustomScaleMetricRuleStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.qubership.org/v1alpha1.CustomScaleMetricRuleSpec">CustomScaleMetricRuleSpec
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.qubership.org/v1alpha1.CustomScaleMetricRule">CustomScaleMetricRule</a>)
</p>
<div>
<p>CustomScaleMetricRuleSpec defines the desired state of CustomScaleMetricRule</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>rules</code><br/>
<em>
<a href="#monitoring.qubership.org/v1alpha1.CustomMetricRuleConfig">
[]CustomMetricRuleConfig
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.qubership.org/v1alpha1.CustomScaleMetricRuleStatus">CustomScaleMetricRuleStatus
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.qubership.org/v1alpha1.CustomScaleMetricRule">CustomScaleMetricRule</a>)
</p>
<div>
<p>CustomScaleMetricRuleStatus defines the observed state of CustomScaleMetricRule</p>
</div>
<h3 id="monitoring.qubership.org/v1alpha1.GroupResource">GroupResource
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.qubership.org/v1alpha1.ResourceMapping">ResourceMapping</a>)
</p>
<div>
<p>GroupResource represents a Kubernetes group-resource.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>group</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>resource</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.qubership.org/v1alpha1.NameMapping">NameMapping
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.qubership.org/v1alpha1.CustomMetricRuleConfig">CustomMetricRuleConfig</a>)
</p>
<div>
<p>NameMapping specifies how to convert Prometheus metrics
to/from custom metrics API resources.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>matches</code><br/>
<em>
string
</em>
</td>
<td>
<p>Matches is a regular expression that is used to match
Prometheus series names.  It may be left blank, in which
case it is equivalent to <code>.*</code>.</p>
</td>
</tr>
<tr>
<td>
<code>as</code><br/>
<em>
string
</em>
</td>
<td>
<p>As is the name used in the API.  Captures from Matches
are available for use here.  If not specified, it defaults
to $0 if no capture groups are present in Matches, or $1
if only one is present, and will error if multiple are.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.qubership.org/v1alpha1.PrometheusAdapter">PrometheusAdapter
</h3>
<div>
<p>PrometheusAdapter is the Schema for the prometheusadapters API</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#monitoring.qubership.org/v1alpha1.PrometheusAdapterSpec">
PrometheusAdapterSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Image to use for a <code>prometheus-adapter</code> deployment.</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Replicas set the expected replicas of the prometheus-adapter. The controller will eventually make the size
of the running replicas equal to the expected size.</p>
</td>
</tr>
<tr>
<td>
<code>prometheusUrl</code><br/>
<em>
string
</em>
</td>
<td>
<p>PrometheusURL used to connect to Prometheus. It will eventually contain query parameters
to configure the connection.</p>
</td>
</tr>
<tr>
<td>
<code>metricsRelistInterval</code><br/>
<em>
string
</em>
</td>
<td>
<p>MetricsRelistInterval is the interval at which to update the cache of available metrics from Prometheus</p>
</td>
</tr>
<tr>
<td>
<code>enableResourceMetrics</code><br/>
<em>
bool
</em>
</td>
<td>
<p>EnableResourceMetrics allows enabling/disabling adapter for <code>metrics.k8s.io</code></p>
</td>
</tr>
<tr>
<td>
<code>enableCustomMetrics</code><br/>
<em>
bool
</em>
</td>
<td>
<p>EnableCustomMetrics allows enabling/disabling adapter for <code>custom.metrics.k8s.io</code></p>
</td>
</tr>
<tr>
<td>
<code>customScaleMetricRulesSelector</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#*k8s.io/apimachinery/pkg/apis/meta/v1.labelselector--">
[]*k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector
</a>
</em>
</td>
<td>
<p>CustomScaleMetricRulesSelector defines label selectors to select
CustomScaleMetricRule resources across the cluster.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines resources requests and limits for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="#monitoring.qubership.org/v1alpha1.SecurityContext">
SecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes.</p>
</td>
</tr>
<tr>
<td>
<code>TLSConfig</code><br/>
<em>
<a href="#monitoring.qubership.org/v1alpha1.TLSConfig">
TLSConfig
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>auth</code><br/>
<em>
<a href="#monitoring.qubership.org/v1alpha1.Auth">
Auth
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.
Specified just as map[string]string. For example: &ldquo;type: compute&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>labels</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Map of string keys and values that can be used to organize and categorize
(scope and select) objects. May match selectors of replication controllers
and services.
More info: <a href="https://kubernetes.io/docs/user-guide/labels">https://kubernetes.io/docs/user-guide/labels</a></p>
</td>
</tr>
<tr>
<td>
<code>annotations</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Annotations is an unstructured key value map stored with a resource that may be
set by external tools to store and retrieve arbitrary metadata. They are not
queryable and should be preserved when modifying objects.
More info: <a href="https://kubernetes.io/docs/user-guide/annotations">https://kubernetes.io/docs/user-guide/annotations</a></p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Affinity is a group of affinity scheduling rules.
More info: <a href="https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node">https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node</a></p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Tolerations allow the pods to schedule onto nodes with matching taints.
More info: <a href="https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration">https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration</a></p>
</td>
</tr>
<tr>
<td>
<code>priorityClassName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>PriorityClassName assigned to the Pods</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#monitoring.qubership.org/v1alpha1.PrometheusAdapterStatus">
PrometheusAdapterStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.qubership.org/v1alpha1.PrometheusAdapterSpec">PrometheusAdapterSpec
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.qubership.org/v1alpha1.PrometheusAdapter">PrometheusAdapter</a>)
</p>
<div>
<p>PrometheusAdapterSpec defines the desired state of PrometheusAdapter</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>image</code><br/>
<em>
string
</em>
</td>
<td>
<p>Image to use for a <code>prometheus-adapter</code> deployment.</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code><br/>
<em>
int32
</em>
</td>
<td>
<p>Replicas set the expected replicas of the prometheus-adapter. The controller will eventually make the size
of the running replicas equal to the expected size.</p>
</td>
</tr>
<tr>
<td>
<code>prometheusUrl</code><br/>
<em>
string
</em>
</td>
<td>
<p>PrometheusURL used to connect to Prometheus. It will eventually contain query parameters
to configure the connection.</p>
</td>
</tr>
<tr>
<td>
<code>metricsRelistInterval</code><br/>
<em>
string
</em>
</td>
<td>
<p>MetricsRelistInterval is the interval at which to update the cache of available metrics from Prometheus</p>
</td>
</tr>
<tr>
<td>
<code>enableResourceMetrics</code><br/>
<em>
bool
</em>
</td>
<td>
<p>EnableResourceMetrics allows enabling/disabling adapter for <code>metrics.k8s.io</code></p>
</td>
</tr>
<tr>
<td>
<code>enableCustomMetrics</code><br/>
<em>
bool
</em>
</td>
<td>
<p>EnableCustomMetrics allows enabling/disabling adapter for <code>custom.metrics.k8s.io</code></p>
</td>
</tr>
<tr>
<td>
<code>customScaleMetricRulesSelector</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#*k8s.io/apimachinery/pkg/apis/meta/v1.labelselector--">
[]*k8s.io/apimachinery/pkg/apis/meta/v1.LabelSelector
</a>
</em>
</td>
<td>
<p>CustomScaleMetricRulesSelector defines label selectors to select
CustomScaleMetricRule resources across the cluster.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resources defines resources requests and limits for single Pods.</p>
</td>
</tr>
<tr>
<td>
<code>securityContext</code><br/>
<em>
<a href="#monitoring.qubership.org/v1alpha1.SecurityContext">
SecurityContext
</a>
</em>
</td>
<td>
<p>SecurityContext holds pod-level security attributes.</p>
</td>
</tr>
<tr>
<td>
<code>TLSConfig</code><br/>
<em>
<a href="#monitoring.qubership.org/v1alpha1.TLSConfig">
TLSConfig
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>auth</code><br/>
<em>
<a href="#monitoring.qubership.org/v1alpha1.Auth">
Auth
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<p>Define which Nodes the Pods are scheduled on.
Specified just as map[string]string. For example: &ldquo;type: compute&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>labels</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Map of string keys and values that can be used to organize and categorize
(scope and select) objects. May match selectors of replication controllers
and services.
More info: <a href="https://kubernetes.io/docs/user-guide/labels">https://kubernetes.io/docs/user-guide/labels</a></p>
</td>
</tr>
<tr>
<td>
<code>annotations</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Annotations is an unstructured key value map stored with a resource that may be
set by external tools to store and retrieve arbitrary metadata. They are not
queryable and should be preserved when modifying objects.
More info: <a href="https://kubernetes.io/docs/user-guide/annotations">https://kubernetes.io/docs/user-guide/annotations</a></p>
</td>
</tr>
<tr>
<td>
<code>affinity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core">
Kubernetes core/v1.Affinity
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Affinity is a group of affinity scheduling rules.
More info: <a href="https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node">https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node</a></p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Tolerations allow the pods to schedule onto nodes with matching taints.
More info: <a href="https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration">https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration</a></p>
</td>
</tr>
<tr>
<td>
<code>priorityClassName</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>PriorityClassName assigned to the Pods</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.qubership.org/v1alpha1.PrometheusAdapterStatus">PrometheusAdapterStatus
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.qubership.org/v1alpha1.PrometheusAdapter">PrometheusAdapter</a>)
</p>
<div>
<p>PrometheusAdapterStatus defines the observed state of PrometheusAdapter</p>
</div>
<h3 id="monitoring.qubership.org/v1alpha1.RegexFilter">RegexFilter
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.qubership.org/v1alpha1.CustomMetricRuleConfig">CustomMetricRuleConfig</a>)
</p>
<div>
<p>RegexFilter is a filter that matches positively or negatively against a regex.
Only one field may be set at a time.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>is</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>isNot</code><br/>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.qubership.org/v1alpha1.ResourceMapping">ResourceMapping
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.qubership.org/v1alpha1.CustomMetricRuleConfig">CustomMetricRuleConfig</a>)
</p>
<div>
<p>ResourceMapping specifies how to map Kubernetes resources to Prometheus labels</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>template</code><br/>
<em>
string
</em>
</td>
<td>
<p>Template specifies a golang string template for converting a Kubernetes
group-resource to a Prometheus label.  The template object contains
the <code>.Group</code> and <code>.Resource</code> fields.  The <code>.Group</code> field will have
dots replaced with underscores, and the <code>.Resource</code> field will be
singularized.  The delimiters are <code>&lt;&lt;</code> and <code>&gt;&gt;</code>.</p>
</td>
</tr>
<tr>
<td>
<code>overrides</code><br/>
<em>
<a href="#monitoring.qubership.org/v1alpha1.GroupResource">
map[string]./v1alpha1/.GroupResource
</a>
</em>
</td>
<td>
<p>Overrides specifies exceptions to the above template, mapping label names
to group-resources</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.qubership.org/v1alpha1.SecurityContext">SecurityContext
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.qubership.org/v1alpha1.PrometheusAdapterSpec">PrometheusAdapterSpec</a>)
</p>
<div>
<p>SecurityContext holds pod-level security attributes.
The parameters are required if a Pod Security Policy is enabled
for Kubernetes cluster and required if a Security Context Constraints is enabled
for Openshift cluster.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>runAsUser</code><br/>
<em>
int64
</em>
</td>
<td>
<p>The UID to run the entrypoint of the container process.
Defaults to user specified in image metadata if unspecified.</p>
</td>
</tr>
<tr>
<td>
<code>fsGroup</code><br/>
<em>
int64
</em>
</td>
<td>
<p>A special supplemental group that applies to all containers in a pod.
Some volume types allow the Kubelet to change the ownership of that volume
to be owned by the pod:</p>
<ol>
<li>The owning GID will be the FSGroup</li>
<li>The setgid bit is set (new files created in the volume will be owned by FSGroup)</li>
<li>The permission bits are OR&rsquo;d with rw-rw&mdash;-</li>
</ol>
<p>If unset, the Kubelet will not modify the ownership and permissions of any volume.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="monitoring.qubership.org/v1alpha1.TLSConfig">TLSConfig
</h3>
<p>
(<em>Appears on:</em><a href="#monitoring.qubership.org/v1alpha1.PrometheusAdapterSpec">PrometheusAdapterSpec</a>)
</p>
<div>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>caSecret</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
<p>Certificate authority used when verifying server certificates.</p>
</td>
</tr>
<tr>
<td>
<code>certSecret</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
<p>Client certificate to present when doing client-authentication.</p>
</td>
</tr>
<tr>
<td>
<code>keySecret</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#secretkeyselector-v1-core">
Kubernetes core/v1.SecretKeySelector
</a>
</em>
</td>
<td>
<p>Secret containing the client key file for the target.</p>
</td>
</tr>
</tbody>
</table>
<hr/>
