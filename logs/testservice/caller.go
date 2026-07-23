// Package testservice exists only so logs' tests have a real, separate
// caller package to prove callerService correctly walks past logs' own
// frames to find the actual external caller.
package testservice

import "github.com/scrambledeggs/booky-go-common/logs"

// CallInfo calls logs.Info from outside the logs package.
func CallInfo(note string) {
	logs.Info(note)
}
