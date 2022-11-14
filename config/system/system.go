package system

// System - системный конфиг.
type System struct {
	// PasswordChars символы используемые для генерации паролей
	// (по умолчанию 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789').
	PasswordChars string `json:"password_chars" toml:"password_chars" yaml:"password_chars" form:"password_chars" description:"Символы используемые для генерации паролей (по умолчанию 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789')."`

	// PasswordLength кол-во символов используемых для генерации паролей
	// (по умолчанию 12).
	PasswordLength int `json:"password_length" toml:"password_length" yaml:"password_length" form:"password_length" description:"Кол-во символов используемых для генерации паролей (по умолчанию 12)."`
}
