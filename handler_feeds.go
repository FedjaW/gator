package main

import (
	"context"
	"fmt"

	"github.com/fedjaw/gator/internal/database"
)

func handlerFeeds(s *state, cmd command) error {
    feeds, err := s.db.GetFeeds(context.Background())
    if err != nil {
        return fmt.Errorf("error retrieving feeds: %w", err)
    }

    for _, feed := range feeds {
        user, err := s.db.GetUserByID(context.Background(), feed.UserID)
        if err != nil {
            return fmt.Errorf("error retrieving user that created the feed: %w", err)
        }
        printFeed(feed, user)
        fmt.Println("======================================")
    }


    return nil
}

func printFeed(feed database.Feed, user database.User) {
    fmt.Printf("* Name:         %s\n", feed.Name)
    fmt.Printf("* URL:          %s\n", feed.Url)
    fmt.Printf("* User:         %s\n", user.Name)
}
