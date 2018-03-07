package password

import (
	"errors"
	"html/template"
	"path"
	"ems/auth"
	"ems/auth/claims"
	"ems/mailer"
	"net/mail"
	"ems/core/utils"
)

var (
	// ConfirmationMailSubject confirmation mail's subject
	ConfirmationMailSubject = "Please confirm your account"

	// ConfirmedAccountFlashMessage confirmed your account message
	ConfirmedAccountFlashMessage = template.HTML("Confirmed your account!")

	// ConfirmFlashMessage confirm account flash message
	ConfirmFlashMessage = template.HTML("Please confirm your account")

	// ErrAlreadyConfirmed account already confirmed error
	ErrAlreadyConfirmed = errors.New("Your account already been confirmed")

	// ErrUnconfirmed unauthorized error
	ErrUnconfirmed = errors.New("You have to confirm your account before continuing")
)


// DefaultConfirmationMailer default confirm mailer
var DefaultConfirmationMailer = func(email string, context *auth.Context, claims *claims.Claims, currentUser interface{}) error {
	claims.Subject = "confirm"

	return context.Auth.Mailer.Send(
		mailer.Email{
			TO:      []mail.Address{{Address: email}},
			Subject: ConfirmationMailSubject,
		}, mailer.Template{
			Name:    "auth/confirmation",
			Data:    context,
			Request: context.Request,
			Writer:  context.Writer,
		}.Funcs(template.FuncMap{
			"current_user": func() interface{} {
				return currentUser
			},
			"confirm_url": func() string {
				confirmURL := utils.GetAbsURL(context.Request)
				confirmURL.Path = path.Join(context.Auth.AuthURL("password/confirm"))
				qry := confirmURL.Query()
				qry.Set("token", context.SessionStorer.SignedToken(claims))
				confirmURL.RawQuery = qry.Encode()
				return confirmURL.String()
			},
		}))
}