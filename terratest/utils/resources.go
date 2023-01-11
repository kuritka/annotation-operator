package utils

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
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"
	networkingv1 "k8s.io/api/networking/v1"
)

type Resources struct {
	i *Instance
}

type IngressStatus struct {
	// Associated Service status
	ServiceHealth map[string]string `json:"serviceHealth"`
	// Current Healthy DNS record structure
	HealthyRecords map[string][]string `json:"healthyRecords"`
	// Cluster Geo Tag
	GeoTag string `json:"geoTag"`
	// Comma-separated list of hosts. Duplicating the value from range .spec.ingress.rules[*].host for printer column
	Hosts string `json:"hosts,omitempty"`
}

func (r *Resources) Ingress() *networkingv1.Ingress {
	return k8s.GetIngress(r.i.w.t, r.i.w.k8sOptions, r.i.w.ingress.name)
}

func (r *Resources) IngressStatus() (IngressStatus, error) {
	const statusAnnotation = "k8gb.io/status"
	status := IngressStatus{}
	js := r.Ingress().Annotations[statusAnnotation]
	err := json.Unmarshal([]byte(js), &status)
	return status, err
}

func (r *Resources) GetLocalDNSEndpoint() DNSEndpoint {
	ep, err := r.getDNSEndpoint(r.i.w.ingress.name, r.i.w.namespace)
	r.i.continueIfK8sResourceNotFound(err)
	return ep
}

func (r *Resources) GetExternalDNSEndpoint() DNSEndpoint {
	ep, err := r.getDNSEndpoint("k8gb-ns-extdns", r.i.w.k8gbNamespace)
	r.i.continueIfK8sResourceNotFound(err)
	return ep
}

func (r *Resources) getDNSEndpoint(epName, ns string) (ep DNSEndpoint, err error) {
	ep = DNSEndpoint{}
	j, err := k8s.RunKubectlAndGetOutputE(r.i.w.t, r.i.w.k8sOptions, "get", "dnsendpoints.externaldns.k8s.io", epName, "-ojson")
	if err != nil {
		return ep, err
	}
	err = json.Unmarshal([]byte(j), &ep)
	return ep, err
}

func (r *Resources) WaitUntilDNSEndpointNotFound() error {
	for i := 0; i < defaultRetries; i++ {
		_, err := r.getDNSEndpoint(r.i.w.ingress.name, r.i.w.namespace)
		if err != nil && strings.HasSuffix(err.Error(), "not found") {
			r.i.w.t.Logf("SUCCEED: DNSEndpoint %s doesn't exists", r.i.w.ingress.name)
			return nil
		}
		r.i.w.t.Logf("Wait until DNSEndpoint %s will be removed", r.i.w.ingress.name)
		time.Sleep(defaultSeconds * time.Second)
	}
	return fmt.Errorf("ERROR: DNSEndpoint %s exists; but should be removed", r.i.w.ingress.name)
}

func (r *Resources) WaitUntilDNSEndpointContainsTargets(host string, targets []string) (err error) {
	var ep Endpoint
	for i := 0; i < defaultRetries; i++ {
		ep = r.GetLocalDNSEndpoint().GetEndpointByName(host)
		if EqualItems(ep.Targets, targets) {
			r.i.w.t.Logf("SUCCEED: DNSEndpoint has expected targets %s for host %s", ep.Targets, host)
			return nil
		}
		r.i.w.t.Logf("Wait until DNSEndpoint has targets. host: %s has targets %s but expect %s", host, ep.Targets, targets)
		time.Sleep(defaultSeconds * time.Second)
	}
	return fmt.Errorf("FAIL: DNSEndpoint has no proper targets. host: %s has targets %s but expect %s", host, ep.Targets, targets)
}
