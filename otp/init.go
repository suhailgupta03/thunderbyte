package otp

import (
	"github.com/Masterminds/sprig"
	"github.com/suhailgupta03/thunderbyte/otp/models"
	"github.com/suhailgupta03/thunderbyte/otp/providers/smtp"
	"github.com/zerodha/logf"
	"html/template"
	"path/filepath"
	"strings"
)

type providerTpl struct {
	subject *template.Template
	body    *template.Template
}

type provider struct {
	provider models.Provider
	tpl      *providerTpl
}

// initProviderTpl loads a provider's optional templates.
func initProviderTpl(subj, tplFile string, lo *logf.Logger) *providerTpl {
	out := &providerTpl{}

	// Template file.
	if tplFile != "" {
		// Parse the template file.
		// tpl, err := template.ParseFiles(tplFile)

		tpl, err := template.New(filepath.Base(tplFile)).Funcs(sprig.FuncMap()).ParseFiles(tplFile)

		if err != nil {
			lo.Fatal("error parsing template file", "tplFile", tplFile, "error", err)
		}
		out.body = tpl
	}

	// Subject template string.
	if subj != "" {
		tpl, err := template.New("subject").Parse(subj)
		if err != nil {
			lo.Fatal("error parsing template subject", "tplFile", tplFile, "error", err)
		}

		out.subject = tpl
	}

	return out
}

// initProviders loads models.Provider plugins from the list of given filenames.
func initProviders(cfg *smtp.Config, templateName string, lo *logf.Logger) map[string]*provider {
	out := make(map[string]*provider)
	// Initialized the in-built providers.
	// SMTP.
	smtpEnabled := true
	if smtpEnabled && cfg != nil {
		p, err := smtp.New(*cfg)
		if err != nil {
			lo.Fatal("error initializing smtp provider", "error", err)
		}

		out["smtp"] = &provider{
			provider: p,
			tpl:      initProviderTpl("Test Subject", templateName, lo),
		}
	}

	if len(out) == 0 {
		lo.Fatal("no providers or webhooks enabled")
	}

	names := []string{}
	for name := range out {
		names = append(names, name)
	}

	lo.Info("enabled providers:", strings.Join(names, ", "))

	return out
}
