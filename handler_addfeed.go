package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fedjaw/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
    if len(cmd.args) != 2 {
        return errors.New("name or url was not provided, usage: go run . addfeed <name> <url>")
    }

    name := cmd.args[0] // name of the feed
    url := cmd.args[1] // url of the feed

    ctx := context.Background()

    fetchFeed(ctx, url)

    createFeedParams := database.CreateFeedParams{
        ID: uuid.New(),
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
        Name: name,
        Url: url,
        UserID: user.ID,
    }

    feed, err := s.db.CreateFeed(ctx, createFeedParams)
    if err != nil {
        return err
    }

    createFeedFollowParams := database.CreateFeedFollowParams{
        ID: uuid.New(),
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
        UserID: user.ID,
        FeedID: feed.ID,
    }

    _, err = s.db.CreateFeedFollow(context.Background(), createFeedFollowParams)
    if err != nil {
        return err
    }

    fmt.Printf("feed created: %+v", feed)

    return nil
}
