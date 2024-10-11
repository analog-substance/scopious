package scopious

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/analog-substance/scopious/pkg/state"
	"github.com/analog-substance/scopious/pkg/utils"
	"golang.org/x/net/publicsuffix"
)

const DefaultScopeDir = "data"
const DefaultScope = "default"
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

	scopeFile := scoper.GetScopePath(scopeName)
	err := os.Mkdir(scopeFile, 0755)
	if err != nil {
		panic(err)
	}

	scoper.Scopes[scopeName] = NewScopeFromPath(scopeFile)
	return scoper.Scopes[scopeName]
}

func (scoper *Scoper) GetScopePath(scopeName string) string {
	return filepath.Join(scoper.ScopeDir, scopeName)
}

func (scoper *Scoper) GetScopeExcludePath(scopeName string) string {
	return filepath.Join(scoper.ScopeDir, scopeName, scopeFileExclude)
}

func (scoper *Scoper) GetScopeIPv4Path(scopeName string) string {
	return filepath.Join(scoper.ScopeDir, scopeName, scopeFileIPv4)
}

func (scoper *Scoper) GetScopeIPv6Path(scopeName string) string {
	return filepath.Join(scoper.ScopeDir, scopeName, scopeFileIPv6)
}

func (scoper *Scoper) GetScopeDomainsPath(scopeName string) string {
	return filepath.Join(scoper.ScopeDir, scopeName, scopeFileDomains)
}

type Scope struct {
	Path              string
	Description       string
	IPv4              map[string]bool
	Domains           map[string]bool
	IPv6              map[string]bool
	Excludes          map[string]bool
	inScopeCIDRs      map[string]*net.IPNet
	excludedCIDRs     map[string]*net.IPNet
	excludedIPAddrs   map[string]bool
	excludedHostnames map[string]bool
	rootDomainMap     map[string]bool
	rootDomainSorted  []string
}

func NewScopeFromPath(path string) *Scope {
	return &Scope{
		Path:     path,
		IPv4:     map[string]bool{},
		IPv6:     map[string]bool{},
		Domains:  map[string]bool{},
		Excludes: map[string]bool{},

		rootDomainMap:    map[string]bool{},
		rootDomainSorted: []string{},
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
			}
		}
	}

	s.populateExcludes()
	s.populateIncludes()
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
			if s.CanAddIP(&ip) {
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

func (s *Scope) IsIPInScope(ip *net.IP, mustBeInScope bool) bool {
	if ip == nil {
		return false
	}

	if s.excludedCIDRs == nil {
		s.populateExcludes()
	}

	// exclude takes precedence
	for _, excludeCIDR := range s.excludedCIDRs {
		if excludeCIDR.Contains(*ip) {
			return false
		}
	}
	_, ok := s.excludedIPAddrs[ip.String()]
	if ok {
		return false
	}

	if !mustBeInScope {
		return true
	}

	_, ok = s.IPv4[ip.String()]
	if ok {
		return true
	}

	if s.inScopeCIDRs == nil {
		s.populateIncludes()
	}

	for _, inScopeCIDR := range s.inScopeCIDRs {
		if inScopeCIDR.Contains(*ip) {
			return true
		}
	}
	return false
}

func (s *Scope) getExcludedHostNames() map[string]bool {

	return s.excludedHostnames
}

func (s *Scope) IsDomainInScope(domain string, mustBeInScope bool) bool {
	if domain == "" {
		return false
	}

	if s.excludedHostnames == nil {
		s.populateExcludes()
	}

	_, ok := s.excludedHostnames[domain]
	if ok {
		return false
	}
	// is domain implicitly blocked via parent domain
	for excludedDomain := range s.excludedHostnames {
		if strings.HasSuffix(domain, "."+excludedDomain) {
			return false
		}
	}

	if !mustBeInScope {
		return true
	}

	_, ok = s.Domains[domain]
	if ok {
		return true
	}
	// is domain implicitly blocked via parent domain
	for _, includedDomain := range s.RootDomains() {
		if strings.HasSuffix(domain, "."+includedDomain) {
			return true
		}
	}

	return false
}

func (s *Scope) Prune(all bool, scopeItemsToCheck ...string) []string {
	scopeCheckResults := map[string]bool{}

	for _, scopeToCheck := range scopeItemsToCheck {
		normalized := normalizedScope(scopeToCheck)
		if normalized == "" {
			continue
		}
		ipAddrs, err := utils.GetAllIPs(normalized, all)
		if err == nil {
			for _, expandedIP := range ipAddrs {
				if s.IsIPInScope(expandedIP, true) {
					scopeCheckResults[expandedIP.String()] = true
				}
			}
		} else {
			if s.IsInScope(normalized) {
				scopeCheckResults[scopeToCheck] = true
			}
		}
	}

	prunedResults := []string{}
	for prunedRes := range scopeCheckResults {
		prunedResults = append(prunedResults, prunedRes)
	}
	return prunedResults
}

func (s *Scope) IsInScope(itemToCheck string) bool {
	normalized := normalizedScope(itemToCheck)
	ip := net.ParseIP(normalized)
	if ip != nil {
		return s.IsIPInScope(&ip, true)
	}

	ips, err := utils.GetAllIPs(normalized, true)
	if err == nil {
		for _, ip := range ips {
			if !s.IsIPInScope(ip, true) {
				return false
			}
		}
		return true
	}

	return s.IsDomainInScope(itemToCheck, true)
}

func (s *Scope) AllExpanded(all bool) []string {
	return s.Prune(all, s.AllIPs()...)
}

func (s *Scope) AllIPs() []string {
	return append(sortedScopeKeys(s.IPv4), sortedScopeKeys(s.IPv6)...)
}

func (s *Scope) RootDomains() []string {
	if len(s.rootDomainMap) == 0 {
		s.rootDomainMap = make(map[string]bool)
		for domain := range s.Domains {
			rootDomain, err := publicsuffix.EffectiveTLDPlusOne(domain)
			if err != nil {
				log.Println("root domain err", err)
			}
			s.rootDomainMap[rootDomain] = true
		}
		s.rootDomainSorted = sortedScopeKeys(s.rootDomainMap)
	}
	return s.rootDomainSorted
}

func (s *Scope) AllDomains() []string {
	return sortedScopeKeys(s.Domains)
}
func (s *Scope) CanAddIP(ipAddr *net.IP) bool {
	// is IP explicitly blocked
	return s.IsIPInScope(ipAddr, false)
}

func (s *Scope) CanAddDomain(domainToCheck string) bool {
	// is domain explicitly blocked
	return s.IsDomainInScope(domainToCheck, false)
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
			hostname := strings.TrimSuffix(parsedURL.Hostname(), ".")
			return hostname
			//_, err := publicsuffix.EffectiveTLDPlusOne(hostname)
			//if err == nil {
			//	if state.Debug {
			//		log.Println("hostname", hostname)
			//	}
			//	return hostname
			//} else {
			//	log.Println("root domain err", err)
			//}
		}
	}

	// must be invalid
	return ""
}

