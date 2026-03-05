package network

import (
	"fmt"
	"log"
	"realess-server/pkg/utils"
)

/*
SetupNatMasquerade configures NAT masquerading for the specified TUN interface using iptables.
Returns an error if the operation fails.
*/
func SetupNatMasquerade(cidr string, physicalInterface string) error {
	// Enable NAT masquerading using iptables
	checkCmd := fmt.Sprintf("iptables -t nat -C POSTROUTING -s %s -o %s -j MASQUERADE", cidr, physicalInterface)
	if err := utils.RunAsRootSilent(checkCmd); err == nil {
		log.Println("[network] NAT masquerading rule configured.")
		return nil
	}

	// If the rule does not exist, add it
	addCmd := fmt.Sprintf("iptables -t nat -A POSTROUTING -s %s -o %s -j MASQUERADE", cidr, physicalInterface)
	if err := utils.RunAsRoot(addCmd); err != nil {
		log.Fatalf("Failed to set up NAT masquerading: %v", err)
	}

	log.Println("[network] NAT masquerading configured.")
	return nil
}

/*
RemoveNatMasquerade removes the NAT masquerading rule for the specified CIDR.
Returns an error if the operation fails.
*/
func RemoveNatMasquerade(cidr string, physicalInterface string) error {
	// Remove the NAT masquerading rule using iptables
	checkCmd := fmt.Sprintf("iptables -t nat -C POSTROUTING -s %s -o %s -j MASQUERADE", cidr, physicalInterface)
	if err := utils.RunAsRootSilent(checkCmd); err != nil {
		return nil
	}

	delCmd := fmt.Sprintf("iptables -t nat -D POSTROUTING -s %s -o %s -j MASQUERADE", cidr, physicalInterface)
	if err := utils.RunAsRoot(delCmd); err != nil {
		return fmt.Errorf("Failed to remove NAT masquerading: %w", err)
	}

	log.Println("[network] NAT masquerading rule removed.")
	return nil
}
