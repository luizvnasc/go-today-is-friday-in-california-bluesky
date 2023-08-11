package main

import (
	"context"
	"errors"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"github.com/luizvnasc/bluesky.bot/post"
	"github.com/mercadolibre/golang-restclient/rest"
	"log"
	"os"
	"time"
)
import bluesky "github.com/karalabe/go-bluesky"

var (
	blueskyHandle = os.Getenv("blueskyHandle")
	blueskyAppkey = os.Getenv("blueskyAppkey")
	repo          = os.Getenv("repo")
	BaseURL       = os.Getenv("baseUrl")
	port          = os.Getenv("PORT")
)

func main() {

	println("Configs")

	println(blueskyAppkey)
	println(blueskyHandle)
	println(repo)
	println(BaseURL)

	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(5).Seconds().Do(func() {
		rest.Get(BaseURL)
	})

	if err != nil {
		log.Fatal(err)
	}

	_, err = s.Every(1).Friday().At("09:30").Do(func() {
		err := todayIsFridayInCalifornia()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("TODAY IS FRIDAY IN CALIFORNIA! SHOOT!")
	})

	if err != nil {
		log.Fatal(err)
	}
	s.StartAt(time.Now().Add(3 * time.Second)).StartAsync()
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))

}

func login(ctx context.Context, handle, appkey string) (*bluesky.Client, error) {
	client, err := bluesky.Dial(ctx, bluesky.ServerBskySocial)
	if err != nil {
		panic(err)
	}

	err = client.Login(ctx, handle, appkey)
	switch {
	case errors.Is(err, bluesky.ErrMasterCredentials):
		log.Println("You're not allowed to use your full-access credentials, please create an appkey")
	case errors.Is(err, bluesky.ErrLoginUnauthorized):
		log.Println("Username of application password seems incorrect, please double check")
	case err != nil:
		log.Println("Something else went wrong, please look at the returned error")
	}
	return client, err
}

func todayIsFridayInCalifornia() (err error) {
	ctx := context.Background()
	client, err := login(ctx, blueskyHandle, blueskyAppkey)
	defer client.Close()

	file, err := os.Open("resources/TodayIsFridayInCalifornia.jpg")
	if err != nil {
		log.Fatal(err)
		return
	}

	image, err := post.UploadBlob(ctx, client, file)
	if err != nil {
		log.Fatal(err)
		return
	}

	p := bsky.FeedPost{
		LexiconTypeID: "app.bsky.feed.post",
		CreatedAt:     time.Now().Format(time.RFC3339),
		Embed: &bsky.FeedPost_Embed{
			EmbedImages: &bsky.EmbedImages{
				Images: []*bsky.EmbedImages_Image{
					{
						Image: image,
					},
				},
			},
		},
	}

	record := &post.Record{
		Collection: "app.bsky.feed.post",
		Repo:       repo,
		Record:     p,
	}

	if err = post.Create(ctx, client, record); err != nil {
		log.Println(err)
		return
	}
	return
}
