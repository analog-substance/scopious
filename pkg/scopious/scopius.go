package scopious

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/analog-substance/scopious/internal/state"
	"github.com/analog-substance/scopious/pkg/utils"
	"golang.org/x/net/publicsuffix"
)

const DefaultScopeDir = "scope"
const DefaultScope = "external"
const scopeFileIPv4 = "ipv4.txt"
const scopeFileIPv6 = "ipv6.txt"
const scopeFileDomains = "domains.txt"
const scopeFileExclude = "exclude.txt"

var ipv6Regexp = regexp.MustCompile("([0-9a-f]{4}::?)+([0-9a-f]{4})")

type Scoper struct {
	Scopes   map[string]*Scope
	ScopeDir string
}

func New() *Scoper {
	return FromPath(DefaultScopeDir)
}

func FromPath(scoperPath string) *Scoper {
	s := &Scoper{
		ScopeDir: scoperPath,
		Scopes:   map[string]*Scope{},
	}

	s.Load()
	return s
}

func (scoper *Scoper) Load() {
	dirs, err := os.ReadDir(scoper.ScopeDir)
	if err != nil {
		err = os.Mkdir(scoper.ScopeDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	for _, dirEntry := range dirs {
		if dirEntry.IsDir() {
			scopeName := dirEntry.Name()
			scoper.Scopes[scopeName] = NewScopeFromPath(filepath.Join(scoper.ScopeDir, scopeName))
			scoper.Scopes[scopeName].Load()
		}
	}

	if len(scoper.Scopes) == 0 {
		// maybe error instead
		scoper.GetScope(DefaultScope)
	}
}
func (scoper *Scoper) Save() {
	for _, scope := range scoper.Scopes {
		scope.Save()
	}
}

func (scoper *Scoper) GetScope(scopeName string) *Scope {
	scope, exists := scoper.Scopes[scopeName]
	if exists {
		return scope
	}

	err := os.Mkdir(filepath.Join(scoper.ScopeDir, scopeName), 0755)
	if err != nil {
		panic(err)
	}

	scoper.Scopes[scopeName] = NewScopeFromPath(filepath.Join(scoper.ScopeDir, scopeName))
	return scoper.Scopes[scopeName]
}

type Scope struct {
	Path              string
	Description       string
	IPv4              map[string]bool
	Domains           map[string]bool
	IPv6              map[string]bool
	Excludes          map[string]bool
	excludedCIDRs     []*net.IPNet
	excludedIPAddrs   []string
	excludedHostnames []string
}

func NewScopeFromPath(path string) *Scope {
	return &Scope{
		Path:              path,
		IPv4:              map[string]bool{},
		IPv6:              map[string]bool{},
		Domains:           map[string]bool{},
		Excludes:          map[string]bool{},
		excludedCIDRs:     []*net.IPNet{},
		excludedIPAddrs:   []string{},
		excludedHostnames: []string{},
	}
}

func (s *Scope) Load() {
	dirs, err := os.ReadDir(s.Path)
	if err != nil {
		panic(err)
	}

	for _, dirEntry := range dirs {
		if !dirEntry.IsDir() {
			if dirEntry.Name() == scopeFileIPv4 {
				s.IPv4, err = utils.ReadLinesMap(filepath.Join(s.Path, scopeFileIPv4))
				if err != nil {
					panic(err)
				}
			}

			if dirEntry.Name() == scopeFileIPv6 {
				s.IPv6, err = utils.ReadLinesMap(filepath.Join(s.Path, scopeFileIPv6))
				if err != nil {
					panic(err)
				}
			}

			if dirEntry.Name() == scopeFileDomains {
				s.Domains, err = utils.ReadLinesMap(filepath.Join(s.Path, scopeFileDomains))
				if err != nil {
					panic(err)
				}
			}

			if dirEntry.Name() == scopeFileExclude {
				s.Excludes, err = utils.ReadLinesMap(filepath.Join(s.Path, scopeFileExclude))
				if err != nil {
					panic(err)
				}
				s.excludedCIDRs, s.excludedIPAddrs, s.excludedHostnames = getCIDRsIPsHostname(s.Excludes)
			}
		}
	}
}

func (s *Scope) Save() {
	err := utils.WriteLines(filepath.Join(s.Path, scopeFileIPv4), sortedScopeKeys(s.IPv4))
	if err != nil && state.Debug {
		log.Println("error saving IPv4:", err)
	}

	err = utils.WriteLines(filepath.Join(s.Path, scopeFileIPv6), sortedScopeKeys(s.IPv6))
	if err != nil && state.Debug {
		log.Println("error saving IPv6:", err)
	}

	err = utils.WriteLines(filepath.Join(s.Path, scopeFileDomains), sortedScopeKeys(s.Domains))
	if err != nil && state.Debug {
		log.Println("error saving Domains:", err)
	}

	err = utils.WriteLines(filepath.Join(s.Path, scopeFileExclude), sortedScopeKeys(s.Excludes))
	if err != nil && state.Debug {
		log.Println("error saving Excludes:", err)
	}
}

func (s *Scope) Add(all bool, scopeItems ...string) {

	for _, scopeItem := range scopeItems {
		scopeItem = normalizedScope(scopeItem)
		if scopeItem == "" {
			continue
		}

		// if we have a direct match in our excludes, do not add the scope item
		_, exists := s.Excludes[scopeItem]
		if exists {
			continue
		}

		if strings.Contains(scopeItem, "/") {
			// perhaps we have a CIDR
			_, err := utils.GetAllIPs(scopeItem, all)
			if err != nil {
				if state.Debug {
					log.Println("error processing cidr", err)
				}
				continue
			}

			// if we have a `:` then we must have an IPv6 address
			if strings.Contains(scopeItem, ":") {
				s.IPv6[scopeItem] = true
				continue
			}

			s.IPv4[scopeItem] = true
			continue
		}

		// if we have a `:` then we must have an IPv6 address
		if strings.Contains(scopeItem, ":") {
			s.IPv6[scopeItem] = true
			continue
		}

		ip := net.ParseIP(scopeItem)
		if ip != nil {
			if s.CanAddIP(ip) {
				s.IPv4[ip.String()] = true
			}
			// item was an IP address, continue now to prevent useless processing
			continue
		}

		// not IPv6 or IPv4... must be a domain
		if s.CanAddDomain(scopeItem) {
			s.Domains[scopeItem] = true
		}
	}
}

func (s *Scope) AddExclude(scopeItems ...string) {

	for _, scopeItem := range scopeItems {
		scopeItem = normalizedScope(scopeItem)
		if scopeItem == "" {
			continue
		}

		s.Excludes[scopeItem] = true
	}
}

func (s *Scope) Prune(all bool, scopeItemsToCheck ...string) []string {

	scopeCheckResults := map[string]bool{}

	includeIPv4CIDRs, includeIPv4Addrs, _ := getCIDRsIPsHostname(s.IPv4)
	includeIPv6CIDRs, includeIPv6Addrs, _ := getCIDRsIPsHostname(s.IPv6)
	_, _, includeDomains := getCIDRsIPsHostname(s.Domains)

	includeCIDRs := append(includeIPv4CIDRs, includeIPv6CIDRs...)
	includeIPAddrs := append(includeIPv4Addrs, includeIPv6Addrs...)
	expandedIPs, ipAddrsToCheck, domainStrsToCheck := normalizeAndExpandStringSlice(scopeItemsToCheck, all)

	scopeItemsToCheck = append(scopeItemsToCheck, expandedIPs...)
CheckIpAddr:
	for _, ipAddr := range ipAddrsToCheck {
		if !s.CanAddIP(ipAddr) {
			continue CheckIpAddr
		}

		// is IP explicitly allowed
		if slices.Contains(includeIPAddrs, ipAddr.String()) {
			scopeCheckResults[ipAddr.String()] = true
			continue CheckIpAddr
		}

		// is IP implicitly allowed through CIDR
		for _, includedCIDR := range includeCIDRs {
			if includedCIDR.Contains(ipAddr) {
				scopeCheckResults[ipAddr.String()] = true
				continue CheckIpAddr
			}
		}
	}

CheckDomains:
	for _, domainToCheck := range domainStrsToCheck {
		if !s.CanAddDomain(domainToCheck) {
			continue CheckDomains
		}

		// is IP explicitly allowed
		if slices.Contains(includeDomains, domainToCheck) {
			scopeCheckResults[domainToCheck] = true
			continue CheckDomains
		}

		// is IP implicitly allowed through CIDR
		for _, includedDomain := range includeDomains {
			if strings.HasSuffix(domainToCheck, "."+includedDomain) {
				scopeCheckResults[domainToCheck] = true
				continue CheckDomains
			}
		}
	}

	prunedResultsMap := map[string]bool{}
	for inScopeItem, _ := range scopeCheckResults {
		for _, originalInput := range scopeItemsToCheck {
			normal := normalizedScope(originalInput)

			if strings.Contains(normal, ":") {
				// must be ipV6
				parsedIP := net.ParseIP(normal)
				if parsedIP != nil && parsedIP.String() == inScopeItem {
					prunedResultsMap[originalInput] = true
					continue
				}
			}

			if normal == inScopeItem {
				prunedResultsMap[originalInput] = true
				continue
			}

		}
	}

	prunedResults := []string{}
	for prunedRes, _ := range prunedResultsMap {
		prunedResults = append(prunedResults, prunedRes)
	}
	sort.Strings(prunedResults)
	return prunedResults
}

func (s *Scope) AllExpanded(all bool) []string {
	return s.Prune(all, s.AllIPs()...)
}

func (s *Scope) AllIPs() []string {
	return append(sortedScopeKeys(s.IPv4), sortedScopeKeys(s.IPv6)...)
}

func (s *Scope) RootDomains() []string {
	rootDomains := map[string]bool{}
	for domain, _ := range s.Domains {
		rootDomain, err := publicsuffix.EffectiveTLDPlusOne(domain)
		if err != nil {
			log.Println("root domain err", err)
		}
		rootDomains[rootDomain] = true
	}
	return sortedScopeKeys(rootDomains)
}

func (s *Scope) AllDomains() []string {
	return sortedScopeKeys(s.Domains)
}
func (s *Scope) CanAddIP(ipAddr net.IP) bool {
	// is IP explicitly blocked
	if slices.Contains(s.excludedIPAddrs, ipAddr.String()) {
		return false
	}

	// is IP implicitly blocked through CIDR
	for _, excludedCidr := range s.excludedCIDRs {
		if excludedCidr.Contains(ipAddr) {
			return false
		}
	}

	return true
}

func (s *Scope) CanAddDomain(domainToCheck string) bool {
	// is domain explicitly blocked
	if slices.Contains(s.excludedHostnames, domainToCheck) {
		return false
	}

	// is domain implicitly blocked via parent domain
	for _, excludedDomain := range s.excludedHostnames {
		if strings.HasSuffix(domainToCheck, "."+excludedDomain) {
			return false
		}
	}

	return true
}

func getExpandedCIDRSFromScopeMap(scopeMap map[string]bool, all bool) []net.IP {
	allScopeIPs := []net.IP{}
	for inScopeStr, _ := range scopeMap {
		inScopeIPaddrs, err := utils.GetAllIPs(inScopeStr, all)
		if err == nil {
			allScopeIPs = append(allScopeIPs, inScopeIPaddrs...)
		}
	}
	return allScopeIPs
}

func normalizedScope(scopeItem string) string {
	scopeItem = strings.TrimSpace(scopeItem)
	if len(scopeItem) == 0 {
		return ""
	}

	containsProto := strings.Contains(scopeItem, "://")
	if !containsProto && strings.Contains(scopeItem, "/") {
		// perhaps we have a CIDR
		_, ipNet, err := net.ParseCIDR(scopeItem)
		if err == nil {
			return ipNet.String()
		}
	}

	possibleIPv6 := ipv6Regexp.MatchString(scopeItem)
	if !containsProto {
		if possibleIPv6 {
			scopeItem = fmt.Sprintf("[%s]", scopeItem)
		}
		scopeItem = fmt.Sprintf("https://%s", scopeItem)
	}

	// we may have a URL
	parsedURL, err := url.Parse("a" + scopeItem)
	if err == nil {
		// no errors, we have a URL
		if len(parsedURL.Host) > 0 {
			return strings.TrimSuffix(parsedURL.Hostname(), ".")
		}
	}

	// must be invalid
	return ""
}

func normalizeAndExpandStringSlice(scopeItemsToCheck []string, all bool) (expandedIPs []string, normalizedIPAddrs []net.IP, normalizedHostnames []string) {

	for _, scopeToCheck := range scopeItemsToCheck {
		normalized := normalizedScope(scopeToCheck)
		if normalized == "" {
			continue
		}
		ipAddrs, err := utils.GetAllIPs(normalized, all)
		if err == nil {
			if strings.Contains(normalized, "/") {
				for _, ip := range ipAddrs {
					expandedIPs = append(expandedIPs, ip.String())
				}
			}
			normalizedIPAddrs = append(normalizedIPAddrs, ipAddrs...)
		} else {
			normalizedHostnames = append(normalizedHostnames, normalized)
		}
	}
	return
}

func getCIDRsIPsHostname(scopeMap map[string]bool) (cidrs []*net.IPNet, ipAddrs []string, hostnames []string) {

	for scopeItem, _ := range scopeMap {
		_, ipNet, err := net.ParseCIDR(scopeItem)
		if err == nil {
			cidrs = append(cidrs, ipNet)
			continue
		}

		ip := net.ParseIP(scopeItem)
		if ip != nil {
			ipAddrs = append(ipAddrs, ip.String())
			continue
		}

		hostnames = append(hostnames, scopeItem)
	}
	return
}

func sortedScopeKeys(mapWithStringKeys map[string]bool) []string {
	keys := []string{}
	for key, _ := range mapWithStringKeys {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
