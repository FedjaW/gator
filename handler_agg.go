package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/fedjaw/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
    if len(cmd.args) != 1 {
        return fmt.Errorf("usage agg <time> (for example agg 1s)")
    }
    timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
    if err != nil {
        return err
    }
    fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)

    ticker := time.NewTicker(timeBetweenReqs)
    for ; ; <- ticker.C {
        scrapeFeeds(s)
    }

    /* feedUrl := "https://www.wagslane.dev/index.xml"
    rssFeed, err := fetchFeed(context.Background(), feedUrl)
    if err != nil {
        return err
    }

    fmt.Printf("%+v\n", rssFeed)

    return nil 
    */
}

func scrapeFeeds(s *state) error {
    feed, err := s.db.GetNextFeedToFetch(context.Background())
    if err != nil {
        return fmt.Errorf("error fetching next feed %w", err)
    }

    err = s.db.MarkFeedFetched(context.Background(), feed.ID)
    if err != nil {
        return fmt.Errorf("error marking feed as marked %w", err)
    }

    rssFeed, err := fetchFeed(context.Background(), feed.Url)
    if err != nil {
        return fmt.Errorf("error fetching feed %w", err)
    }

    for _, post := range rssFeed.Channel.Item {
        fmt.Printf("%s\n", post.Title)

        publishedAt, _ := time.Parse(time.RFC1123, post.PubDate)

        createPostParams := database.CreatePostParams{
            ID: uuid.New(),
            CreatedAt: time.Now().UTC(),
            UpdatedAt:time.Now().UTC(),
            Title: post.Title,
            Url: post.Link,
            Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
            PublishedAt: publishedAt,
            FeedID: feed.ID,
        }
        _, err := s.db.CreatePost(context.Background(), createPostParams)
        if err != nil {
            fmt.Printf("%v\n", err)
            continue
        }
    }

    return nil
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("User-Agent", "gator")

    client := http.Client{}
    res, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    resInBytes, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }

    // var rssFeed RSSFeed // this is also valid
    rssFeed := RSSFeed{} // this is shor for: var rssFeed RSSFeed = RSSFeed{}
    xml.Unmarshal(resInBytes, &rssFeed)

    rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
    rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

    for i, item := range rssFeed.Channel.Item {
        item.Title = html.UnescapeString(item.Title)
        item.Description = html.UnescapeString(item.Description)
        rssFeed.Channel.Item[i] = item
    }

    return &rssFeed, nil
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}


/*  ORIGINAL SOLUTION OF scrapeFeed function
func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}
	for _, item := range feedData.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID:    feed.ID,
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}
*/