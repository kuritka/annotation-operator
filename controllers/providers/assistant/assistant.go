package assistant

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
	"time"

	externaldns "sigs.k8s.io/external-dns/endpoint"
)

type Assistant interface {
	// CoreDNSExposedIPs retrieves list of exposed IP by CoreDNS
	CoreDNSExposedIPs() ([]string, error)
	// GetExternalTargets retrieves slice of targets from external clusters
	GetExternalTargets(host, primaryGeoTag string, extClusterNsNames map[string]string) (targets Targets)
	// SaveDNSEndpoint update DNS endpoint or create new one if doesnt exist
	SaveDNSEndpoint(namespace string, i *externaldns.DNSEndpoint) error
	// RemoveEndpoint removes endpoint
	RemoveEndpoint(endpointName string) error
	// InspectTXTThreshold inspects fqdn TXT record from edgeDNSServer. If record doesn't exists or timestamp is greater than
	// splitBrainThreshold the error is returned. In case fakeDNSEnabled is true, 127.0.0.1:7753 is used as edgeDNSServer
	InspectTXTThreshold(fqdn string, splitBrainThreshold time.Duration) error
}
