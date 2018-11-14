package googlecloudstorageprovider

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

// GoogleCloudStorageConfig is the structure for Google Cloud storage Config
type GoogleCloudStorageConfig struct {
	ProjectId         string `json:"projectId"`
	BucketName        string `json:"bucketName"`
	ServiceAccountPEM string `json:"serviceAccountPem"`
	GoogleAccessId    string `json:"GoogleAccessId"`
}

// Service is the structure for Service
type Service struct {
	PrivateKey []byte
	// TODO Create a client if we need.
}

// New will reture a new service
func New(serviceAccountPEM string) *Service {
	privateKey, err := ioutil.ReadFile(serviceAccountPEM)
	if err != nil {
		log.Warn("error while read config file: ", err)
	}

	return &Service{
		PrivateKey: privateKey,
	}
}