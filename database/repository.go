package database

import "github.com/jmoiron/sqlx"

const (
	SETTINGS_REPO = "thunderbyte_settings"
	AUTH_USER     = "auth_users"
	AUTH_PASSWORD = "auth_passwords"
)

type ThunderByteSetting struct {
	ID    int    `db:"id"`
	Key   string `db:"key"`
	Value string `db:"value"`
}

type ThunderByteSettings []ThunderByteSetting

type VerifiedUser struct {
	UserId   int64  `db:"userid"`
	Username string `db:"username"`
}

type AuthProfile struct {
	UserId   int64  `db:"id"`
	Username string `db:"username"`
}

// Queries contains all prepared SQL queries.
type Queries interface{}

type ThunderbyteQueries struct {
	GetAllSettings             *sqlx.Stmt `query:"get-all-settings"`
	GetSettingByKey            *sqlx.Stmt `query:"get-setting-by-key"`
	VerifyCredentials          *sqlx.Stmt `query:"verify-creds"`
	FetchAuthProfileByUsername *sqlx.Stmt `query:"fetch-auth-profile-by-username"`
	FetchAuthProfileById       *sqlx.Stmt `query:"fetch-auth-profile-by-id"`
	CreateAuthProfile          *sqlx.Stmt `query:"create-auth-profile"`
	CreatePassword             *sqlx.Stmt `query:"create-password"`
}
