package titan

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/curve25519"
)

// Identity encapsulates the Titan platform component identity
type Identity struct {
	msid         int
	titanType    int
	titanTypeVer int
	key          [32]byte
	signature    [64]byte
	privateKey   [32]byte
}

func (ti *Identity) GetMsid() int {
	return ti.msid
}

func (ti *Identity) GetType() int {
	return ti.titanType
}

func (ti *Identity) GetTypeVer() int {
	return ti.titanTypeVer
}

func (ti *Identity) GetKey() [32]byte {
	return ti.key
}

func (ti *Identity) GetKeyAsString() string {
	return base64.StdEncoding.EncodeToString(ti.key[:])
}

func (ti *Identity) GetSignature() [64]byte {
	return ti.signature
}

func (ti *Identity) GetSignatureAsString() string {
	return base64.StdEncoding.EncodeToString(ti.signature[:])
}

func (ti *Identity) GetCertificate() [64]byte {
	return ti.signature
}

func (ti *Identity) GetCertificateAsString() string {
	return base64.StdEncoding.EncodeToString(ti.signature[:])
}

// PrettyPrint is a helper method to create a logfile string containing the state
// of the object
func (ti *Identity) PrettyPrint() string {
	var sb strings.Builder

	sb.WriteString("Titan Identity:\n")
	fmt.Fprintf(&sb, "\tType        : 0x%08x\n", ti.titanType)
	fmt.Fprintf(&sb, "\tType Version: 0x%08x\n", ti.titanTypeVer)
	fmt.Fprintf(&sb, "\tMSID        : 0x%08x\n", ti.msid)
	fmt.Fprintf(&sb, "\tPublic Key  : %s\n", ti.GetKeyAsString())
	fmt.Fprintf(&sb, "\tCertificate : %s\n", ti.GetCertificateAsString())

	return sb.String()
}

// NewIdentity creates a new Titan Identity
func NewIdentity(msid int, titanType int, titanTypeVer int) (*Identity, error) {
	id := Identity{msid: msid, titanType: titanType, titanTypeVer: titanType}

	// Generate private key
	_, err := rand.Read(id.privateKey[:])
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}
	// Calculate public key
	curve25519.ScalarBaseMult(&id.key, &id.privateKey)

	return &id, nil
}
