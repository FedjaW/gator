package main

import (
	"context"
	"fmt"
	"log"
)

func handlerReset(s *state, cmd command) error {
    err := s.db.ResetUser(context.Background())
    if err != nil {
        log.Fatal(err)
        return fmt.Errorf("couldn't delete users: %w", err)
    }
    fmt.Println("Successfully removed all users")
    return nil
}
