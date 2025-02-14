# README

Postgresql is running locally.

Connection string to the `gator` database is: `psql "postgres://postgres:<password>@localhost:5433/gator"`
Password is my linux user pw. Somehow it was not changed successfully with the command `sudo passwd postgres`.

Connect to database vie `psql`: `sudo -u postgres psql`

## TODO

- [ ] Create a README.md file in the root of your repo if you don't have one already. (You should be tracking your changes with Git.)
- [ ] Explain to the user that they'll need Postgres and Go installed to run the program.
- [ ] Explain to the user how to install the gator CLI using go install.
- [ ] Explain to the user how to set up the config file and run the program. Tell them about a few of the commands they can run.