func normalizeAndExpandStringSlice(scopeItemsToCheck []string, all bool) (expandedIPs []string, normalizedIPAddrs []*net.IP, normalizedHostnames []string) {

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

func (s *Scope) populateIncludes() {
	s.inScopeCIDRs = map[string]*net.IPNet{}
	s.populateInScopeCIDRs(s.IPv4)
	s.populateInScopeCIDRs(s.IPv6)
}

func (s *Scope) populateExcludes() {

	s.excludedCIDRs = map[string]*net.IPNet{}
	s.excludedIPAddrs = map[string]bool{}
	s.excludedHostnames = map[string]bool{}

	for scopeItem := range s.Excludes {
		_, ok := s.excludedCIDRs[scopeItem]
		if !ok {
			_, ipNet, err := net.ParseCIDR(scopeItem)
			if err == nil {
				s.excludedCIDRs[scopeItem] = ipNet
				continue
			}
		}

		_, ok = s.excludedIPAddrs[scopeItem]
		if !ok {
			ip := net.ParseIP(scopeItem)
			if ip != nil {
				s.excludedIPAddrs[ip.String()] = true
				continue
			}
		}

		s.excludedHostnames[scopeItem] = true
	}
}

func getCIDRsIPsHostname(scopeMap map[string]bool) (cidrs []*net.IPNet, ipAddrs []string, hostnames []string) {

	for scopeItem := range scopeMap {
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
	for key := range mapWithStringKeys {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (s *Scope) populateInScopeCIDRs(ipScopeMap map[string]bool) {

	for ip := range ipScopeMap {
		if strings.Contains(ip, "/") {
			_, ok := s.inScopeCIDRs[ip]
			if !ok {
				_, ipNet, err := net.ParseCIDR(ip)
				if err != nil {
					panic(err)
				}

				s.inScopeCIDRs[ip] = ipNet
			}
		}
	}
}
