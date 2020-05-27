// Author       kevin
// Time         2019-09-27 13:31
// File Desc    JWT 签名相关

package jwt

import (
	"sync"
)

var signingMethods = map[string]func() SigningMethod{}
var signingMethodLock = new(sync.RWMutex)

// Implement SigningMethod to add new methods for signing or verifying tokens.
type SigningMethod interface {
	// Verify signature, return nil if signature is valid
	Verify(signingString, signature string, key interface{}) error
	// Returns encoded signature or error
	Sign(signingString string, key interface{}) (string, error)
	// returns the alg identifier for this method (example: 'HS256')
	Alg() string
}

// Register the "alg" name and a factory function for signing method.
// This is typically done during init() in the method's implementation
func RegisterSigningMethod(alg string, f func() SigningMethod) {
	signingMethodLock.Lock()
	defer signingMethodLock.Unlock()
	signingMethods[alg] = f
}

// Get a signing method from an "alg" string
func GetSigningMethod(alg string) (method SigningMethod) {
	signingMethodLock.RLock()
	defer signingMethodLock.RUnlock()
	if methodF, ok := signingMethods[alg]; ok {
		method = methodF()
	}
	return
}
