# checklist
basic checklist app written in go and react

## running this application

prerequisites:
- npm 10.8 (used 10.8.2)
- go 1.22 (used 1.22.6)
- postgreSQL (used 14.11 homebrew)

##### step 1:
create a .env file to specify information

Variables required in .env:
```env
# defaults to development if empty
MODE="development" | "production"
POSTGRES_CONNECTION_URI="postgres://username:password@postgresServerIP:portIfNecessary/databaseName"
```

##### step 2:
change working directory to frontend and run the command: `npm i && npm run build && mv dist ../static` (this will build the vite application and move it to the main directory for the application as "static").

##### step 3:
go back to the main directory for application/repository and run the go program:

build with `go build -o checklistapp *.go` and run the program with `./checklistapp`.

It should now be running on port :8030.
