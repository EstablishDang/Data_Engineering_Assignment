package main

import (
	proto_mssql_server "app_mssql_server/mssql_server/proto/v1"
	mssql_server "app_mssql_server/mssql_server/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	_ "github.com/microsoft/go-mssqldb"
)

// // Replace with your own connection parameters
// var server = "3.0.135.17"
// var port = 1433
// var user = "admin"
// var password = "lhp21051999_phat"
// var database = "AdventureWorks2008R2"
// var db *sql.DB

func SetupHttpServer(port string, certFile string, keyFile string) {
	var svr mssql_server.Service
	var err error

	mux := runtime.NewServeMux()
	proto_mssql_server.RegisterAppMssqlMgmtServiceHandlerServer(context.Background(), mux, &svr)
	addr := ":" + port

	if certFile == "" || keyFile == "" {
		err = http.ListenAndServe(addr, mux)
		if err != nil {
			log.Printf("[ERROR]_Start HTTP server fail")
			log.Fatal(err)
		}
	} else if certFile != "" && keyFile != "" {
		err = http.ListenAndServeTLS(addr, certFile, keyFile, mux)
		if err != nil {
			log.Printf("[ERROR]_Start HTTPS server fail")
			log.Fatal(err)
		}
	}
}

func main() {
	os.Setenv("TZ", "UTC")
	fmt.Printf("Hello!\n")
	log.SetOutput(os.Stdout)

	mssql_server.InitDB()

	httpPort := "20000"
	SetupHttpServer(httpPort, "", "")
	log.Printf("Start HTTP server at [::]:%s", httpPort)

	// //Test
	// // Create connection string
	// var err error
	// connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
	// 	server, user, password, port, database)

	// // Create connection pool
	// db, err = sql.Open("sqlserver", connString)
	// if err != nil {
	// 	log.Fatal("Error creating connection pool: " + err.Error())
	// }
	// log.Printf("Connected!\n")

	// // Close the database connection pool after program executes
	// defer db.Close()

	// SelectVersion()

	// // Read employees
	// count, err := ReadEmployees()
	// if err != nil {
	// 	log.Fatal("Error reading Employees: ", err.Error())
	// }
	// fmt.Printf("Read %d row(s) successfully.\n", count)
}
