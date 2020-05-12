package collector

import (
	"nsxt_exporter/client"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	nsxt "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
)

func init() {
	registerCollector("transport_zone", newTransportZoneCollector)
}

type transportZoneCollector struct {
	transportZoneClient client.TransportZoneClient
	logger              log.Logger

	transportZoneLogicalPort   *prometheus.Desc
	transportZoneTransportNode *prometheus.Desc
}

func newTransportZoneCollector(apiClient *nsxt.APIClient, logger log.Logger) prometheus.Collector {
	nsxtClient := client.NewNSXTClient(apiClient, logger)
	transportZoneLogicalPort := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "transport_zone", "logical_port_total"),
		"Total number of logical port in transport zone",
		[]string{"id", "name"},
		nil,
	)
	transportZoneTransportNode := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "transport_zone", "transport_node_total"),
		"Total number of transport node in transport zone",
		[]string{"id", "name", "status"},
		nil,
	)
	return &transportZoneCollector{
		transportZoneClient:        nsxtClient,
		logger:                     logger,
		transportZoneLogicalPort:   transportZoneLogicalPort,
		transportZoneTransportNode: transportZoneTransportNode,
	}
}

// Describe implements the prometheus.Collector interface.
func (c *transportZoneCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.transportZoneLogicalPort
	ch <- c.transportZoneTransportNode
}

// Collect implements the prometheus.Collector interface.
func (c *transportZoneCollector) Collect(ch chan<- prometheus.Metric) {
	transportZones, err := c.transportZoneClient.ListAllTransportZones()
	if err != nil {
		level.Error(c.logger).Log("msg", "Unable to list transport zones", "err", err)
		return
	}
	c.collectTransportZonesStatus(transportZones, ch)
	c.collectTransportZonesHeatmapStatus(transportZones, ch)
}

func (c *transportZoneCollector) collectTransportZonesStatus(transportZones []manager.TransportZone, ch chan<- prometheus.Metric) {
	for _, transportZone := range transportZones {
		transportZoneStatus, err := c.transportZoneClient.GetTransportZoneStatus(transportZone.Id)
		if err != nil {
			level.Error(c.logger).Log("msg", "Unable to get transport zone status", "id", transportZone.Id, "err", err)
			continue
		}
		transportZoneLabels := []string{transportZone.Id, transportZone.DisplayName}
		ch <- prometheus.MustNewConstMetric(c.transportZoneLogicalPort, prometheus.GaugeValue, float64(transportZoneStatus.NumLogicalPorts), transportZoneLabels...)
	}
}

func (c *transportZoneCollector) collectTransportZonesHeatmapStatus(transportZones []manager.TransportZone, ch chan<- prometheus.Metric) {
	for _, transportZone := range transportZones {
		transportZoneHeatmapStatus, err := c.transportZoneClient.GetHeatmapTransportZoneStatus(transportZone.Id)
		if err != nil {
			level.Error(c.logger).Log("msg", "Unable to get transport zone heatmap status", "id", transportZone.Id, "err", err)
			continue
		}
		ch <- c.constructTransportZoneTransportNodeMetric(transportZone, transportZoneHeatmapStatus.DegradedCount, "degraded")
		ch <- c.constructTransportZoneTransportNodeMetric(transportZone, transportZoneHeatmapStatus.DownCount, "down")
		ch <- c.constructTransportZoneTransportNodeMetric(transportZone, transportZoneHeatmapStatus.UnknownCount, "unknown")
		ch <- c.constructTransportZoneTransportNodeMetric(transportZone, transportZoneHeatmapStatus.UpCount, "up")
	}
}

func (c *transportZoneCollector) constructTransportZoneTransportNodeMetric(transportZone manager.TransportZone, count int32, status string) prometheus.Metric {
	labels := []string{transportZone.Id, transportZone.DisplayName, status}
	return prometheus.MustNewConstMetric(c.transportZoneTransportNode, prometheus.GaugeValue, float64(count), labels...)
}
