package network

import (
	"fmt"
	"log"
	"realess-server/pkg/utils"
)

/*
SetupForwardRules configures firewall rules to allow traffic between the TUN interface and the internet.
Returns an error if the operation fails.
*/
func SetupForwardRules(physicalInterface string, tunName string) error {
	// Allow forwarding from the TUN interface to the internet
	checkCmd1 := fmt.Sprintf("iptables -C FORWARD -i %s -o %s -j ACCEPT", tunName, physicalInterface)
	if err := utils.RunAsRootSilent(checkCmd1); err != nil {
		addCmd1 := fmt.Sprintf("iptables -A FORWARD -i %s -o %s -j ACCEPT", tunName, physicalInterface)
		if err := utils.RunAsRoot(addCmd1); err != nil {
			log.Fatalf("Failed to set up forward rule (TUN to internet): %v", err)
		}
	}
	log.Println("[network] Forward rule (TUN to internet) configured.")

	// Allow forwarding from the internet to the TUN interface
	checkCmd2 := fmt.Sprintf("iptables -C FORWARD -i %s -o %s -m state --state RELATED,ESTABLISHED -j ACCEPT", physicalInterface, tunName)
	if err := utils.RunAsRootSilent(checkCmd2); err != nil {
		addCmd2 := fmt.Sprintf("iptables -A FORWARD -i %s -o %s -m state --state RELATED,ESTABLISHED -j ACCEPT", physicalInterface, tunName)
		if err := utils.RunAsRoot(addCmd2); err != nil {
			log.Fatalf("Failed to set up forward rule (internet to TUN): %v", err)
		}
	}
	log.Println("[network] Forward rule (internet to TUN) configured.")

	return nil
}

/*
RemoveForwardRules removes the firewall rules that allow traffic between the TUN interface and the internet.
Returns an error if the operation fails.
*/
func RemoveForwardRules(physicalInterface string, tunName string) error {
	// Remove forwarding from the TUN interface to the internet
	checkCmd1 := fmt.Sprintf("iptables -C FORWARD -i %s -o %s -j ACCEPT", tunName, physicalInterface)
	if err := utils.RunAsRootSilent(checkCmd1); err == nil {
		delCmd1 := fmt.Sprintf("iptables -D FORWARD -i %s -o %s -j ACCEPT", tunName, physicalInterface)
		if err := utils.RunAsRoot(delCmd1); err != nil {
			return fmt.Errorf("Failed to remove forward rule (TUN to internet): %w", err)
		}
		log.Println("[network] Forward rule (TUN to internet) removed.")
	}

	// Remove forwarding from the internet to the TUN interface
	checkCmd2 := fmt.Sprintf("iptables -C FORWARD -i %s -o %s -m state --state RELATED,ESTABLISHED -j ACCEPT", physicalInterface, tunName)
	if err := utils.RunAsRootSilent(checkCmd2); err == nil {
		delCmd2 := fmt.Sprintf("iptables -D FORWARD -i %s -o %s -m state --state RELATED,ESTABLISHED -j ACCEPT", physicalInterface, tunName)
		if err := utils.RunAsRoot(delCmd2); err != nil {
			return fmt.Errorf("Failed to remove forward rule (internet to TUN): %w", err)
		}
		log.Println("[network] Forward rule (internet to TUN) removed.")
	}

	return nil
}
