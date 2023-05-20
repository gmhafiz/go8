package authentication

type RespondCsrf struct {
	CsrfToken string `json:"csrf_token"`
}
