// SPDX-License-Identifier: Apache-2.0
//
// Copyright The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build windows

package hyperv

import (
	"fmt"

	"github.com/prometheus-community/windows_exporter/internal/pdh"
	"github.com/prometheus-community/windows_exporter/internal/types"
	"github.com/prometheus/client_golang/prometheus"
)

type collectorReplicaVM struct {
	perfDataCollectorReplicaVM *pdh.Collector
	perfDataObjectReplicaVM    []perfDataCounterValuesReplicaVM

	// \Hyper-V Replica VM(*)\Last Replication Size
	replicaVMLastReplicationSizeBytes *prometheus.Desc
}

type perfDataCounterValuesReplicaVM struct {
	Name string

	ReplicaVMLastReplicationSizeBytes float64 `perfdata:"Last Replication Size"`
}

func (c *Collector) buildReplicaVM() error {
	var err error

	c.perfDataCollectorReplicaVM, err = pdh.NewCollector[perfDataCounterValuesReplicaVM](c.logger, pdh.CounterTypeRaw, "Hyper-V Replica VM", pdh.InstancesAll)
	if err != nil {
		return fmt.Errorf("failed to create Hyper-V Replica VM collector: %w", err)
	}

	c.replicaVMLastReplicationSizeBytes = prometheus.NewDesc(
		prometheus.BuildFQName(types.Namespace, Name, "replica_vm_last_replication_size_bytes"),
		"Represents the size of the last replication in bytes",
		[]string{"vm"},
		nil,
	)

	return nil
}

func (c *Collector) collectReplicaVM(ch chan<- prometheus.Metric) error {
	err := c.perfDataCollectorReplicaVM.Collect(&c.perfDataObjectReplicaVM)
	if err != nil {
		return fmt.Errorf("failed to collect Hyper-V Replica VM metrics: %w", err)
	}

	for _, data := range c.perfDataObjectReplicaVM {
		ch <- prometheus.MustNewConstMetric(
			c.replicaVMLastReplicationSizeBytes,
			prometheus.GaugeValue,
			data.ReplicaVMLastReplicationSizeBytes,
			data.Name,
		)
	}

	return nil
}
