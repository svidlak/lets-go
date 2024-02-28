package main

import (
	"errors"
	"net/http"

	"github.com/svidlak/lets-go/internal/models"
	"github.com/svidlak/lets-go/internal/validator"
)

type userSignupForm struct {
	Name     string
	Email    string
	Password string
	validator.Validator
}

type userLoginForm struct {
	Email    string
	Password string
	validator.Validator
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	templateData := app.initTemplateData(r)
	templateData.Form = userSignupForm{}

	app.render(w, http.StatusOK, "signup", templateData)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	form := userSignupForm{
		Name:     name,
		Email:    email,
		Password: password,
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Name, 100), "name", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Email, 100), "email", "This field cannot be more than 100 characters long")
	form.CheckField(validator.ValidEmail(form.Email), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Password, 100), "password", "This field cannot be more than 100 characters long")

	if !form.Valid() {
		data := app.initTemplateData(r)
		data.Form = form

		app.render(w, http.StatusUnprocessableEntity, "signup", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)

	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			templateData := app.initTemplateData(r)
			templateData.Form = form

			app.render(w, http.StatusUnprocessableEntity, "signup", templateData)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "You have been registered")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	templateData := app.initTemplateData(r)
	templateData.Form = userLoginForm{}

	app.render(w, http.StatusOK, "login", templateData)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	form := userLoginForm{Email: email, Password: password}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Email, 100), "email", "This field cannot be more than 100 characters long")
	form.CheckField(validator.ValidEmail(form.Email), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Password, 100), "password", "This field cannot be more than 100 characters long")

	if !form.Valid() {
		data := app.initTemplateData(r)
		data.Form = form

		app.render(w, http.StatusUnprocessableEntity, "login", data)
		return
	}

	user, err := app.users.Authenticate(form.Email, form.Password)

	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.initTemplateData(r)
			data.Form = form

			app.render(w, http.StatusUnprocessableEntity, "login", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", user.ID)
	app.sessionManager.Put(r.Context(), "UserName", user.Name)

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
