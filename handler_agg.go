package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

func handlerAgg(s *state, cmd command) error {
    if len(cmd.args) != 1 {
        return fmt.Errorf("usage agg <time>")
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

    for _, item := range rssFeed.Channel.Item {
        fmt.Printf("%s\n", item.Title)
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
