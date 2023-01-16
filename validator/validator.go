package validator

type SignUpForm struct {
	CompanyName          string `validate:"required"`
	FirstName            string `validate:"required"`
	LastName             string `validate:"required"`
	Email                string `validate:"required|email"`
	Password             string `validate:"required|min_len:8"`
	PasswordConfirmation string `validate:"required|min_len:8"`
}

type SignInForm struct {
	Email    string `validate:"required|email"`
	Password string `validate:"required|min_len:8"`
}

type SignInCodeForm struct {
	Code string `validate:"required"`
}

type EmailForm struct {
	Email string `validate:"required|email"`
}

type EmailResetForm struct {
	Email             string `validate:"required|email"`
	EmailConfirmation string `validate:"required|email"`
	Token             string `validate:"required"`
}

type PasswordResetForm struct {
	Password             string `validate:"required|min_len:8"`
	PasswordConfirmation string `validate:"required|min_len:8"`
	Token                string `validate:"required"`
}

type ProfileForm struct {
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Phone     string `validate:"required"`
}

type AccountForm struct {
	Name string `validate:"required"`
}

type RoleForm struct {
	Name string `validate:"required"`
}

type PermissionForm struct {
	Name string `validate:"required"`
}
