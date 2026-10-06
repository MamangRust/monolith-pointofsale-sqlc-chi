package adapter

import (
	"github.com/MamangRust/monolith-point-of-sale-pkg/resilience"
)

// GuardSetter is implemented by every gRPC repository so WithDependencyGuard
// can attach a guard without changing each constructor's client parameter.
type GuardSetter interface {
	SetGuard(*resilience.DependencyGuard)
}

// GuardOption configures a GuardSetter (a gRPC-backed repository). Repositories
// accept zero or more options at construction; the guard wraps every outbound
// gRPC call with a per-call deadline, bulkhead and circuit breaker.
type GuardOption func(GuardSetter)

// WithDependencyGuard returns a GuardOption that attaches a dependency guard
// (per-call timeout + circuit breaker + bulkhead) to a gRPC repository.
// Passing nil disables guarding.
func WithDependencyGuard(guard *resilience.DependencyGuard) GuardOption {
	return func(s GuardSetter) {
		s.SetGuard(guard)
	}
}
