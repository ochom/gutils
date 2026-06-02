// Package gutils provides a curated collection of production-ready utility
// packages for building modern Go applications.
//
// The module is organized into focused subpackages so you can import only what
// you need:
//
//   - auth: JWT token generation, validation, and encryption helpers.
//   - cache: Pluggable cache abstraction with memory and Redis drivers.
//   - env: Typed environment variable access with default and required getters.
//   - errors: HTTP-aware error types for API and service applications.
//   - gttp: Unified HTTP client interface with multiple backend implementations.
//   - helpers: Common helpers for strings, CSV, phone numbers, emails, and more.
//   - images: Image compression and conversion helpers.
//   - jsonx: Generic JSON encode/decode helpers.
//   - logs: Lightweight, colorized logging utilities.
//   - nosql: MongoDB helpers with generic CRUD operations.
//   - pubsub: RabbitMQ publisher/consumer utilities.
//   - sqlr: Generic GORM-based SQL data access utilities.
//   - ussd: USSD session and step navigation utilities.
//   - uuid: UUID and short ID helpers.
//   - arrays: Generic array/slice utility helpers.
//
// Installation:
//
//	go get github.com/ochom/gutils
//
// Import the packages you need directly:
//
//	import "github.com/ochom/gutils/env"
//	import "github.com/ochom/gutils/sqlr"
//
// Full API docs: https://pkg.go.dev/github.com/ochom/gutils
package gutils
