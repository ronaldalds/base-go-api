package i18n

// Mensagens de sucesso General
const (
	MSG_START_SERVER = "Server started on port %d"
)

// Mensagens de erro General
const (
	ERR_INVALID_TIMEZONE = "invalid timezone: %s"
)

// Mensagens de erro JWT
const (
	ERR_INVALID_JWT = "invalid jwt token"
	ERR_SIGN_JWT    = "could not sign access token string %v"
)

// Mensagens de erro Variaveis de Ambiente
const (
	ERR_VARIABLE_IS_REQUIRED = "variable %s is required"
	ERR_CONVERT_TO_INTEGER   = "failed to convert %s to integer: %v"
)

// Mensagens de erro Password
const (
	ERR_CRYPTING_PASSWORD_FAILED = "crypting password failed: %s"
)

// Mensagens de erro Validação
const (
	ERR_PASSWORD_INVALID_UPPERCASE = "password must contain at least one uppercase letter."
	ERR_PASSWORD_INVALID_SYMBOL    = "password must contain at least one symbol."
	ERR_VALIDATE_STRUCTURE_DEFAULT = "validate Structure Failed!"
	ERR_INVALID_ID_PARAMS          = "ID must be a valid number"
	ERR_INVALID_BODY               = "invalid body: %s"
	ERR_INVALID_QUERY_STRING       = "invalid query string: %s"
	ERR_INVALID_FIELD              = "invalid field: %s"
)

// Mensagens de erro Database
const (
	ERR_SET_KEY_REDIS        = "failed to set key redis: %s"
	ERR_CREATE_USER          = "failed to create user: %s"
	ERR_SELECT_CONSENT_ID    = "failed to select consent id: %s"
	ERR_SELECT_USER_OR_EMAIL = "user with username '%s' or email '%s' already exists"
	ERR_LOGIN                = "failed to login: username or password is incorrect"
)

// Mensagens de erro Client S3
const (
	ERR_CREATE_S3_CLIENT = "failed to create S3 client: %s"
)

// Mensagens de erro Presigner
const (
	ERR_PRESIGNER_FAILED            = "presigner Failed!"
	ERR_GENERATE_PRESIGN_URL_FRONT  = "failed to generate presign url front: %s"
	ERR_GENERATE_PRESIGN_URL_BACK   = "failed to generate presign url back: %s"
	ERR_GENERATE_PRESIGN_URL_SELFIE = "failed to generate presign url selfie: %s"
)

// Mensagens de erro Paginação
const (
	ERR_PAGE_NOT_EXIST = "page does not exist"
)
