package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/chzyer/readline"
	"github.com/evanxg852000/foxdb/internal/core"
	wire "github.com/jeroenrinzema/psql-wire"
)

func main() {
	dirPtr := flag.String("dir", "./data", "database directory path")
	servePtr := flag.Bool("serve", false, "start the database server")
	portPtr := flag.Int("port", 5444, "server port to listen on")
	flag.Parse()

	db, err := core.Open(*dirPtr)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		return
	}

	if !*servePtr {
		runInteractiveMode(db)
		return
	}

	runServerMode(db, *portPtr)
}

func runInteractiveMode(db *core.Database) {
	fmt.Println("Running in interactive mode. Type '\\exit' to quit.")

	config := &readline.Config{
		Prompt:                 "> ",
		DisableAutoSaveHistory: true,
	}

	rl, err := readline.NewEx(config)
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		sqlCommand, err := readSqlInput(rl)
		if err != nil {
			break // exit on error or EOF
		}

		if strings.TrimSpace(sqlCommand) != "" {
			rl.SaveHistory(sqlCommand)
			err = executeCommand(db, sqlCommand)
			if err != nil {
				fmt.Printf("Error executing command: %v\n", err)
			}
		}
	}

}

func runServerMode(db *core.Database, port int) {
	serverAddress := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting server on %s...\n", serverAddress)
	if err := wire.ListenAndServe(serverAddress, createRequestHandler(db)); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func createRequestHandler(db *core.Database) wire.ParseFn {
	return func(ctx context.Context, sqlStmt string) (wire.PreparedStatements, error) {
		// session, exist := wire.GetSession(ctx)
		clientParams := wire.ClientParameters(ctx)
		userName := clientParams["user"]
		dbName := clientParams["database"]
		return wire.Prepared(wire.NewStatement(func(ctx context.Context, writer wire.DataWriter, parameters []wire.Parameter) error {
			fmt.Printf("Connected user: %s\n", userName)
			fmt.Printf("Connected to database: %s\n", dbName)
			fmt.Printf("Received query: %s\n", sqlStmt)
			fmt.Printf("Received parameters: %v\n", parameters)
			return writer.Complete("OK")
		})), nil
	}
}

func readSqlInput(rl *readline.Instance) (string, error) {
	sqlCommand := ""
	for {
		if sqlCommand == "" {
			rl.SetPrompt("> ")
		} else {
			rl.SetPrompt("... ")
		}

		line, err := rl.Readline()
		if err != nil { // io.EOF
			return "", err
		}

		if line == "" {
			continue
		}

		// special commands start with '\'
		if sqlCommand == "" && strings.HasPrefix(line, "\\") {
			return line, nil
		}

		sqlCommand += line
		if strings.HasSuffix(line, ";") {
			return sqlCommand, nil
		}
		sqlCommand += "\n"
	}

}

func executeCommand(db *core.Database, sqlCommand string) error {
	fmt.Printf("Executing SQL command:\n%s\n", sqlCommand)
	data, err := db.Run(context.TODO(), sqlCommand)
	if err != nil {
		return err
	}
	//TODO:
	fmt.Printf("Command result:\n%s\n", data)
	return nil
}
