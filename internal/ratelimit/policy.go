package ratelimit

import "time"

type Policy struct {
	Limit  int
	Window time.Duration
}

var Policies = struct {
	Default     Policy
	AuthRefresh Policy
	SocialLogin Policy
}{
	Default: Policy{
		Limit:  100,
		Window: time.Minute,
	},

	AuthRefresh: Policy{
		Limit:  10,
		Window: time.Minute,
	},

	SocialLogin: Policy{
		Limit:  5,
		Window: time.Minute,
	},
}
