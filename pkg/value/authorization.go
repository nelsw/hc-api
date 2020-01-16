package value

type Authorization struct {
	UserId         string `json:"user_id"`
	SourceIp       string `json:"source_ip"`
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token,omitempty"`
	TemporaryToken string `json:"temporary_token,omitempty"`
}
