package utils

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"os"
	"sort"
)

/*
GetPhysicalInterface identifies the primary physical network interface by examining the system's routing table.
It checks both IPv4 and IPv6 routes, prioritizing interfaces based on route metrics, and validates that the interface is physical by checking for the presence of a device entry in /sys/class/net/<interface>/device.

Returns the name of the physical network interface or an error if no valid interface is found.
*/
func GetPhysicalInterface() (string, error) {
	// 1. Define the address families to check: first V4, then V6
	families := []int{netlink.FAMILY_V4, netlink.FAMILY_V6}

	for _, family := range families {
		routes, err := netlink.RouteListFiltered(family, &netlink.Route{Dst: nil}, netlink.RT_FILTER_DST)
		if err != nil || len(routes) == 0 {
			continue
		}

		// Sort routes by priority (lower value means higher priority)
		sort.Slice(routes, func(i, j int) bool {
			return routes[i].Priority < routes[j].Priority
		})

		// Iterate through the sorted routes and find the first valid physical interface
		for _, route := range routes {
			link, _ := netlink.LinkByIndex(route.LinkIndex)
			if link == nil {
				continue
			}

			// Validate (check /sys/class/net/<name>/device)
			if _, err := os.Stat(fmt.Sprintf("/sys/class/net/%s/device", link.Attrs().Name)); err == nil {
				return link.Attrs().Name, nil
			}
		}
	}

	return "", fmt.Errorf("No valid physical network interface found.")
}
