package cachekey

// Predefined namespaces. New entries SHOULD be added here rather than calling
// NewNamespace at the call site, so all Redis keys remain discoverable.
var (
	// OAuthEmailVerification keys are "oauth:email_verification:<token>".
	OAuthEmailVerification = NewNamespace("oauth", "email_verification")
)
