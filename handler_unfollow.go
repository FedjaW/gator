package main

import (
	"context"
	"fmt"

	"github.com/fedjaw/gator/internal/database"
)

func handlerUnfollow(s *state, cmd command, user database.User) error { 
    if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.name)
	}
    url := cmd.args[0]
    feed, err := s.db.GetFeedByURL(context.Background(), url)
    if err != nil {
        return err
    }

    params := database.DeleteFeedFollowParams{
        FeedID: feed.ID,
        UserID: user.ID,
    }
    err = s.db.DeleteFeedFollow(context.Background(), params)
    if err != nil {
        return err
    }

    fmt.Printf("%s unfollowed successfully!\n", feed.Name)

    return nil
}
