package db

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Johnhi19/TreeSpotter_backend/models"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client
var db *sql.DB

func Connect(txtFile string) {
	file, err := os.Open(txtFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())
	scanner.Scan()
	password := strings.TrimSpace(scanner.Text())

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// MySQL connection string: username:password@tcp(host:port)/database
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/treeSpotter", username, password)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MySQL!")
}

func FindAllMeadows() []models.Meadow {
	var meadows []models.Meadow

	rows, err := db.Query("SELECT * FROM Meadow")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var med models.Meadow
		if err := rows.Scan(&med.ID, &med.Name, &med.TreeIds, &med.Size, &med.Location); err != nil {
			panic(err)
		}
		meadows = append(meadows, med)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return meadows
}

func FindAllTreesForMeadow(meadowId int) []models.Tree {
	var trees []models.Tree

	rows, err := db.Query("SELECT * FROM Tree WHERE MeadowId = ?", meadowId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var tree models.Tree
		if err := rows.Scan(&tree.ID, &tree.Type, &tree.Age, &tree.MeadowId, &tree.Position); err != nil {
			panic(err)
		}
		trees = append(trees, tree)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return trees
}

func InsertOneMeadow(meadow models.Meadow) any {
	result, err := db.Exec("INSERT INTO Meadow (Name, TreeIds, Size, Location) VALUES (?, ?, ?, ?)",
		meadow.Name, meadow.TreeIds, meadow.Size, meadow.Location)
	if err != nil {
		panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted a meadow with ID:", id)
	return id
}

func InsertOneTree(tree models.Tree) any {
	result, err := db.Exec("INSERT INTO Tree (Type, Age, MeadowId, Position) VALUES (?, ?, ?, ?)",
		tree.Type, tree.Age, tree.MeadowId, tree.Position)
	if err != nil {
		panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted a tree with ID:", id)
	return id
}

func FindOneMeadowById(meadowId int) models.Meadow {
	var meadow models.Meadow

	row := db.QueryRow("SELECT * FROM Meadow WHERE ID = ?", meadowId)
	if err := row.Scan(&meadow.ID, &meadow.Name, &meadow.TreeIds, &meadow.Size, &meadow.Location); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No meadow found with ID:", meadowId)
			return meadow
		}
		panic(err)
	}
	fmt.Println("Found meadow:", meadow)
	return meadow
}

func FindOneTreeById(treeId int) models.Tree {
	var tree models.Tree

	row := db.QueryRow("SELECT * FROM Tree WHERE ID = ?", treeId)
	if err := row.Scan(&tree.ID, &tree.Type, &tree.Age, &tree.MeadowId, &tree.Position); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No tree found with ID:", treeId)
			return tree
		}
		panic(err)
	}
	fmt.Println("Found tree:", tree)
	return tree
}

func Disconnect() {
	if err := db.Close(); err != nil {
		log.Fatal("Error closing database connection:", err)
	} else {
		fmt.Println("Database connection closed.")
	}
}
