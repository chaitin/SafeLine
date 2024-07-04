package constants

const (
	SuperUser      = "admin"
	ProductName    = "长亭雷池 WAF 社区版"
	ProductVersion = ""
	ConfigFilePath = "config.yml"
	CertsPath      = "certs"
)

const (
	NotUpgrade int = iota
	RecommendedUpgrade
	MustUpgrade
)
