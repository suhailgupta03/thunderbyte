package otp

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/suhailgupta03/thunderbyte/otp/models"
	"github.com/suhailgupta03/thunderbyte/otp/providers/smtp"
	"github.com/suhailgupta03/thunderbyte/otp/store"
	"github.com/zerodha/logf"
	"time"
)

const (
	alphaChars     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	numChars       = "0123456789"
	alphaNumChars  = alphaChars + numChars
	uriViewOTP     = "/otp/%s/%s"
	uriViewAddress = "/otp/%s/%s/address"
	uriCheck       = "/otp/%s/%s?otp=%s&action=check"
)

type SetOTPRequest struct {
	// The URL where the server is running
	RootURL            string
	Namespace          string
	Provider           string
	ID                 string
	To                 string
	OtpTTL             time.Duration
	RawMaxAttempts     int
	Extra              []byte
	SMTPConfig         *smtp.Config
	HTMLTemplateName   string
	Subject            string
	Lo                 *logf.Logger
	Store              store.Store
	ChannelDescription string
	AddressDescription string
}

type VerifyOTPRequest struct {
	Namespace string
	Provider  string
	ID        string
	OTPVal    string
	Lo        *logf.Logger
	Store     store.Store
}

type CheckOTPStatus struct {
	Namespace string
	Provider  string
	ID        string
	OTPVal    string
	Lo        *logf.Logger
	Store     store.Store
}

type pushTpl struct {
	To        string
	Namespace string
	Channel   string
	OTP       string
	OTPURL    string
	OTPTTL    time.Duration
}

type OTPResp struct {
	models.OTP
	URL string `json:"url"`
}

// generateRandomString generates a cryptographically random,
// alphanumeric string of length n.
func generateRandomString(totalLen int, chars string) (string, error) {
	bytes := make([]byte, totalLen)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for k, v := range bytes {
		bytes[k] = chars[v%byte(len(chars))]
	}
	return string(bytes), nil
}

// isLocked tells if an OTP is locked after exceeding attempts.
func isLocked(otp models.OTP) bool {
	if otp.Attempts >= otp.MaxAttempts {
		return true
	}
	return false
}

func getURL(rootURL string, otp models.OTP, check bool) string {
	if check {
		return rootURL + fmt.Sprintf(uriCheck, otp.Namespace, otp.ID, otp.OTP)
	}
	return rootURL + fmt.Sprintf(uriViewOTP, otp.Namespace, otp.ID)
}

// push compiles a message template and pushes it to the provider.
func push(otp models.OTP, p *provider, rootURL string, otpTTL time.Duration) error {
	var (
		subj = &bytes.Buffer{}
		out  = &bytes.Buffer{}

		data = pushTpl{
			Channel:   p.provider.ChannelName(),
			Namespace: otp.Namespace,
			To:        otp.To,
			OTP:       otp.OTP,
			OTPURL:    getURL(rootURL, otp, true),
			OTPTTL:    otpTTL,
		}
	)

	if p.tpl != nil {
		if p.tpl.subject != nil {
			if err := p.tpl.subject.Execute(subj, data); err != nil {
				return err
			}
		}

		if p.tpl.body != nil {
			if err := p.tpl.body.Execute(out, data); err != nil {
				return err
			}
		}
	}
	return p.provider.Push(otp, subj.String(), out.Bytes())
}

// verifyOTP validates an OTP against user input.
func verifyOTP(namespace, id, otp string, deleteOnVerify bool, s store.Store, lo *logf.Logger) (models.OTP, error) {
	// Check the OTP.
	out, err := s.Check(namespace, id, true)
	if err != nil {
		if err != store.ErrNotExist {
			lo.Error("error checking OTP", "error", err)
			return out, err
		}
		return out, errors.New("error checking OTP.")
	}

	errMsg := ""
	if isLocked(out) {
		errMsg = fmt.Sprintf("Too many attempts. Please retry after %0.f seconds.",
			out.TTL.Seconds())
	} else if out.OTP != otp {
		errMsg = "Incorrect OTP"
	}

	// There was an error.
	if errMsg != "" {
		return out, errors.New(errMsg)
	}

	// Delete the OTP?
	if deleteOnVerify {
		s.Delete(namespace, id)
	}

	s.Close(namespace, id)
	out.Closed = true
	return out, err
}

