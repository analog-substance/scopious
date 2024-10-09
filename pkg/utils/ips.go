package utils

import "net"

func GetAllIPs(cidr string, all bool) ([]*net.IP, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	allIPs := []*net.IP{}
	if err != nil {
		ip = net.ParseIP(cidr)
		if ip != nil {
			allIPs = append(allIPs, &ip)
			return allIPs, nil
		}
		return nil, err
	}

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		parsedIp := net.ParseIP(ip.String())
		allIPs = append(allIPs, &parsedIp)
	}
	if len(allIPs) == 1 {
		return allIPs, nil
	}

	if all {
		return allIPs, nil
	}
	// else remove network address and broadcast address
	return allIPs[1 : len(allIPs)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
