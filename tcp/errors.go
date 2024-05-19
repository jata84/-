/**
 * Copyright gitteamer 2020
 * @date: 2020/11/10
 * @note: error
 */
package tcp

import (
	"errors"
	"fmt"
)

var (
	errRemoteForceDisconnect = errors.New("An existing connection was forcibly closed by the remote host.")
	errRecvEOF               = errors.New("receive EOF from connect.")
	errDataLenOvertakeMaxLen = errors.New("receive data length overtake max data length.")
)

func panicToError(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Panic error: %v", r)
		}
	}()

	fn()
	return
}
