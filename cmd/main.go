package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/chzyer/readline"
	wire "github.com/jeroenrinzema/psql-wire"
)

func main() {
	servePtr := flag.Bool("serve", false, "start the database server")
	portPtr := flag.Int("port", 5444, "server port to listen on")
	flag.Parse()
	if !*servePtr {
		runInteractiveMode()
		return
	}

	runServerMode(*portPtr)

}

func runInteractiveMode() {
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
			err = executeCommand(sqlCommand)
			if err != nil {
				fmt.Printf("Error executing command: %v\n", err)
			}
		}
	}

}

func runServerMode(port int) {
	serverAddress := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting server on %s...\n", serverAddress)
	if err := wire.ListenAndServe(serverAddress, handleRequest); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func handleRequest(ctx context.Context, sqlStmt string) (wire.PreparedStatements, error) {
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

func executeCommand(command string) error {
	fmt.Printf("Executing SQL command:\n%s\n", command)
	// Here you would add the logic to parse and execute the SQL command
	return nil
}
