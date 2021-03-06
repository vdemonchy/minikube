/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

package kubeadm

import "html/template"

var kubeadmConfigTemplate = template.Must(template.New("kubeadmConfigTemplate").Parse(`apiVersion: kubeadm.k8s.io/v1alpha1
kind: MasterConfiguration
api:
  advertiseAddress: {{.AdvertiseAddress}}
  bindPort: {{.APIServerPort}}
kubernetesVersion: {{.KubernetesVersion}}
certificatesDir: {{.CertDir}}
networking:
  serviceSubnet: {{.ServiceCIDR}}
etcd:
  dataDir: {{.EtcdDataDir}}
nodeName: {{.NodeName}}
{{range .ExtraArgs}}{{.Component}}:{{range $key, $value := .Options}}
  {{$key}}: {{$value}}
{{end}}{{end}}`))

var kubeletSystemdTemplate = template.Must(template.New("kubeletSystemdTemplate").Parse(`
[Service]
Environment="KUBELET_KUBECONFIG_ARGS=--kubeconfig=/etc/kubernetes/kubelet.conf --require-kubeconfig=true"
Environment="KUBELET_SYSTEM_PODS_ARGS=--pod-manifest-path=/etc/kubernetes/manifests --allow-privileged=true"
Environment="KUBELET_DNS_ARGS=--cluster-dns=10.0.0.10 --cluster-domain=cluster.local"
Environment="KUBELET_CADVISOR_ARGS=--cadvisor-port=0"
Environment="KUBELET_CGROUP_ARGS=--cgroup-driver=cgroupfs"
ExecStart=
ExecStart=/usr/bin/kubelet $KUBELET_KUBECONFIG_ARGS $KUBELET_SYSTEM_PODS_ARGS $KUBELET_DNS_ARGS $KUBELET_CADVISOR_ARGS $KUBELET_CGROUP_ARGS {{.ExtraOptions}}
`))

const kubeletService = `
[Unit]
Description=kubelet: The Kubernetes Node Agent
Documentation=http://kubernetes.io/docs/

[Service]
ExecStart=/usr/bin/kubelet
Restart=always
StartLimitInterval=0
RestartSec=10

[Install]
WantedBy=multi-user.target
`

var kubeadmRestoreTemplate = template.Must(template.New("kubeadmRestoreTemplate").Parse(`
sudo kubeadm alpha phase certs all --config {{.KubeadmConfigFile}} &&
sudo /usr/bin/kubeadm alpha phase kubeconfig all --config {{.KubeadmConfigFile}} &&
sudo /usr/bin/kubeadm alpha phase controlplane all --config {{.KubeadmConfigFile}} &&
sudo /usr/bin/kubeadm alpha phase etcd local --config {{.KubeadmConfigFile}}
`))

var kubeadmInitTemplate = template.Must(template.New("kubeadmInitTemplate").Parse("sudo /usr/bin/kubeadm init --config {{.KubeadmConfigFile}} --skip-preflight-checks"))
