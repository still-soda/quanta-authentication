package config

type ResponseType string

const (
	ResponseTypeCode        ResponseType = "code"
	ResponseTypeIDToken     ResponseType = "id_token"
	ResponseTypeToken       ResponseType = "token"
	ResponseTypeCodeIDToken ResponseType = "code id_token"
)

type Scope string

const (
	ScopeOpenID  Scope = "openid"
	ScopeProfile Scope = "profile"
	ScopeEmail   Scope = "email"
	ScopeRoles   Scope = "roles"
)

type ResponseMode string

const (
	ResponseModeQuery    ResponseMode = "query"
	ResponseModeFragment ResponseMode = "fragment"
	ResponseModeFormPost ResponseMode = "form_post"
)

type GrantType string

const (
	GrantTypeAuthorizationCode GrantType = "authorization_code"
	GrantTypeImplicit          GrantType = "__implicit"
	GrantTypeClientCredentials GrantType = "client_credentials"
	GrantTypePassword          GrantType = "password"
	GrantTypeRefreshToken      GrantType = "refresh_token"
)

type SubjectType string

const (
	SubjectTypePublic SubjectType = "public"
)

type IDTokenSigningAlg string

const (
	IDTokenSigningAlgRS256 IDTokenSigningAlg = "RS256"
)

type TokenEndpointAuthMethod string

const (
	TokenEndpointAuthMethodClientSecretBasic TokenEndpointAuthMethod = "client_secret_basic"
	TokenEndpointAuthMethodClientSecretPost  TokenEndpointAuthMethod = "client_secret_post"
)

type CodeChallengeMethod string

const (
	CodeChallengeMethodS256  CodeChallengeMethod = "S256"
	CodeChallengeMethodPlain CodeChallengeMethod = "plain"
)

type Claim string

const (
	// basic
	ClaimSub Claim = "sub"
	ClaimIss Claim = "iss"
	ClaimAud Claim = "aud"
	ClaimExp Claim = "exp"
	ClaimIat Claim = "iat"

	// profile
	ClaimName        Claim = "name"
	ClaimDisplayName Claim = "display_name"
	ClaimStudentID   Claim = "student_id"
	ClaimPicture     Claim = "picture"

	// email
	ClaimEmail         Claim = "email"
	ClaimEmailVerified Claim = "email_verified"

	// roles
	ClaimRoles Claim = "roles"
)
