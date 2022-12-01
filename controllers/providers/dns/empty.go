package dns

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
	"cloud.example.com/annotation-operator/controllers/depresolver"
	"cloud.example.com/annotation-operator/controllers/providers/assistant"
	"cloud.example.com/annotation-operator/controllers/reconciliation"
	externaldns "sigs.k8s.io/external-dns/endpoint"
)

// EmptyDNSProvider is executed when fakeDNSEnabled is true.
type EmptyDNSProvider struct {
	assistant assistant.Assistant
	config    depresolver.Config
}

func NewEmptyDNS(config depresolver.Config, assistant assistant.Assistant) *EmptyDNSProvider {
	return &EmptyDNSProvider{
		config:    config,
		assistant: assistant,
	}
}

func (p *EmptyDNSProvider) CreateZoneDelegationForExternalDNS(*reconciliation.LoopState) (err error) {
	return
}

func (p *EmptyDNSProvider) IngressExposedIPs(gslb *reconciliation.LoopState) (r []string, err error) {
	return p.assistant.IngressExposedIPs(gslb)
}

func (p *EmptyDNSProvider) GetExternalTargets(host string) (targets assistant.Targets) {
	return p.assistant.GetExternalTargets(host, p.config.GetExternalClusterNSNames())
}

func (p *EmptyDNSProvider) SaveDNSEndpoint(gslb *reconciliation.LoopState, i *externaldns.DNSEndpoint) error {
	return p.assistant.SaveDNSEndpoint(gslb.NamespacedName.Namespace, i)
}

func (p *EmptyDNSProvider) String() string {
	return "EMPTY"
}
