package aws

// Host representa o host da aplicacao
type Host struct {
	URL      string
	Port     int64
	Protocol string
	Arn      string
}

// GatewayResource representa o host da aplicacao
type GatewayResource struct {
	ID   string
	Path string
}
