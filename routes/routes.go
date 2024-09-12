package routes

import (
	"github.com/Falokut/images_storage_service/controller"
	"net/http"

	"github.com/Falokut/go-kit/http/endpoint"
	"github.com/Falokut/go-kit/http/router"
)

type Router struct {
	Images controller.Images
}

func (r Router) InitRoutes(wrapper endpoint.Wrapper) *router.Router {
	mux := router.New()
	for _, desc := range endpointDescriptors(r) {
		mux.Handler(desc.Method, desc.Path, wrapper.Endpoint(desc.Handler))
	}

	return mux
}

type EndpointDescriptor struct {
	Method  string
	Path    string
	Handler any
}

func endpointDescriptors(r Router) []EndpointDescriptor {
	return []EndpointDescriptor{
		{
			Method:  http.MethodPost,
			Path:    "/image/:category",
			Handler: r.Images.UploadImage,
		},
		{
			Method:  http.MethodGet,
			Path:    "/image/:category/:imageId",
			Handler: r.Images.GetImage,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/image/:category/:imageId",
			Handler: r.Images.DeleteImage,
		},
		{
			Method:  http.MethodGet,
			Path:    "/image/:category/:imageId/exist",
			Handler: r.Images.IsImageExist,
		},
		{
			Method:  http.MethodPost,
			Path:    "/image/:category/:imageId/replace",
			Handler: r.Images.ReplaceImage,
		},
	}
}
