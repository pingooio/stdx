package autoupdate

import "github.com/pingooio/stdx/crypto"

const (
	SaltSize = crypto.KeySize256

	ReleaseManifestFilename = "release.json"

	DefaultUserAgent = "Mozilla/5.0 (compatible; +autoupdate)"
)
