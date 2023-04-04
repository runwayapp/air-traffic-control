package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var db *sql.DB

type Command struct {
	Id           string
	Organization string
	Repository   string
	Name         string
	Data         string
	Created_at   string
	Updated_at   string
}

func main() {
	var err error
	// Load in the `.env` file in development
	if os.Getenv("ENV") != "production" {
		err = godotenv.Load()
		if err != nil {
			log.Fatal("failed to load env", err)
		}
	}

	// Open a connection to the database
	db, err = sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatal("failed to open db connection", err)
	}

	// Build router & define routes
	router := gin.Default()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	router.GET("/commands/:org/:repo", GetRepoCommands)
	router.GET("/commands/:org/:repo/:commandId", GetSingleCommand)
	router.POST("/commands", CreateCommand)
	router.PUT("/commands/:org/:repo/:commandId", UpdateCommand)
	router.DELETE("/commands/:org/:repo/:commandId", DeleteCommand)

	// Run the router
	router.Run()
}

func GetRepoCommands(c *gin.Context) {
	org := c.Param("org")
	org = strings.ReplaceAll(org, "/", "")
	repo := c.Param("repo")
	repo = strings.ReplaceAll(repo, "/", "")

	query := `SELECT * FROM commands WHERE organization = ? AND repository = ?`
	res, err := db.Query(query, org, repo)
	defer res.Close()
	if err != nil {
		msg, _ := fmt.Printf("(GetCommands) db.Query %s", err)
		panic(msg)
	}

	commands := []Command{}
	for res.Next() {
		var command Command
		err := res.Scan(&command.Id, &command.Organization, &command.Repository, &command.Name, &command.Data, &command.Created_at, &command.Updated_at)
		if err != nil {
			msg, _ := fmt.Printf("(GetCommands) res.Scan %s", err)
			panic(msg)
		}
		commands = append(commands, command)
	}

	c.JSON(http.StatusOK, commands)
}

func GetSingleCommand(c *gin.Context) {
	org := c.Param("org")
	org = strings.ReplaceAll(org, "/", "")
	repo := c.Param("repo")
	repo = strings.ReplaceAll(repo, "/", "")
	commandId := c.Param("commandId")
	commandId = strings.ReplaceAll(commandId, "/", "")

	var command Command
	query := `SELECT * FROM commands WHERE id = ? AND organization = ? AND repository = ?`
	err := db.QueryRow(query, commandId, org, repo).Scan(&command.Id, &command.Organization, &command.Repository, &command.Name, &command.Data, &command.Created_at, &command.Updated_at)
	if err != nil {
		msg, _ := fmt.Printf("(GetSingleCommand) db.Exec %s", err)
		panic(msg)
	}

	c.JSON(http.StatusOK, command)
}

func CreateCommand(c *gin.Context) {
	id := uuid.New().String()

	var newCommand Command
	err := c.BindJSON(&newCommand)
	if err != nil {
		msg, _ := fmt.Printf("(CreateCommand) c.BindJSON %s", err)
		panic(msg)
	}

	// add id to newCommand
	newCommand.Id = id

	// Check all required inputs
	if newCommand.Organization == "" || newCommand.Repository == "" || newCommand.Name == "" || newCommand.Data == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization, repository, name, and data are required"})
		return
	}

	query := `INSERT INTO commands (id, organization, repository, name, data) VALUES (?, ?, ?, ?, ?)`
	res, err := db.Exec(query, newCommand.Id, newCommand.Organization, newCommand.Repository, newCommand.Name, newCommand.Data)
	if err != nil {
		msg, _ := fmt.Printf("(CreateCommand) db.Exec %s", err)
		panic(msg)
	}

	_, err = res.LastInsertId()

	if err != nil {
		msg, _ := fmt.Printf("(CreateCommand) res.LastInsertId %s", err)
		panic(msg)
	}

	c.JSON(http.StatusOK, newCommand)
}

func UpdateCommand(c *gin.Context) {
	var updates Command
	err := c.BindJSON(&updates)
	if err != nil {
		msg, _ := fmt.Printf("(UpdateCommand) c.BindJSON %s", err)
		panic(msg)
	}

	org := c.Param("org")
	org = strings.ReplaceAll(org, "/", "")
	repo := c.Param("repo")
	repo = strings.ReplaceAll(repo, "/", "")
	commandId := c.Param("commandId")
	commandId = strings.ReplaceAll(commandId, "/", "")

	// check all required inputs
	if updates.Name == "" || updates.Data == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and data are required"})
		return
	}

	query := `UPDATE commands SET name = ?, data = ? WHERE id = ? AND organization = ? AND repository = ?`
	result, err := db.Exec(query, updates.Name, updates.Data, commandId, org, repo)
	if err != nil {
		msg, _ := fmt.Printf("(UpdateCommand) db.Exec %s", err)
		panic(msg)
	}

	// if no rows were affected, return an error
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		msg, _ := fmt.Printf("(DeleteCommand) result.RowsAffected %s", err)
		panic(msg)
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "command not found"})
		return
	}

	c.Status(http.StatusOK)
}

func DeleteCommand(c *gin.Context) {
	org := c.Param("org")
	org = strings.ReplaceAll(org, "/", "")
	repo := c.Param("repo")
	repo = strings.ReplaceAll(repo, "/", "")
	commandId := c.Param("commandId")
	commandId = strings.ReplaceAll(commandId, "/", "")

	// if an org or repo is not provided, return an error
	if org == "" || repo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "org and repo are required"})
		return
	}

	// if a commandId is not provided, return an error
	if commandId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "commandId is required"})
		return
	}

	query := `DELETE FROM commands WHERE id = ? AND organization = ? AND repository = ?`
	result, err := db.Exec(query, commandId, org, repo)
	if err != nil {
		msg, _ := fmt.Printf("(DeleteCommand) db.Exec %s", err)
		panic(msg)
	}

	// if no rows were affected, return an error
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		msg, _ := fmt.Printf("(DeleteCommand) result.RowsAffected %s", err)
		panic(msg)
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "command not found"})
		return
	}

	c.Status(http.StatusOK)
}
