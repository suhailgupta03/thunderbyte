package otp

type OTPErrorCode int

const (
	OTPErrorUnknown OTPErrorCode = iota
	MaxAttemptsExceeded
	InvalidOTP
	OTPExpired
	ProviderNotSupported
	ConnectionToStoreFailed
	SettingOTPFailed
	SendingOTPFailed
)

type OTPError struct {
	Message    string
	ErrorCode  OTPErrorCode
	RetryAfter float64
}

func (e *OTPError) Error() string {
	return e.Message
}

func NewOTPError(code OTPErrorCode, message string) *OTPError {
	return &OTPError{
		Message:   message,
		ErrorCode: code,
	}
}
