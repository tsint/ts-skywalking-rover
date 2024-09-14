// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package common

import (
	"strings"

	"github.com/apache/skywalking-rover/pkg/process/api"
	"github.com/apache/skywalking-rover/pkg/process/finders/kubernetes"
)

type MonitorFilter interface {
	// ShouldExclude returns true if the process should be excluded from monitoring.
	ShouldExclude(process []api.ProcessInterface) bool
	// ExcludeNamespaces returns a list of namespaces that should be excluded from monitoring.
	ExcludeNamespaces() []string
}

type StaticMonitorFilter struct {
	namespaces         map[string]bool
	clusters           map[string]bool
	originalNamespaces []string
}

func NewStaticMonitorFilter(namespaces, clusters []string) *StaticMonitorFilter {
	return &StaticMonitorFilter{
		namespaces:         convertArrayToMapBool(namespaces),
		clusters:           convertArrayToMapBool(clusters),
		originalNamespaces: namespaces,
	}
}

func (s *StaticMonitorFilter) ShouldExclude(processes []api.ProcessInterface) bool {
	containsNotExcludeCluster := false
	for _, entity := range processes {
		if entity.DetectType() != api.Kubernetes { // for now, we only have the kubernetes detected processes
			continue
		}
		namespace := entity.DetectProcess().(*kubernetes.Process).PodContainer().Pod.Namespace
		if s.namespaces[namespace] {
			return true
		}
		if cluster, _, found := strings.Cut(entity.Entity().ServiceName, "::"); found {
			if !s.clusters[cluster] {
				containsNotExcludeCluster = true
			}
		} else {
			containsNotExcludeCluster = true
			break
		}
	}
	return !containsNotExcludeCluster
}

func (s *StaticMonitorFilter) ExcludeNamespaces() []string {
	return s.originalNamespaces
}

func convertArrayToMapBool(a []string) map[string]bool {
	m := make(map[string]bool, len(a))
	for _, v := range a {
		m[v] = true
	}
	return m
}