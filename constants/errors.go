package constants

import "errors"

var ErrInsufficientFunds = errors.New("insufficient funds in account balance")
var ErrThirdPartyFailure = errors.New("third-party failure")
