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

package localkube

import (
	apiserver "k8s.io/kubernetes/cmd/kube-apiserver/app"
	"k8s.io/kubernetes/cmd/kube-apiserver/app/options"

	"k8s.io/kubernetes/pkg/storage/storagebackend"
)

func (lk LocalkubeServer) NewAPIServer() Server {
	return NewSimpleServer("apiserver", serverInterval, StartAPIServer(lk))
}

func StartAPIServer(lk LocalkubeServer) func() error {
	config := options.NewServerRunOptions()

	config.GenericServerRunOptions.BindAddress = lk.APIServerAddress
	config.GenericServerRunOptions.SecurePort = lk.APIServerPort
	config.GenericServerRunOptions.InsecureBindAddress = lk.APIServerInsecureAddress
	config.GenericServerRunOptions.InsecurePort = lk.APIServerInsecurePort

	config.GenericServerRunOptions.ClientCAFile = lk.GetCAPublicKeyCertPath()
	config.GenericServerRunOptions.TLSCertFile = lk.GetPublicKeyCertPath()
	config.GenericServerRunOptions.TLSPrivateKeyFile = lk.GetPrivateKeyCertPath()
	config.GenericServerRunOptions.AdmissionControl = "NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,ResourceQuota"

	// use localkube etcd
	config.GenericServerRunOptions.StorageConfig = storagebackend.Config{ServerList: KubeEtcdClientURLs}

	// set Service IP range
	config.GenericServerRunOptions.ServiceClusterIPRange = lk.ServiceClusterIPRange

	// defaults from apiserver command
	config.GenericServerRunOptions.EnableProfiling = true
	config.GenericServerRunOptions.EnableWatchCache = true
	config.GenericServerRunOptions.MinRequestTimeout = 1800

	config.AllowPrivileged = true

	config.GenericServerRunOptions.RuntimeConfig = lk.RuntimeConfig

	lk.SetExtraConfigForComponent("apiserver", &config)

	return func() error {
		return apiserver.Run(config)
	}
}
