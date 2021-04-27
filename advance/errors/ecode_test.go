package errors

import (
	"fmt"
	"testing"

	ecode "gotest/advance/errors"

	"github.com/pkg/errors"
)

func errNative() error {
	errCode := ecode.ServerErr
	err := errors.Wrap(errCode, "native error")
	return err
}

func errWithMsg() error {
	err := errNative()
	return errors.WithMessage(err, "with message")
}

func errProcess() error {
	err := errNative()
	if Cause(err) == ecode.ServerErr {
		//fmt.Printf("[returnProcessErr]: %+v \n", err)
		err = ecode.ServerErr2
	}

	return err
}

func TestErrors(t *testing.T) {
	err := errProcess()
	fmt.Printf("[ErrProcess]: %+v \n", GRpcStatusFromEcode(err))

	err = errWithMsg()
	fmt.Printf("[ErrWithMsg]: uid: 1001, %+v \n", err)
	fmt.Println("<--------------------")
}
