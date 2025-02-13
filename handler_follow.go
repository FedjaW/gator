package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fedjaw/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
    if len(cmd.args) != 1 {
        return errors.New("provide a url as argument")
    }

    url := cmd.args[0]

    feed, err := s.db.GetFeedByURL(context.Background(), url)
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

    fmt.Printf("Feed Name: %s\n", feed.Name)
    fmt.Printf("User Name: %s", user.Name)

    return nil
}
