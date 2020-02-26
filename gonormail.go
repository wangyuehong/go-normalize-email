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
)

var (
	defaultEmailNormalizer = DefaultEmailNormalizer()
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

// EmailNormalizer struct that holding normalizaters.
type EmailNormalizer struct {
	mux         sync.Mutex
	normalizers []Normalizer
}

// DefaultNormalizer ...
func DefaultEmailNormalizer() *EmailNormalizer {
	return NewEmailNormalizer(
		NormalizeFunc(ToLowerCase),
		NewDomainAlias(map[string]string{"googlemail.com": "gmail.com"}),
		NewRemoveLocalDots("gmail.com"),
		NewRemoveSubAddressing(map[string]string{"gmail.com": "+"}),
	)
}

// NewEmailNormalizer create new EmailNormalizer by given Normalizer
func NewEmailNormalizer(nrs ...Normalizer) *EmailNormalizer {
	normalizers := make([]Normalizer, 0, len(nrs))
	for _, nr := range nrs {
		if nr != nil {
			normalizers = append(normalizers, nr)
		}
	}
	return &EmailNormalizer{normalizers: normalizers}
}

// Normalize normalize given email by registered Normalizer.
func (n *EmailNormalizer) Normalize(email string) string {
	emailAddress := NewEmailAddress(email)
	if emailAddress == nil {
		return email
	}

	for _, nr := range n.normalizers {
		nr.Normalize(emailAddress)
	}
	return emailAddress.String()
}

// AddNormalizer append normalizers.
func (n *EmailNormalizer) AddNormalizer(nrs ...Normalizer) *EmailNormalizer {
	n.mux.Lock()
	defer n.mux.Unlock()

	for _, nr := range nrs {
		if nr != nil {
			n.normalizers = append(n.normalizers, nr)
		}
	}
	return n
}

// Normalize normalize given email by default EmailNormalizer
func Normalize(email string) string {
	return defaultEmailNormalizer.Normalize(email)
}

// ToLowerCase normalize local part and domain part to lower case.
func ToLowerCase(email *EmailAddress) {
	email.Local = strings.ToLower(email.Local)
	email.Domain = strings.ToLower(email.Domain)
}

type RemoveLocalDots struct {
	domains map[string]struct{}
}

// NewRemoveLocalDots ...
func NewRemoveLocalDots(domains ...string) *RemoveLocalDots {
	domainMap := make(map[string]struct{}, len(domains))
	for _, domain := range domains {
		domainMap[domain] = struct{}{}
	}
	return &RemoveLocalDots{domains: domainMap}
}

// Normalize ...
func (d *RemoveLocalDots) Normalize(email *EmailAddress) {
	if _, ok := d.domains[email.Domain]; ok {
		email.Local = strings.ReplaceAll(email.Local, DOT, EMPTY)
	}
}

type RemoveSubAddressing struct {
	tags map[string]string
}

// NewRemoveSubAddressing ...
func NewRemoveSubAddressing(tags map[string]string) *RemoveSubAddressing {
	return &RemoveSubAddressing{tags: tags}
}

// Normalize ...
func (s *RemoveSubAddressing) Normalize(email *EmailAddress) {
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
