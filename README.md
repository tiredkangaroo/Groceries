# checklist

This is a basic checklist app written in golang and react. 

## prerequisites
- golang (v 1.22 was used in this app)
- postgreSQL (v 14 was used in this app)

### development prerequisites
- node (v 21.6.2 was used in this app)
- npm (v 10.2.4 was used in this app)

## using a systemd service file (via systemctl)

The .service file is to run on a Linux machine that uses systemd. .service files belong in `/etc/systemd/system`.

#### modify file

However, before moving it into the directory, edit the relevant .service file as necessary (change the paths to YOUR system's paths).

#### move file
Move the file into the `/etc/systemd/system`.

#### start app
Once you have the modified .service file moved in to the directory, run `systemctl start checklist.service`. 

#### status check
Check its status:

`systemctl status groceries.service`

If necessary, fix errors.

#### visit app
If the app is working, you may now visit it at `[::1]:8030` or via the IP of the computer at port 8030 on any other computer in the network.

#### run app on system start
Run `systemctl enable checklist.service` if you would like it to run on startup.

## development info
db.go is a rough general database manager. It is not specific to this app. db.go only supports strings and integers. It lacks a lot of basic postgres.

main.go defines server routes, establishes a connection to the database, and creates the Product database.

APIServer.go has all the functions to actually handle the server routes (except for /, which is handled by main.go and serves the frontend)

All files in the directory static are served as static files on /.

models.go just defines tables (think django models).

All files in the frontend are the HTML, React (TSX, SWC) code. It uses Vite for the server. 

### commands
Commands that might be needed for development.
#### frontend
install dependencies: `npm i`

run dev server: `npm run dev`

build: `npm run build`

#### backend
build: `go build`

run file: `go run /path/to/file`
