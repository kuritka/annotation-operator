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
	"cloud.example.com/annotation-operator/controllers/mapper"
	"cloud.example.com/annotation-operator/controllers/providers/assistant"
	externaldns "sigs.k8s.io/external-dns/endpoint"
)

type Provider interface {
	// CreateZoneDelegationForExternalDNS handles delegated zone in Edge DNS
	CreateZoneDelegationForExternalDNS(*mapper.LoopState) error
	// GetExternalTargets retrieves list of external targets for specified host
	GetExternalTargets(string, string) assistant.Targets
	// SaveDNSEndpoint update DNS endpoint in gslb or create new one if doesn't exist
	SaveDNSEndpoint(*mapper.LoopState, *externaldns.DNSEndpoint) error
	// String see: Stringer interface
	String() string
	// RequireFinalizer tells whether provider requires to collect any resources
	RequireFinalizer() bool
	// Finalize would be implemented when RequireFinalizer is true
	Finalize(*mapper.LoopState) error
}
