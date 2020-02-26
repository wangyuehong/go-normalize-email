package gonormail

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		email string
		want  string
	}{
		{email: "Not A Email", want: "Not A Email"},
		{email: "Not@A@Email", want: "Not@A@Email"},
		{email: "A.B.c@Gmail.com", want: "abc@gmail.com"},
		{email: "a.B..c@Gmail.com", want: "abc@gmail.com"},
		{email: "a.b.c+001@googlemail.com", want: "abc@gmail.com"},
		{email: "a.b.c+001@whatever.com", want: "a.b.c+001@whatever.com"},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if got := Normalize(tt.email); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestNormalizer_Register(t *testing.T) {
// 	type fields struct {
// 		localFuncs         NormalizeFuncs
// 		domainFuncs        NormalizeFuncs
// 		localFuncsByDomain map[string]NormalizeFuncs
// 	}
// 	type args struct {
// 		domain string
// 		funcs  []NormalizeFunc
// 	}
// 	tests := []struct {
// 		fields fields
// 		argss  []args
// 		email  string
// 		want   string
// 	}{
// 		{
// 			fields: fields{
// 				localFuncs:         nil,
// 				domainFuncs:        nil,
// 				localFuncsByDomain: nil,
// 			},
// 			argss: []args{
// 				{
// 					domain: "",
// 					funcs:  nil,
// 				},
// 				{
// 					domain: "",
// 					funcs:  nil,
// 				},
// 			},
// 			email: "abc@email.com",
// 			want:  "abc@email.com",
// 		},
// 		{
// 			fields: fields{
// 				localFuncs:         defaultFuncs,
// 				domainFuncs:        defaultFuncs,
// 				localFuncsByDomain: nil,
// 			},
// 			argss: []args{
// 				{
// 					domain: "email.COM",
// 					funcs: NormalizeFuncs{
// 						func(s string) string { return s + "+" },
// 						func(s string) string { return s + "m" },
// 					},
// 				},
// 				{
// 					domain: "EMAIL.com",
// 					funcs: NormalizeFuncs{
// 						nil,
// 						func(s string) string { return s + "n" },
// 					},
// 				},
// 			},
// 			email: "ABC@EMAIL.COM",
// 			want:  "abc+mn@email.com",
// 		},
// 	}
// 	for i, tt := range tests {
// 		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
// 			n := NewNormalizer(tt.fields.domainFuncs, tt.fields.localFuncs, tt.fields.localFuncsByDomain)
// 			for _, args := range tt.argss {
// 				n = n.Register(args.domain, args.funcs...)
// 			}
// 			if got := n.Normalize(tt.email); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Normalize() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestNormalizer_Normalize(t *testing.T) {
// 	type fields struct {
// 		localFuncs         NormalizeFuncs
// 		domainFuncs        NormalizeFuncs
// 		localFuncsByDomain map[string]NormalizeFuncs
// 	}
// 	tests := []struct {
// 		fields fields
// 		email  string
// 		want   string
// 	}{
// 		{
// 			fields: fields{
// 				localFuncs:         nil,
// 				domainFuncs:        nil,
// 				localFuncsByDomain: nil,
// 			},
// 			email: "abc@email.com",
// 			want:  "abc@email.com",
// 		},
// 		{
// 			fields: fields{
// 				localFuncs:         nil,
// 				domainFuncs:        nil,
// 				localFuncsByDomain: map[string]NormalizeFuncs{"email.com": nil},
// 			},
// 			email: "abc@email.com",
// 			want:  "abc@email.com",
// 		},
// 		{
// 			fields: fields{
// 				localFuncs:  NormalizeFuncs{strings.ToUpper},
// 				domainFuncs: NormalizeFuncs{strings.ToUpper},
// 				localFuncsByDomain: map[string]NormalizeFuncs{
// 					"EMAIL.COM": NormalizeFuncs{func(s string) string { return s + "+s" }},
// 				},
// 			},
// 			email: "abc@email.com",
// 			want:  "ABC+s@EMAIL.COM",
// 		},
// 	}
// 	for i, tt := range tests {
// 		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
// 			n := NewNormalizer(tt.fields.domainFuncs, tt.fields.localFuncs, tt.fields.localFuncsByDomain)
// 			if got := n.Normalize(tt.email); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Normalizer.Normalize() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestDeleteDots(t *testing.T) {
// 	tests := []struct {
// 		localPart string
// 		want      string
// 	}{
// 		{localPart: "", want: ""},
// 		{localPart: ".", want: ""},
// 		{localPart: "a.b", want: "ab"},
// 		{localPart: "a.b.c", want: "abc"},
// 		{localPart: ".a.b.c.", want: "abc"},
// 		{localPart: "a..b...c", want: "abc"},
// 	}
// 	for i, tt := range tests {
// 		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
// 			if got := DeleteDots(tt.localPart); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("DeleteDots() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestDeleteSubAddr(t *testing.T) {
// 	tests := []struct {
// 		localPart string
// 		want      string
// 	}{
// 		{localPart: "", want: ""},
// 		{localPart: "+", want: ""},
// 		{localPart: "a+b", want: "a"},
// 		{localPart: "a+b+c", want: "a"},
// 		{localPart: "+c", want: ""},
// 	}
// 	for i, tt := range tests {
// 		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
// 			if got := DeleteSubAddr(tt.localPart); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("DeleteSubAddr() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_normalize(t *testing.T) {
// 	tests := []struct {
// 		funcs NormalizeFuncs
// 		str   string
// 		want  string
// 	}{
// 		{
// 			funcs: nil,
// 			str:   "a",
// 			want:  "a",
// 		},
// 		{
// 			funcs: NormalizeFuncs{nil},
// 			str:   "a",
// 			want:  "a",
// 		},
// 		{
// 			funcs: NormalizeFuncs{
// 				func(s string) string { return s + "j" },
// 				func(s string) string { return s + "k" },
// 			},
// 			str:  "a",
// 			want: "ajk",
// 		},
// 	}
// 	for i, tt := range tests {
// 		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
// 			if got := normalize(tt.funcs, tt.str); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("normalize() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestRemoveLocalDots_Normalize(t *testing.T) {
	tests := []struct {
		domains []string
		email   string
		want    string
	}{
		{
			domains: []string{"email.com"},
			email:   "a.b.c@email.com",
			want:    "abc@email.com",
		},
		{
			domains: []string{"email.com"},
			email:   "a..b..c..@email.com",
			want:    "abc@email.com",
		},
		{
			domains: []string{"email.com"},
			email:   "a.b.c@cmail.com",
			want:    "a.b.c@cmail.com",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			ea := NewEmailAddress(tt.email)
			NewRemoveLocalDots(tt.domains...).Normalize(ea)
			if got := ea.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveLocalDots.Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveSubAddressing_Normalize(t *testing.T) {
	tests := []struct {
		tags  map[string]string
		email string
		want  string
	}{
		{
			tags:  map[string]string{"email.com": "+"},
			email: "a@email.com",
			want:  "a@email.com",
		},
		{
			tags:  map[string]string{"email.com": "+"},
			email: "a+b+c@email.com",
			want:  "a@email.com",
		},
		{
			tags:  map[string]string{"email.com": "-"},
			email: "a--b-c@email.com",
			want:  "a@email.com",
		},
		{
			tags:  map[string]string{},
			email: "a+b@email.com",
			want:  "a+b@email.com",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			ea := NewEmailAddress(tt.email)
			NewRemoveSubAddressing(tt.tags).Normalize(ea)
			if got := ea.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveSubAddressing.Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainAlias_Normalize(t *testing.T) {
	tests := []struct {
		aliases map[string]string
		email   string
		want    string
	}{
		{
			aliases: map[string]string{
				"examplemail.com": "email.com",
			},
			email: "a@examplemail.com",
			want:  "a@email.com",
		},
		{
			aliases: map[string]string{},
			email:   "a@examplemail.com",
			want:    "a@examplemail.com",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			ea := NewEmailAddress(tt.email)
			NewDomainAlias(tt.aliases).Normalize(ea)
			if got := ea.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DomainAlias.Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToLowerCase(t *testing.T) {
	tests := []struct {
		email string
		want  string
	}{
		{
			email: "Abc@Email.Com",
			want:  "abc@email.com",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			ea := NewEmailAddress(tt.email)
			ToLowerCase(ea)
			if got := ea.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToLowerCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
