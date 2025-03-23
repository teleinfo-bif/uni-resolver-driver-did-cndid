package document

// DIDResolution CN-DID 文档结构
type DIDResolution struct {
	Context        []string      `json:"context"`        // 一组 URL 数组
	Version        string        `json:"version"`        // CN-DID 文档的版本
	ID             string        `json:"id"`             // 解析的 CN-DID
	PublicKey      []PublicKey   `json:"publicKey"`      // 公钥数组
	Authentication []string      `json:"authentication"` // 一组公钥
	AlsoKnownAs    []AlsoKnownAs `json:"alsoKnownAs"`    // 关联 ID 数组
	Extension      Extension     `json:"extension"`      // 扩展字段
	Service        []Service     `json:"service"`        // 服务地址数组
	Created        string        `json:"created"`        // 创建时间
	Updated        string        `json:"updated"`        // 上次更新时间
	Proof          Proof         `json:"proof"`          // 签名信息
}

// PublicKey 公钥结构
type PublicKey struct {
	ID           string `json:"id"`           // 公钥 ID
	Type         string `json:"type"`         // 公钥算法类型
	Controller   string `json:"controller"`   // 归属的 BID
	PublicKeyHex string `json:"publicKeyHex"` // 十六进制公钥
}

// AlsoKnownAs 关联ID
type AlsoKnownAs struct {
	Type int    `json:"type"` // 关联 ID 的类型
	ID   string `json:"id"`   // 关联 ID
}

// Extension 扩展字段结构
type Extension struct {
	Recovery              []string               `json:"recovery"`              // 一组公钥 ID
	TTL                   uint                   `json:"ttl"`                   // 缓存时间，单位秒
	DelegateSign          DelegateSign           `json:"delegateSign"`          // 第三方对 publicKey 的签名
	Type                  uint                   `json:"type"`                  // 属性类型
	Attributes            []Attribute            `json:"attributes"`            // 一组属性
	VerifiableCredentials []VerifiableCredential `json:"verifiableCredentials"` // 凭证列表
}

// DelegateSign 委托签名结构
type DelegateSign struct {
	Signer         string `json:"signer"`         // 签名公钥 ID
	SignatureValue string `json:"signatureValue"` // 签名的 base64 编码????
}

// Attribute 表示一个属性
type Attribute struct {
	Key     string `json:"key"`     // 属性的 key
	Desc    string `json:"desc"`    // 属性的描述
	Encrypt uint   `json:"encrypt"` // 是否加密，0 非加密，1 加密
	Format  string `json:"format"`  // 数据类型，如 image、text、video、mixture 等
	Value   string `json:"value"`   // 属性自定义 value
}

// VerifiableCredential 表示一个可验证凭证
type VerifiableCredential struct {
	ID   string `json:"id"`   // 凭证 ID
	Type uint   `json:"type"` // 凭证类型
}

// Service 服务结构
type Service struct {
	ID              string `json:"id"`              // 服务地址的 ID
	Type            string `json:"type"`            // 服务类型
	Version         string `json:"version"`         // 解析服务支持的协议版本
	ServerType      uint   `json:"serverType"`      // 解析地址类型
	Protocol        uint   `json:"protocol"`        // 解析服务支持的传输协议
	ServiceEndpoint string `json:"serviceEndpoint"` // 服务的 URL 地址
	Port            uint   `json:"port"`            // 解析端口
}

// Proof 签名结构
type Proof struct {
	Creator        string `json:"creator"`        // 签名公钥 ID
	SignatureValue string `json:"signatureValue"` // 签名的 base58 编码
}
