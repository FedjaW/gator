package main

import (
	"context"
	"fmt"

	"github.com/fedjaw/gator/internal/database"
)

func handlerFollowing(s *state, cmd command, user database.User) error {
    follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
    if err != nil {
        return fmt.Errorf("error retrieving feed follows: %w", err)
    }

    fmt.Printf("Feed follows for user %s:\n", user.Name)
    for _, follow := range follows {
        fmt.Printf("* %s\n", follow.FeedName)
    }
    fmt.Println("======================================")

    return nil
}
