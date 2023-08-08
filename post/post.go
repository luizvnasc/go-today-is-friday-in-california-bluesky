package post

import (
	"context"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/karalabe/go-bluesky"
	"io"
	"log"
)

type Record struct {
	Collection string        `json:"collection" cborgen:"collection"`
	Repo       string        `json:"repo" cborgen:"repo"`
	Record     bsky.FeedPost `json:"record" cborgen:"record"`
}

func Create(ctx context.Context, client *bluesky.Client, post *Record) error {
	return client.CustomCall(func(api *xrpc.Client) error {
		var err error
		var out atproto.ServerDescribeServer_Output
		if err = api.Do(ctx, xrpc.Procedure, "application/json", "com.atproto.repo.createRecord", nil, post, &out); err != nil {
			log.Println(err)
		}
		return err
	})

}

func UploadBlob(ctx context.Context, client *bluesky.Client, blob io.Reader) (out *util.LexBlob, err error) {
	err = client.CustomCall(func(api *xrpc.Client) error {
		if res, err := atproto.RepoUploadBlob(ctx, api, blob); err != nil {
			log.Println(err)
			return err
		} else {
			out = res.Blob
			return err
		}
	})
	return out, err
}
