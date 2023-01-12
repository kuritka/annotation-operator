package terratest

import (
	"testing"

	"github.com/kuritka/annotation-operator/terratest"
	"github.com/kuritka/annotation-operator/terratest/utils"
	"github.com/stretchr/testify/require"
)

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

func TestFailoverLifecycleOnThreeClusters(t *testing.T) {
	const ingressPath = "./resources/ingress_fo3.yaml"
	//const digHits = 10
	//const wgetHits = 8
	instanceEU, err := utils.NewWorkflow(t, terratest.Environment.EUCluster, terratest.Environment.EUClusterPort).
		WithIngress(ingressPath).
		WithTestApp(terratest.Environment.EUCluster).
		WithBusybox().
		Start()
	require.NoError(t, err)
	defer instanceEU.Kill()

	instanceUS, err := utils.NewWorkflow(t, terratest.Environment.USCluster, terratest.Environment.USClusterPort).
		WithIngress(ingressPath).
		WithTestApp(terratest.Environment.USCluster).
		WithBusybox().
		Start()
	require.NoError(t, err)
	defer instanceUS.Kill()

	instanceZA, err := utils.NewWorkflow(t, terratest.Environment.ZACluster, terratest.Environment.ZAClusterPort).
		WithIngress(ingressPath).
		WithTestApp(terratest.Environment.ZACluster).
		WithBusybox().
		Start()
	require.NoError(t, err)
	defer instanceZA.Kill()

	t.Run("Wait until EU, US, ZA clusters are ready", func(t *testing.T) {
		usClusterIPs := instanceUS.GetInfo().NodeIPs
		err = instanceEU.Resources().WaitUntilDNSEndpointContainsTargets(instanceEU.GetInfo().Host, usClusterIPs)
		require.NoError(t, err)
		err = instanceUS.Resources().WaitUntilDNSEndpointContainsTargets(instanceUS.GetInfo().Host, usClusterIPs)
		require.NoError(t, err)
		err = instanceZA.Resources().WaitUntilDNSEndpointContainsTargets(instanceZA.GetInfo().Host, usClusterIPs)
		require.NoError(t, err)
	})

	t.Logf("All clusters are running 🚜💨! 🇪🇺 %s;🇺🇲 %s; 🇿🇦 %s;",
		terratest.Environment.EUCluster,
		terratest.Environment.USCluster,
		terratest.Environment.ZACluster)

	//t.Run("🇪🇺 Digging US,EU,ZA cluster, IPs of EU are returned", func(t *testing.T) {
	//	euClusterIPs := utils.Merge(instanceEU.GetInfo().NodeIPs)
	//	ips := instanceEU.Tools().DigNCoreDNS(digHits)
	//	require.True(t, utils.MapHasOnlyKeys(ips, euClusterIPs...))
	//	ips = instanceUS.Tools().DigNCoreDNS(digHits)
	//	require.True(t, utils.MapHasOnlyKeys(ips, euClusterIPs...))
	//})
	//
	//t.Run("🇪🇺 Digging US,EU cluster, IPs of EU are returned", func(t *testing.T) {
	//	euClusterIPs := utils.Merge(instanceEU.GetInfo().NodeIPs)
	//	ips := instanceEU.Tools().DigNCoreDNS(digHits)
	//	require.True(t, utils.MapHasOnlyKeys(ips, euClusterIPs...))
	//	ips = instanceUS.Tools().DigNCoreDNS(digHits)
	//	require.True(t, utils.MapHasOnlyKeys(ips, euClusterIPs...))
	//})
	//
	//t.Run("🇪🇺 Wget application, EU,US clusters, returns only EU app", func(t *testing.T) {
	//	instanceHit := instanceEU.Tools().WgetNTestApp(wgetHits)
	//	require.True(t, utils.MapHasOnlyKeys(instanceHit, terratest.Environment.EUCluster))
	//	instanceHit = instanceUS.Tools().WgetNTestApp(wgetHits)
	//	require.True(t, utils.MapHasOnlyKeys(instanceHit, terratest.Environment.EUCluster))
	//})
	//
	//t.Run("💀🇪🇺 Killing app on Primary EU cluster", func(t *testing.T) {
	//	usClusterIPs := utils.Merge(instanceUS.GetInfo().NodeIPs)
	//	instanceEU.App().StopTestApp()
	//	err = instanceEU.Resources().WaitUntilDNSEndpointContainsTargets(instanceEU.GetInfo().Host, usClusterIPs)
	//	require.NoError(t, err, "WARNING: If you running test locally, ensure the App with same host IS NOT running in forgotten namespaces")
	//	err = instanceUS.Resources().WaitUntilDNSEndpointContainsTargets(instanceUS.GetInfo().Host, usClusterIPs)
	//	require.NoError(t, err, "WARNING: If you running test locally, ensure the App with same host IS NOT running in forgotten namespaces")
	//})
	//
	//t.Run("🇺🇸 Digging US,EU cluster, IPs of US are returned", func(t *testing.T) {
	//	usClusterIPs := utils.Merge(instanceUS.GetInfo().NodeIPs)
	//	ips := instanceEU.Tools().DigNCoreDNS(digHits)
	//	require.True(t, utils.MapHasOnlyKeys(ips, usClusterIPs...))
	//	ips = instanceUS.Tools().DigNCoreDNS(digHits)
	//	require.True(t, utils.MapHasOnlyKeys(ips, usClusterIPs...))
	//})
	//
	//t.Run("🇺🇸 Wget application US clusters, returns only US app", func(t *testing.T) {
	//	instanceHit := instanceUS.Tools().WgetNTestApp(wgetHits)
	//	require.True(t, utils.MapHasOnlyKeys(instanceHit, terratest.Environment.USCluster))
	//})
	//
	//t.Run("💀🇺🇸 Killing app on Secondary US cluster", func(t *testing.T) {
	//	instanceUS.App().StopTestApp()
	//	err = instanceUS.Resources().WaitUntilDNSEndpointContainsTargets(instanceUS.GetInfo().Host, []string{})
	//	require.NoError(t, err, "WARNING: If you running test locally, ensure the App with same host IS NOT running in forgotten namespaces")
	//	err = instanceEU.Resources().WaitUntilDNSEndpointContainsTargets(instanceEU.GetInfo().Host, []string{})
	//	require.NoError(t, err, "WARNING: If you running test locally, ensure the App with same host IS NOT running in forgotten namespaces")
	//})
	//
	//t.Run("💀💀 Digging US,EU cluster, empty IPs are returned", func(t *testing.T) {
	//	ips := instanceEU.Tools().DigNCoreDNS(digHits)
	//	require.True(t, utils.MapHasOnlyKeys(ips))
	//	ips = instanceUS.Tools().DigNCoreDNS(digHits)
	//	require.True(t, utils.MapHasOnlyKeys(ips))
	//})
	//
	//t.Run("⚡🇪🇺 Spin up Primary cluster - EU", func(t *testing.T) {
	//	instanceEU.App().StartTestApp()
	//	euClusterIPs := utils.Merge(instanceEU.GetInfo().NodeIPs)
	//	err = instanceUS.Resources().WaitUntilDNSEndpointContainsTargets(instanceUS.GetInfo().Host, euClusterIPs)
	//	require.NoError(t, err, "WARNING: If you running test locally, ensure the App with same host IS NOT running in forgotten namespaces")
	//	err = instanceEU.Resources().WaitUntilDNSEndpointContainsTargets(instanceEU.GetInfo().Host, euClusterIPs)
	//	require.NoError(t, err, "WARNING: If you running test locally, ensure the App with same host IS NOT running in forgotten namespaces")
	//})
	//
	//t.Run("🇪🇺 Digging US,EU cluster, IPs of US are returned", func(t *testing.T) {
	//	euClusterIPs := utils.Merge(instanceEU.GetInfo().NodeIPs)
	//	ips := instanceEU.Tools().DigNCoreDNS(digHits)
	//	require.True(t, utils.MapHasOnlyKeys(ips, euClusterIPs...))
	//	ips = instanceEU.Tools().DigNCoreDNS(digHits)
	//	require.True(t, utils.MapHasOnlyKeys(ips, euClusterIPs...))
	//})
	//
	//t.Run("🇪🇺 Wget application US clusters, returns only EU app", func(t *testing.T) {
	//	instanceHit := instanceEU.Tools().WgetNTestApp(wgetHits)
	//	require.True(t, utils.MapHasOnlyKeys(instanceHit, terratest.Environment.EUCluster))
	//	instanceHit = instanceUS.Tools().WgetNTestApp(wgetHits)
	//	require.True(t, utils.MapHasOnlyKeys(instanceHit, terratest.Environment.EUCluster))
	//})
	//
	//t.Run("⚡🇺🇸🇪🇺 Spin up Secondary cluster - US", func(t *testing.T) {
	//	instanceUS.App().StartTestApp()
	//	euIPs := instanceEU.GetInfo().NodeIPs
	//	err = instanceUS.Resources().WaitUntilDNSEndpointContainsTargets(instanceUS.GetInfo().Host, euIPs)
	//	require.NoError(t, err, "WARNING: If you running test locally, ensure the App with same host IS NOT running in forgotten namespaces")
	//	err = instanceEU.Resources().WaitUntilDNSEndpointContainsTargets(instanceEU.GetInfo().Host, euIPs)
	//	require.NoError(t, err, "WARNING: If you running test locally, ensure the App with same host IS NOT running in forgotten namespaces")
	//})
	//
	//t.Run("🇪🇺 Digging US,EU cluster, IPs of Primary EU are returned", func(t *testing.T) {
	//	euClusterIPs := instanceEU.GetInfo().NodeIPs
	//	ips := instanceEU.Tools().DigNCoreDNS(digHits)
	//	require.True(t, utils.MapHasOnlyKeys(ips, euClusterIPs...))
	//	ips = instanceEU.Tools().DigNCoreDNS(digHits)
	//	require.True(t, utils.MapHasOnlyKeys(ips, euClusterIPs...))
	//})
	//
	//t.Run("🇪🇺 Wget application US clusters, returns only EU app", func(t *testing.T) {
	//	instanceHit := instanceEU.Tools().WgetNTestApp(wgetHits)
	//	require.True(t, utils.MapHasOnlyKeys(instanceHit, terratest.Environment.EUCluster))
	//	instanceHit = instanceUS.Tools().WgetNTestApp(wgetHits)
	//	require.True(t, utils.MapHasOnlyKeys(instanceHit, terratest.Environment.EUCluster))
	//})
}
