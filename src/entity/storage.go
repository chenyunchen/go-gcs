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
	To string `json:"to" validate:"required"`
}

// GroupPayload is the structure for group payload
type GroupPayload struct {
	GroupId string `json:"groupId" validate:"required"`
}

// SignedUrlRequest is the structure for signed url response
type SignedUrl struct {
	Url           string `json:"url" validate:"required"`
	UploadHeaders `json:"uploadHeaders" validate:"required"`
	UploadQueries `json:"uploadQueries" validate:"required"`
}

// UploadHeaders is the structure for upload header:
type UploadHeaders struct {
	ContentType        string `json:"Content-Type" validate:"required"`
	ContentLengthRange string `json:"x-goog-content-length-range" validate:"required"`
}

// UploadQueries is the structure for upload query:
type UploadQueries struct {
	Expires        string `json:"Expires" validate:"required"`
	GoogleAccessId string `json:"GoogleAccessId" validate:"required"`
	Signature      string `json:"Signature" validate:"required"`
}

// ResizeImageRequest is the structure for resize image request
type ResizeImageRequest struct {
	Url         string `json:"url" validate:"required"`
	ContentType string `json:"contentType" validate:"required"`
}

// ResizeImage is the structure for resize image response
type ResizeImage struct {
	Origin         string `json:"origin" validate:"required"`
	ThumbWidth100  string `json:"100" validate:"required"`
	ThumbWidth150  string `json:"150" validate:"required"`
	ThumbWidth300  string `json:"300" validate:"required"`
	ThumbWidth640  string `json:"640" validate:"required"`
	ThumbWidth1080 string `json:"1080" validate:"required"`
}
