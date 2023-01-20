package main

// > go mod init pluginpricing
// > go get github.com/aws/aws-sdk-go-v2
// > go get github.com/aws/aws-sdk-go-v2/config
// > go get github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk
// > go get github.com/aws/aws-sdk-go-v2/service/rds

import (
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "context"
    "database/sql"
    "fmt"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go-v2/feature/rds/auth"
    _ "github.com/go-sql-driver/mysql"
)

func main() {

    var dbName string "database-1-instance-1"
    var dbUser string "admin"
    var dbHost string = "database-1-instance-1.cl7yuwuzzgbr.us-west-2.rds.amazonaws.com"
    var dbPort int = 3306
    var dbEndpoint string = fmt.Sprintf("%s:%d", dbHost, dbPort)
    var region string = "us-west-2"

    sess := session.Must(session.NewSession())
    creds := sess.Config.Credentials

    port := os.Getenv("PORT")
    if port == "" {
	port = "5000"
    }

    authenticationToken, err := auth.BuildAuthToken(
    	context.TODO(), dbEndpoint, region, dbUser, cfg.Credentials)
    if err != nil {
	    panic("failed to create authentication token: " + err.Error())
    }

    f, _ := os.Create("/var/log/golang/golang-server.log")
    defer f.Close()
    log.SetOutput(f)

    const indexPage = "public/index.html"
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
	    if buf, err := ioutil.ReadAll(r.Body); err == nil {
		log.Printf("Received message: %s\n", string(buf))
	    }
	} else {
	    log.Printf("Serving %s to %s...\n", indexPage, r.RemoteAddr)
	    http.ServeFile(w, r, indexPage)
	}
    })

    http.HandleFunc("/scheduled", func(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST" {
	    log.Printf("Received task %s scheduled at %s\n", r.Header.Get("X-Aws-Sqsd-Taskname"), r.Header.Get("X-Aws-Sqsd-Scheduled-At"))
	}
    })

    log.Printf("Listening on port %s\n\n", port)
    http.ListenAndServe(":"+port, nil)
}
