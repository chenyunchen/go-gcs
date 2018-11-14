package entity

// SignedUrlRequest is the structure for signed url request
type SignedUrlRequest struct {
	FileName    string `json:"fileName" validate:"required"`
	ContentType string `json:"contentType" validate:"required"`
	Tag         string `json:"tag" validate:"required"`
	Payload     string `json:"payload" validate:"required"`
}

// SinglePayload is the structure for single payload
type SinglePayload struct {
	From string `json:"from" validate:"required"`
	To   string `json:"to" validate:"required"`
}

// GroupPayload is the structure for group payload
type GroupPayload struct {
	From    string `json:"from" validate:"required"`
	GroupId string `json:"groupId" validate:"required"`
}

// SignedUrlRequest is the structure for signed url response
type SignedUrl struct {
	Url          string `json:"url" validate:"required"`
	UploadQuerys `json:"uploadQuerys" validate:"required"`
}

// GroupPayload is the structure for upload query:
type UploadQuerys struct {
	Expires        string `json:"Expires" validate:"required"`
	GoogleAccessId string `json:"GoogleAccessId" validate:"required"`
	Signature      string `json:"Signature" validate:"required"`
}
