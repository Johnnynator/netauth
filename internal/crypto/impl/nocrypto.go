// +build nocrypto

package impl

import (
	// The blank import here permits the init() within the
	// nocrypto module to register its implementation to the
	// crypto plugin system.
	_ "github.com/NetAuth/NetAuth/internal/crypto/impl/nocrypto"
)
