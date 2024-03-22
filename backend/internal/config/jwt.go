package config

// Jwt
// @Description: jwt 认证相关配置
type Jwt struct {
	SigningKey     string `mapstructure:"signing_key" yaml:"signing_key" json:"signing_key,omitempty"`             // 用于签名和验证 JWT 的秘钥
	ExpirationTime string `mapstructure:"expiration_time" yaml:"expiration_time" json:"expiration_time,omitempty"` // 令牌过期时间
	BufferTime     string `mapstructure:"buffer_time" yaml:"buffer_time" json:"buffer_time,omitempty"`             // 令牌缓冲时间
	Issuer         string `mapstructure:"issuer" yaml:"issuer" json:"issuer,omitempty"`                            // 令牌发行人标识
}
