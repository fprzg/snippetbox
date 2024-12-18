package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.fepg.org/internal/models"
	"snippetbox.fepg.org/internal/validator"

	//TODO(Farid): Este paquete hace que tengamos que generar demasiado código repetido (declarar un nuevo struct cada que queremos parsear un formulario "html form"). Solucionarlo
	"github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type nameUpdateForm struct {
	Name                string `form:"name"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type emailUpdateForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type passwordUpdateForm struct {
	OldPasswd           string `form:"old_password"`
	NewPasswd           string `form:"new_password"`
	ConfirmNewPasswd    string `form:"confirm_new_password"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "about.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.snippetNotFound(w, r, id)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.snippetNotFound(w, r, id)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetNotFound(w http.ResponseWriter, r *http.Request, id int) {
	data := app.newTemplateData(r)
	data.Snippet = &models.Snippet{ID: id}
	app.render(w, http.StatusNotFound, "snippet_not_found.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	var form snippetCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "Couldn't decode post form for snippet create POST")
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		//fmt.Printf("Shit ain't right\n")
		data := app.newTemplateData(r)
		data.Form = form
		//fmt.Printf("'%s' '%s' '%d'\n", form.Title, form.Content, form.Expires)
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "Couldn't decode post form for user signup POST")
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckEmail(form.Email, "email")
	form.CheckPassword(form.Password, "password")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, err)
		}

		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "Couldn't decode form for user login POST")
		return
	}

	form.CheckEmail(form.Email, "email")
	form.CheckPassword(form.Password, "password")

	redirectPageAfterLogin := app.sessionManager.PopString(r.Context(), "redirectPageAfterLogin")
	app.infoLog.Printf("Redirecting to '%s'\n", redirectPageAfterLogin)
	if redirectPageAfterLogin == "" {
		redirectPageAfterLogin = "/snippet/create"
	}

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect ")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
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

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, redirectPageAfterLogin, http.StatusSeeOther)
}

func (app *application) accountView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	userData, err := app.users.GetUserData(id)
	if err != nil {
		return
	}

	data.UserData = userData
	app.render(w, http.StatusOK, "account.tmpl", data)
}

func (app *application) accountNameUpdateGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = nameUpdateForm{}
	app.render(w, http.StatusOK, "update_name.tmpl", data)
}

func (app *application) accountNameUpdatePost(w http.ResponseWriter, r *http.Request) {
	var form nameUpdateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, "Couldn't decode post form for user name update POST")
		return
	}

	form.CheckPassword(form.Password, "password")
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be empty")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "update_name.tmpl", data)
		return
	}

	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	err = app.users.UpdateName(id, form.Name, form.Password)
	if err != nil {
		data := app.newTemplateData(r)
		data.Form = form

		if errors.Is(err, models.ErrInvalidCredentials) {
			data.Flash = "Invalid credentials"
		} else {
			app.errorLog.Printf("Error occurred. Try again")
			data.Flash = "Error occurred. Try again"
		}

		app.render(w, http.StatusUnprocessableEntity, "update_name.tmpl", data)

		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Name successfully updated!")
	http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}

func (app *application) accountEmailUpdateGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = emailUpdateForm{}
	app.render(w, http.StatusOK, "update_email.tmpl", data)
}

func (app *application) accountEmailUpdatePost(w http.ResponseWriter, r *http.Request) {
	var form emailUpdateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.infoLog.Printf("ERROR DECODING")
		app.clientError(w, http.StatusBadRequest, "Couldn't decode form for user email update POST")
		return
	}

	form.CheckPassword(form.Password, "password")
	form.CheckEmail(form.Email, "email")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "update_email.tmpl", data)
		return
	}

	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	err = app.users.UpdateEmail(id, form.Email, form.Password)
	if err != nil {
		data := app.newTemplateData(r)
		data.Form = form

		if errors.Is(err, models.ErrInvalidCredentials) {
			data.Flash = "Invalid credentials"
		} else {
			app.errorLog.Printf("Error %s", err)
			data.Flash = "Error occurred. Try again"
		}

		app.render(w, http.StatusUnprocessableEntity, "update_email.tmpl", data)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Email successfully updated!")
	http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}

func (app *application) accountPasswordUpdateGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = passwordUpdateForm{}
	app.render(w, http.StatusOK, "update_password.tmpl", data)
}

func (app *application) accountPasswordUpdatePost(w http.ResponseWriter, r *http.Request) {
	const redirectIfFail = "update_password.tmpl"
	var form passwordUpdateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.infoLog.Printf("ERROR DECODING")
		app.clientError(w, http.StatusBadRequest, "Couldn't decode post form for password update POST")
		return
	}

	form.CheckPassword(form.OldPasswd, "old_password")
	form.CheckPassword(form.NewPasswd, "new_password")
	form.CheckPassword(form.ConfirmNewPasswd, "confirm_new_password")
	form.CheckField(validator.NotEqual(form.NewPasswd, form.OldPasswd), "new_password", "New and old password must be different")
	form.CheckField(validator.Equal(form.NewPasswd, form.ConfirmNewPasswd), "confirm_new_password", "Fields must be identical")

	if !form.Valid() {
		//app.infoLog.Printf("The password update form is not valid")
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, redirectIfFail, data)
		return
	}

	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	err = app.users.UpdatePassword(id, form.OldPasswd, form.NewPasswd)
	if err != nil {
		data := app.newTemplateData(r)
		data.Form = form

		if errors.Is(err, models.ErrInvalidCredentials) {
			//app.errorLog.Printf("The fucking credentials again, BITCH!!!")
			data.Flash = "Invalid credentials"
		} else {
			//app.errorLog.Printf("Error ocurred when user %d tried to change password: '%s'", id, err)
			data.Flash = "An error occurred. Try again"
		}

		app.render(w, http.StatusUnprocessableEntity, redirectIfFail, data)

		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Password successfully updated!")
	http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}

func (app *application) userLogoutGet(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
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

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
