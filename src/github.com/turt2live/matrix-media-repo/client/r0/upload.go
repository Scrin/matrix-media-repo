package r0

import (
	"io"
	"net/http"

	"github.com/turt2live/matrix-media-repo/client"
	"github.com/turt2live/matrix-media-repo/config"
	"github.com/turt2live/matrix-media-repo/media_handler"
	"github.com/turt2live/matrix-media-repo/storage"
	"github.com/turt2live/matrix-media-repo/util"
)

// Request:
//   QS: ?filename=
//   Headers: Content-Type
//   Body: <byte[]>
//
// Response:
//   Body: {"content_uri":"mxc://domain.com/media_id"}

type MediaUploadedResponse struct {
	ContentUri string `json:"content_uri"`
}

func UploadMedia(w http.ResponseWriter, r *http.Request, db storage.Database, c config.MediaRepoConfig) interface{} {
	if !util.IsServerOurs(r.Host, c) {
		return client.AuthFailed()
	}
	accessToken := util.GetAccessTokenFromRequest(r)
	userId, err := util.GetUserIdFromToken(r.Context(), r.Host, accessToken, c)
	if err != nil || userId == "" {
		return client.AuthFailed()
	}

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		filename = "upload.bin"
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream" // binary
	}

	var reader io.Reader
	reader = r.Body
	if c.Uploads.MaxSizeBytes > 0 {
		reader = io.LimitReader(r.Body, c.Uploads.MaxSizeBytes)
	}

	request := &media_handler.MediaUploadRequest{
		UploadedBy:      userId,
		ContentType:     contentType,
		DesiredFilename: filename,
		Host:            r.Host,
		Contents:        reader,
	}

	mxc, err := request.StoreAndGetMxcUri(r.Context(), c, db)
	if err != nil {
		return client.InternalServerError(err.Error())
	}

	return &MediaUploadedResponse{mxc}
}