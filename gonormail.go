package gonormail

import (
	"fmt"
	"strings"
	"sync"
)

const (
	AT    = "@"
	DOT   = "."
	EMPTY = ""

	domainGmail = "gmail.com"
)

var (
	defaultNormalizers = DefaultNormalizers()
)

type Normalizer interface {
	Normalize(email *EmailAddress)
}

type NormalizeFunc func(*EmailAddress)

func (n NormalizeFunc) Normalize(email *EmailAddress) {
	n(email)
}

type EmailAddress struct {
	Local  string
	Domain string
}

func NewEmailAddress(email string) *EmailAddress {
	splitted := strings.Split(email, AT)
	if len(splitted) != 2 {
		return nil
	}

	return &EmailAddress{Local: splitted[0], Domain: splitted[1]}
}

func (e *EmailAddress) String() string {
	return fmt.Sprintf("%s%s%s", e.Local, AT, e.Domain)
}

// Normalizers struct that holding normalizaters.
type Normalizers struct {
	mux         sync.Mutex
	normalizers []Normalizer
}

// DefaultNormalizers ...
func DefaultNormalizers() *Normalizers {
	return NewNormalizers(
		NormalizeFunc(ToLower),
		NewDomainAlias(map[string]string{"googlemail.com": domainGmail}),
		NewDeleteLocalDots(domainGmail),
		NewCutSubAddressing(map[string]string{domainGmail: "+"}),
	)
}

// NewNormalizers create new Normalizers by given Normalizer
func NewNormalizers(nrs ...Normalizer) *Normalizers {
	normalizers := make([]Normalizer, 0, len(nrs))
	for _, nr := range nrs {
		if nr != nil {
			normalizers = append(normalizers, nr)
		}
	}
	return &Normalizers{normalizers: normalizers}
}

// Normalize normalize given email by registered Normalizer.
func (n *Normalizers) Normalize(email *EmailAddress) {
	for _, nr := range n.normalizers {
		if nr != nil && email != nil {
			nr.Normalize(email)
		}
	}
	return
}

// AppendNormalizer append normalizers.
func (n *Normalizers) AppendNormalizer(nrs ...Normalizer) *Normalizers {
	n.mux.Lock()
	defer n.mux.Unlock()

	for _, nr := range nrs {
		if nr != nil {
			n.normalizers = append(n.normalizers, nr)
		}
	}
	return n
}

// Normalize normalize given email by default Normalizers
func Normalize(email *EmailAddress) {
	defaultNormalizers.Normalize(email)
}

// ToLower normalize local part and domain part to lower case.
func ToLower(email *EmailAddress) {
	email.Local = strings.ToLower(email.Local)
	email.Domain = strings.ToLower(email.Domain)
}

type DeleteLocalDots struct {
	domains map[string]struct{}
}

// NewDeleteLocalDots ...
func NewDeleteLocalDots(domains ...string) *DeleteLocalDots {
	domainMap := make(map[string]struct{}, len(domains))
	for _, domain := range domains {
		domainMap[domain] = struct{}{}
	}
	return &DeleteLocalDots{domains: domainMap}
}

// Normalize ...
func (d *DeleteLocalDots) Normalize(email *EmailAddress) {
	if _, ok := d.domains[email.Domain]; ok {
		email.Local = strings.ReplaceAll(email.Local, DOT, EMPTY)
	}
}

type CutSubAddressing struct {
	tags map[string]string
}

// NewCutSubAddressing ...
func NewCutSubAddressing(tags map[string]string) *CutSubAddressing {
	return &CutSubAddressing{tags: tags}
}

// Normalize ...
func (s *CutSubAddressing) Normalize(email *EmailAddress) {
	if tag, ok := s.tags[email.Domain]; ok {
		email.Local = strings.Split(email.Local, tag)[0]
	}
}

// DomainAlias holding the map of alias -> domain
type DomainAlias struct {
	aliases map[string]string
}

// NewDomainAlias returns a new Normalizer that transfers domain alias to normalized domain.
func NewDomainAlias(aliases map[string]string) *DomainAlias {
	return &DomainAlias{aliases: aliases}
}

// Normalize normalizes domain part of the given email by aliases map.
func (d *DomainAlias) Normalize(email *EmailAddress) {
	if alias, ok := d.aliases[email.Domain]; ok {
		email.Domain = alias
	}
}
