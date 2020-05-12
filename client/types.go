package client

import "github.com/vmware/go-vmware-nsxt/manager"

// LogicalPortClient represents API group logical port for NSX-T client.
type LogicalPortClient interface {
	ListLogicalPorts(localVarOptionals map[string]interface{}) (manager.LogicalPortListResult, error)
	GetLogicalPortStatusSummary(localVarOptionals map[string]interface{}) (manager.LogicalPortStatusSummary, error)
	GetLogicalPortOperationalStatus(lportId string, localVarOptionals map[string]interface{}) (manager.LogicalPortOperationalStatus, error)
}

// LogicalRouterPortClient represents API group logical router port for NSX-T client.
type LogicalRouterPortClient interface {
	ListLogicalRouterPorts(localVarOptionals map[string]interface{}) (manager.LogicalRouterPortListResult, error)
	GetLogicalRouterPortStatisticsSummary(lrportID string) (manager.LogicalRouterPortStatisticsSummary, error)
}

// DHCPClient represents API group DHCP for NSX-T client.
type DHCPClient interface {
	ListDhcpServers(localVarOptionals map[string]interface{}) (manager.LogicalDhcpServerListResult, error)
	GetDhcpStatus(dhcpID string, localVarOptionals map[string]interface{}) (manager.DhcpServerStatus, error)
}
