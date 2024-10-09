package scopious

import (
	"reflect"
	"slices"
	"testing"
)

func TestScope_Add_IPv4(t *testing.T) {
	type fields struct {
		Path        string
		Description string
		IPv4        map[string]bool
		Domains     map[string]bool
		IPv6        map[string]bool
		Exclude     map[string]bool
	}

	emptyScope := fields{
		IPv4:    map[string]bool{},
		Domains: map[string]bool{},
		IPv6:    map[string]bool{},
		Exclude: map[string]bool{},
	}

	emptyScope2 := fields{
		IPv4:    map[string]bool{},
		Domains: map[string]bool{},
		IPv6:    map[string]bool{},
		Exclude: map[string]bool{},
	}

	emptyScope3 := fields{
		IPv4:    map[string]bool{},
		Domains: map[string]bool{},
		IPv6:    map[string]bool{},
		Exclude: map[string]bool{},
	}

	type args struct {
		scopeItems []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]bool
	}{
		{name: "Should add IPv4 to IPv4", fields: emptyScope, args: args{[]string{"10.0.0.1"}}, want: map[string]bool{"10.0.0.1": true}},
		{name: "Should not add domain to IPv4", fields: emptyScope2, args: args{[]string{"10.0.0.1.xip.io"}}, want: map[string]bool{}},
		{name: "Should not add IPv6 to IPv4", fields: emptyScope3, args: args{[]string{"fda4:20e2:424d:cad4:e96:4b42:2fe1:46fb"}}, want: map[string]bool{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scope{
				Path:        tt.fields.Path,
				Description: tt.fields.Description,
				IPv4:        tt.fields.IPv4,
				Domains:     tt.fields.Domains,
				IPv6:        tt.fields.IPv6,
				Excludes:    tt.fields.Exclude,
			}
			s.Add(false, tt.args.scopeItems...)

			if got := s.IPv4; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IPv4 = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScope_Add_Domains(t *testing.T) {
	type fields struct {
		Path        string
		Description string
		IPv4        map[string]bool
		Domains     map[string]bool
		IPv6        map[string]bool
		Exclude     map[string]bool
	}

	emptyScope := fields{
		IPv4:    map[string]bool{},
		Domains: map[string]bool{},
		IPv6:    map[string]bool{},
		Exclude: map[string]bool{},
	}

	emptyScope2 := fields{
		IPv4:    map[string]bool{},
		Domains: map[string]bool{},
		IPv6:    map[string]bool{},
		Exclude: map[string]bool{},
	}

	emptyScope3 := fields{
		IPv4:    map[string]bool{},
		Domains: map[string]bool{},
		IPv6:    map[string]bool{},
		Exclude: map[string]bool{},
	}

	type args struct {
		scopeItems []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]bool
	}{
		{name: "Should add Domain to Domans", fields: emptyScope, args: args{[]string{"10.0.0.1.xip.io"}}, want: map[string]bool{"10.0.0.1.xip.io": true}},
		{name: "Should not add IPv4 to domains", fields: emptyScope2, args: args{[]string{"10.0.0.1"}}, want: map[string]bool{}},
		{name: "Should not add IPv6 to domains", fields: emptyScope3, args: args{[]string{"fda4:20e2:424d:cad4:e96:4b42:2fe1:46fb"}}, want: map[string]bool{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scope{
				Path:        tt.fields.Path,
				Description: tt.fields.Description,
				IPv4:        tt.fields.IPv4,
				Domains:     tt.fields.Domains,
				IPv6:        tt.fields.IPv6,
				Excludes:    tt.fields.Exclude,
			}
			s.Add(false, tt.args.scopeItems...)

			if got := s.Domains; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Domains = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScope_Add_IPv6(t *testing.T) {
	type fields struct {
		Path        string
		Description string
		IPv4        map[string]bool
		Domains     map[string]bool
		IPv6        map[string]bool
		Exclude     map[string]bool
	}

	emptyScope := fields{
		IPv4:    map[string]bool{},
		Domains: map[string]bool{},
		IPv6:    map[string]bool{},
		Exclude: map[string]bool{},
	}

	emptyScope2 := fields{
		IPv4:    map[string]bool{},
		Domains: map[string]bool{},
		IPv6:    map[string]bool{},
		Exclude: map[string]bool{},
	}

	emptyScope3 := fields{
		IPv4:    map[string]bool{},
		Domains: map[string]bool{},
		IPv6:    map[string]bool{},
		Exclude: map[string]bool{},
	}

	type args struct {
		scopeItems []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]bool
	}{
		{name: "Should add IPv6 to IPv6", fields: emptyScope, args: args{[]string{"fda4:20e2:424d:cad4:e96:4b42:2fe1:46fb"}}, want: map[string]bool{"fda4:20e2:424d:cad4:e96:4b42:2fe1:46fb": true}},
		{name: "Should not add IPv4 to IPv6", fields: emptyScope2, args: args{[]string{"10.0.0.1"}}, want: map[string]bool{}},
		{name: "Should not add Domain to IPv6", fields: emptyScope3, args: args{[]string{"10.0.0.1.xip.io"}}, want: map[string]bool{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scope{
				Path:        tt.fields.Path,
				Description: tt.fields.Description,
				IPv4:        tt.fields.IPv4,
				Domains:     tt.fields.Domains,
				IPv6:        tt.fields.IPv6,
				Excludes:    tt.fields.Exclude,
			}
			s.Add(false, tt.args.scopeItems...)

			if got := s.IPv6; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IPv6 = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalizedScope(t *testing.T) {
	type args struct {
		scopeItem string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Plain IPv4 Address", args: args{"10.0.0.1"}, want: "10.0.0.1"},
		{name: "Plain IPv4 Address with space", args: args{"  10.0.0.1  "}, want: "10.0.0.1"},
		{name: "Plain CIDR", args: args{"10.0.0.1/24"}, want: "10.0.0.0/24"},
		{name: "Plain domain", args: args{"whatever.dead"}, want: "whatever.dead"},
		{name: "URL", args: args{"https://whatever.dead"}, want: "whatever.dead"},
		{name: "URL with params", args: args{"https://whatever.dead/place/?t=test.com"}, want: "whatever.dead"},
		{name: "Plain IPv6", args: args{"fda4:20e2:424d:cad4:e96:4b42:2fe1:46fb"}, want: "fda4:20e2:424d:cad4:e96:4b42:2fe1:46fb"},
		{name: "Plain IPv6 CIDR", args: args{"fda4:20e2:424d:cad4:e96:4b42:2fe1:46fb/62"}, want: "fda4:20e2:424d:cad4::/62"},
		{name: "domain with port", args: args{"asdf.com.test:443"}, want: "asdf.com.test"},
		{name: "URL with port", args: args{"https://asdf.com.test:443/face"}, want: "asdf.com.test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizedScope(tt.args.scopeItem); got != tt.want {
				t.Errorf("normalizedScope() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScope_Prune(t *testing.T) {
	type fields struct {
		Path        string
		Description string
		IPv4        map[string]bool
		Domains     map[string]bool
		IPv6        map[string]bool
		Excludes    map[string]bool
	}
	type args struct {
		scopeItemsToCheck []string
	}

	scopeFields := fields{
		IPv4: map[string]bool{
			"10.42.0.0/30": true,
			"10.42.2.42":   true,
		},
		Domains: map[string]bool{
			"inscope.tld":        true,
			"alsovalidscope.tld": true,
		},
		IPv6: map[string]bool{
			"2001:db8::/64": true,
		},
		Excludes: map[string]bool{
			"10.42.0.0/31":           true,
			"notinscope.inscope.tld": true,
		},
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name:   "kitchen sink check",
			fields: scopeFields,
			args: args{
				[]string{
					"10.13.0.2",
					"10.29.0.0/29",
					"10.42.0.42",
					"console.inscope.tld",
					"10.42.0.10",
					"https://admin.stillinscope.inscope.tld:8443/garbage.html",
					"api.notinscope.inscope.tld",
					"api.stillinscope.inscope.tld",
					"10.42.0.0/24",
					"2001:0db8:0000:0000:0000:ff00:0042:8329",
				},
			},
			want: []string{
				"10.42.0.2",
				"10.42.0.3",
				"console.inscope.tld",
				"api.stillinscope.inscope.tld",
				"2001:db8::ff00:42:8329",
				"https://admin.stillinscope.inscope.tld:8443/garbage.html",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scope{
				Path:        tt.fields.Path,
				Description: tt.fields.Description,
				IPv4:        tt.fields.IPv4,
				Domains:     tt.fields.Domains,
				IPv6:        tt.fields.IPv6,
				Excludes:    tt.fields.Excludes,
			}
			got := s.Prune(false, tt.args.scopeItemsToCheck...)

			if len(got) != len(tt.want) {
				t.Errorf("Prune returned wrong length: got:%v, wanted:%v", got, tt.want)
			}

			didntContain := []string{}
			for _, want := range tt.want {
				if !slices.Contains(got, want) {
					didntContain = append(didntContain, want)
				}
			}

			if len(didntContain) > 0 {
				t.Errorf("Prune() didnt contain desired values = %v, want %v", got, didntContain)
			}
		})
	}
}

func TestScope_AllIPs(t *testing.T) {
	type fields struct {
		Path        string
		Description string
		IPv4        map[string]bool
		Domains     map[string]bool
		IPv6        map[string]bool
		Excludes    map[string]bool
	}

	scopeFields := fields{
		IPv4: map[string]bool{
			"10.42.0.0/24": true,
			"10.42.2.42":   true,
			"10.42.2.43":   true,
		},
		Domains: map[string]bool{
			"inscope.tld":        true,
			"alsovalidscope.tld": true,
		},
		IPv6: map[string]bool{
			"2001:db8::/64": true,
		},
		Excludes: map[string]bool{
			"10.42.0.69":        true,
			"10.42.0.0/28":      true,
			"admin.inscope.tld": true,
		},
	}

	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{name: "returns inscope things", fields: scopeFields, want: []string{"10.42.0.0/24", "10.42.2.42", "10.42.2.43", "2001:db8::/64"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scope{
				Path:        tt.fields.Path,
				Description: tt.fields.Description,
				IPv4:        tt.fields.IPv4,
				Domains:     tt.fields.Domains,
				IPv6:        tt.fields.IPv6,
				Excludes:    tt.fields.Excludes,
			}
			if got := s.AllIPs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllIPs() = %v, want %v", got, tt.want)
			}
		})
	}
}
