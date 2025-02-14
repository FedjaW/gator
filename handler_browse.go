package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/fedjaw/gator/internal/database"
)


func handlerBrowse(s *state, cmd command, user database.User) error {
    var limit int32 = 2
    if len(cmd.args) == 1 {
        i, _ := strconv.Atoi(cmd.args[0])
        limit = int32(i)
    }

    fmt.Printf("limit: %d\n", limit)

    getPostsParams := database.GetPostsParams{
            UserID: user.ID,
            Limit: limit,
    }

    posts, err := s.db.GetPosts(context.Background(), getPostsParams)
    if err != nil {
        fmt.Printf("error: %s\n", err)
        return err
    }

    for _, p := range posts {
        fmt.Printf("- %s\n", p.Title)
    }

    return nil
}
