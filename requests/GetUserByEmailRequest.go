package requests

type GetUserByEmailRequest struct {
	Email string `query:"email"`
}
