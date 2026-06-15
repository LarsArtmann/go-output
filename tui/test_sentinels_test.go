package tui

import "errors"

// Sentinel errors for tests (err113: avoid dynamic errors.New in test code).
var (
	errTestFailure = errors.New("test failure")
	errTestFail    = errors.New("fail")
	errDiskFull    = errors.New("disk full")
)
