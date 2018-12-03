package server

import (
	"github.com/chenyunchen/goscraper"
	"go-gcs/src/entity"
	"go-gcs/src/net/context"
	response "go-gcs/src/net/http"
)

func previewUrlHandler(ctx *context.Context) {
	_, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response
	url := req.QueryParameter("url")
	s, err := goscraper.Scrape(url, 5)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	resp.WriteEntity(entity.PreviewedUrl{
		Type:        "default",
		Icon:        s.Preview.Icon,
		Name:        s.Preview.Name,
		Title:       s.Preview.Title,
		Description: s.Preview.Description,
		Image:       s.Preview.Images[0],
		Url:         s.Preview.Link,
	})
}
