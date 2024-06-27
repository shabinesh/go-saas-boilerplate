migrate-up:
	#go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	~/go/bin/migrate -path migrations -database "sqlite3:///Users/shabinesh/Workshop/transcription/test.db" -verbose up

migrate-down:
	~/go/bin/migrate -path migrations -database "sqlite3:///Users/shabinesh/Workshop/transcription/test.db" -verbose down