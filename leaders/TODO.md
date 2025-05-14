TODO:

CONSUME API:
  - read from redis
  - parse payload
  - calculate stuff
  - insert into DB

DB:
  - create DB view

READ API:
  - fetch leaderboard from DB



* goroutine to listen for events and send them to a channel
* another goroutine to read from channnel and process events
* thing to listen for signals and shutdown gracefully

