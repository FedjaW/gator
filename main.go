package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/fedjaw/gator/internal/config"
	"github.com/fedjaw/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
    cfg *config.Config
    db *database.Queries
}

func main() {
    cfg, _ := config.Read()
    db, err := sql.Open("postgres", cfg.DbUrl)
    if err != nil {
        log.Fatalf("error connecting to db: %v", err)
    }
    defer db.Close()
    dbQueries := database.New(db)

    s := state{
        cfg: &cfg,
        db: dbQueries,
    }

    cmds := commands{
        cmds: make(map[string]func(*state, command) error),
    }

    cmds.register("login", handlerLogin)
    cmds.register("register", handlerRegister)
    cmds.register("reset", handlerReset)
    cmds.register("users", handlerUsers)
    cmds.register("agg", handlerAgg)
    cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
    cmds.register("feeds", handlerFeeds)
    cmds.register("follow", middlewareLoggedIn(handlerFollow))
    cmds.register("following", middlewareLoggedIn(handlerFollowing))
    cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
    cmds.register("browse", middlewareLoggedIn(handlerBrowse))

    args := os.Args
    if len(args) < 2 {
        log.Fatalf("error reading arguments")
    }
    commandName := args[1] // args[0] is the program name
    commandArgs := args[2:]

    cmd := command{
        name: commandName,
        args: commandArgs,
    }

    err = cmds.run(&s ,cmd)
    if err != nil {
        log.Fatal(err)
    }
}

type command struct {
    name string
    args []string
}

type commands struct {
    cmds map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
    c.cmds[name] = f
}

func (c *commands) run(state *state, cmd command) error {
    err := c.cmds[cmd.name](state, cmd)

    return err
}
