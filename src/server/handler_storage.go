package server

import (
	"fmt"

	"go-gcs/src/entity"
	"go-gcs/src/googlecloud/storage"
	"go-gcs/src/net/context"
	response "go-gcs/src/net/http"
)

func createGCSSignedUrlHandler(ctx *context.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	userId, ok := req.Attribute("userId").(string)
	if !ok {
		response.Unauthorized(req.Request, resp.ResponseWriter, fmt.Errorf("Unauthorized: User ID is empty"))
		return
	}

	r := entity.SignedUrlRequest{}
	if err := req.ReadEntity(&r); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}
	if err := sp.Validator.Struct(r); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	var signedUrl entity.SignedUrl
	var err error
	switch r.Tag {
	case "single":
		signedUrl, err = storage.CreateGCSSingleSignedUrl(sp, userId, r.FileName, r.ContentType, r.Payload)
	case "group":
		signedUrl, err = storage.CreateGCSGroupSignedUrl(sp, userId, r.FileName, r.ContentType, r.Payload)
	}
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}
	resp.WriteEntity(signedUrl)
}

func resizeGCSImageHandler(ctx *context.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	r := entity.ResizeImageRequest{}
	if err := req.ReadEntity(&r); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}
	if err := sp.Validator.Struct(r); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	resizeImage, err := storage.ResizeGCSImage(sp, r.Url, r.ContentType)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}
	resp.WriteEntity(resizeImage)
}
