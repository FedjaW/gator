package main

import (
	"context"
	"fmt"
)

func handlerUsers(s *state, cmd command) error {
    users, err := s.db.GetUsers(context.Background())
    if err != nil {
        return err
    }
    currentUser := s.cfg.CurrentUserName
    for _, user := range users {
        if currentUser == user.Name {
            fmt.Printf("* %s (current)\n", user.Name)
        } else {
            fmt.Printf("* %s\n", user.Name)
        }
    }
    return nil
}
