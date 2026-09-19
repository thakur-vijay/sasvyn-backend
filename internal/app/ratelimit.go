package app

import "github.com/sasvyn/backend/internal/ratelimit"

type RateLimiters struct {
	Default     *ratelimit.Limiter
	AuthRefresh *ratelimit.Limiter
	SocialLogin *ratelimit.Limiter
}

func NewRateLimiters() (*RateLimiters, error) {
	defaultLimiter, err := ratelimit.New(
		ratelimit.Policies.Default.Limit,
		ratelimit.Policies.Default.Window,
	)
	if err != nil {
		return nil, err
	}

	authRefreshLimiter, err := ratelimit.New(
		ratelimit.Policies.AuthRefresh.Limit,
		ratelimit.Policies.AuthRefresh.Window,
	)
	if err != nil {
		defaultLimiter.Close()
		return nil, err
	}

	socialLoginLimiter, err := ratelimit.New(
		ratelimit.Policies.SocialLogin.Limit,
		ratelimit.Policies.SocialLogin.Window,
	)
	if err != nil {
		defaultLimiter.Close()
		authRefreshLimiter.Close()
		return nil, err
	}

	return &RateLimiters{
		Default:     defaultLimiter,
		AuthRefresh: authRefreshLimiter,
		SocialLogin: socialLoginLimiter,
	}, nil
}

func (r *RateLimiters) Close() {
	if r == nil {
		return
	}

	if r.Default != nil {
		r.Default.Close()
	}

	if r.AuthRefresh != nil {
		r.AuthRefresh.Close()
	}

	if r.SocialLogin != nil {
		r.SocialLogin.Close()
	}
}
