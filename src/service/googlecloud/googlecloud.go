package googlecloud

// Config is the structure for Google Cloud Config
type Config struct {
	Type                string `json:"type"`
	ProjectId           string `json:"project_id"`
	ProjectKeyId        string `json:"project_key_id"`
	PrivateKey          string `json:"private_key"`
	ClientEmail         string `json:"client_email"`
	ClientId            string `json:"client_id"`
	AuthUri             string `json:"auth_uri"`
	TokenUri            string `json:"token_uri"`
	AuthProviderCertUrl string `json:"auth_provider_x509_cert_url"`
	ClientCertUrl       string `json:"client_x509_cert_url"`
}