// HandleSetOTP creates a new OTP while respecting maximum attempts
// and TTL values.
func HandleSetOTP(req SetOTPRequest) (*OTPResp, error) {
	// TODO: Make the args of initProviders generic to reflect multiple providers
	providers := initProviders(req.SMTPConfig, req.HTMLTemplateName, req.Subject, req.Lo)
	p, ok := providers[req.Provider]
	if !ok {
		req.Lo.Error("Provider not supported. Failed to set OTP", "provider", req.Provider)
		return nil, errors.New(fmt.Sprintf("%s provider not supported. Failed to set OTP", req.Provider))
	}

	// Validate the 'to' address with the provider if one is given.
	if req.To != "" {
		if err := p.provider.ValidateAddress(req.To); err != nil {
			req.Lo.Error("Invalid `to` address", "error", err)
			return nil, errors.New(fmt.Sprintf("Invalid `to` address: %v", err))
		}
	}

	if req.OtpTTL == time.Duration(0) {
		req.Lo.Error("TTL value cannot be empty")
		return nil, errors.New(fmt.Sprintf("TTL value cannot be empty"))
	}
	ttl := time.Second * req.OtpTTL

	if req.RawMaxAttempts == 0 || req.RawMaxAttempts < 1 {
		req.Lo.Error("Max attempts for OTP cannot be empty")
		return nil, errors.New("Max attempts for OTP cannot be empty")
	}

	maxAttempts := req.RawMaxAttempts
	id := req.ID
	if id == "" {
		if oid, err := generateRandomString(32, alphaNumChars); err != nil {
			req.Lo.Error("error generating ID", "error", err)
			return nil, errors.New(fmt.Sprintf("error generating ID %v", err))
		} else {
			id = oid
		}
	}

	otpVal, err := generateRandomString(p.provider.MaxOTPLen(), numChars)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error generating OTP %v", err))
	}

	// Check if the OTP attempts have exceeded the quota.
	otp, err := req.Store.Check(req.Namespace, id, false)
	if err != nil && err != store.ErrNotExist {
		req.Lo.Error("error checking OTP status", "error", err)
		return nil, errors.New(fmt.Sprintf("error checking OTP status %v", err))
	}

	// There's an existing OTP that's locked.
	if err != store.ErrNotExist && isLocked(otp) {
		req.Lo.Error(fmt.Sprintf("OTP attempts exceeded. Retry after %0.f seconds.", otp.TTL.Seconds()))
		return nil, errors.New(fmt.Sprintf("OTP attempts exceeded. Retry after %0.f seconds.", otp.TTL.Seconds()))
	}

	// Create the OTP.
	newOTP, err := req.Store.Set(req.Namespace, id, models.OTP{
		OTP:         otpVal,
		To:          req.To,
		ChannelDesc: req.ChannelDescription,
		AddressDesc: req.AddressDescription,
		Extra:       []byte("{}"),
		Provider:    req.Provider,
		TTL:         ttl,
		MaxAttempts: maxAttempts,
	})

	if err != nil {
		req.Lo.Error("Error setting OTP", "error", err)
		return nil, errors.New(fmt.Sprintf("Error setting OTP %v", err))
	}

	// Push the OTP out.
	if req.To != "" {
		if err := push(newOTP, p, req.RootURL, req.OtpTTL); err != nil {
			req.Lo.Error("error sending OTP", "error", err, "provider", p.provider.ID())
			return nil, errors.New(fmt.Sprintf("Error sending OTP %v provider %s", err, p.provider.ID()))
		}
		req.Lo.Debug("sending otp", "to", newOTP.To, "provider", p.provider.ID(), "namespace", otp.Namespace)

	}

	out := OTPResp{newOTP, getURL(req.RootURL, newOTP, false)}
	return &out, nil
}

// HandleVerifyOTP checks the user input against a stored OTP.
func HandleVerifyOTP(req *VerifyOTPRequest) (*models.OTP, error) {
	if len(req.ID) < 6 {
		req.Lo.Error("ID should be min 6 chars")
		return nil, errors.New("ID should be min 6 chars")
	}
	if req.OTPVal == "" {
		req.Lo.Error("`otp` is empty.")
		return nil, errors.New("`otp` is empty.")
	}

	out, err := verifyOTP(req.Namespace, req.ID, req.OTPVal, true, req.Store, req.Lo)
	return &out, err
}

// HandleCheckOTPStatus checks the user input against a stored OTP.
func HandleCheckOTPStatus(req *CheckOTPStatus) (*models.OTP, error) {
	if len(req.ID) < 6 {
		req.Lo.Error("ID should be min 6 chars.")
		return nil, errors.New("ID should be min 6 chars.")
	}

	// Check the OTP status.
	out, err := req.Store.Check(req.Namespace, req.ID, false)
	if out.Closed {
		req.Store.Delete(req.Namespace, req.ID)
	}
	return &out, err
}
