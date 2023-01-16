package request

type SignUpInput struct {
	CompanyName          string `json:"company_name"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInCodeInput struct {
	Code string `json:"code"`
}

type EmailInput struct {
	Email string `json:"email"`
}

type EmailResetInput struct {
	Email             string `json:"email"`
	EmailConfirmation string `json:"email_confirmation"`
	Token             string `json:"token"`
}

type PasswordResetInput struct {
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
	Token                string `json:"token"`
}

type TokenInput struct {
	Token string `json:"token"`
}

type ProfileInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

type AccountInput struct {
	Name string `json:"name"`
}

type RoleInput struct {
	Name string `json:"name"`
}

type PermissionInput struct {
	Name string `json:"name"`
}
