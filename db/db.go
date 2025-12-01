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
)

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
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/treeSpotter?parseTime=true", username, password)
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

	rows, err := db.Query("SELECT ID, Location, Name, Size, TreeIds FROM Meadow")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var med models.Meadow
		if err := rows.Scan(&med.ID, &med.Location, &med.Name, &med.Size, &med.TreeIds); err != nil {
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
	meadow := FindOneMeadowById(meadowId)

	if len(meadow.TreeIds) == 0 {
		return []models.Tree{}
	}

	var trees []models.Tree

	placeholders := make([]string, len(meadow.TreeIds))
	args := make([]any, len(meadow.TreeIds))
	for i, treeId := range meadow.TreeIds {
		placeholders[i] = "?"
		args[i] = treeId
	}

	query := fmt.Sprintf("SELECT ID, PlantDate, MeadowId, Position, Type FROM Tree WHERE ID IN (%s)",
		strings.Join(placeholders, ","))

	rows, err := db.Query(query, args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var tree models.Tree
		if err := rows.Scan(&tree.ID, &tree.PlantDate, &tree.MeadowId, &tree.Position, &tree.Type); err != nil {
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
	result, err := db.Exec("INSERT INTO Meadow (Location, Name, Size, TreeIds) VALUES (?, ?, ?, ?)",
		meadow.Location, meadow.Name, meadow.Size, meadow.TreeIds)
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

func FindOneMeadowById(meadowId int) models.Meadow {
	var meadow models.Meadow

	row := db.QueryRow("SELECT ID, Location, Name, Size, TreeIds FROM Meadow WHERE ID = ?", meadowId)
	if err := row.Scan(&meadow.ID, &meadow.Location, &meadow.Name, &meadow.Size, &meadow.TreeIds); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No meadow found with ID:", meadowId)
			return meadow
		}
		panic(err)
	}
	fmt.Println("Found meadow:", meadow)
	return meadow
}

func DeleteOneMeadow(meadowId int) error {
	// First, get all tree IDs associated with the meadow
	meadow := FindOneMeadowById(meadowId)
	if meadow.ID == 0 {
		return fmt.Errorf("meadow with ID %d not found", meadowId)
	}

	// Delete all associated trees
	for _, treeId := range meadow.TreeIds {
		if err := deleteTreeOnly(treeId); err != nil {
			fmt.Printf("Warning: Failed to delete tree ID %d: %v\n", treeId, err)
		}
	}

	// Now delete the meadow itself
	result, err := db.Exec("DELETE FROM Meadow WHERE ID = ?", meadowId)
	if err != nil {
		return fmt.Errorf("failed to delete meadow: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no meadow found with ID %d", meadowId)
	}

	fmt.Printf("Deleted meadow with ID: %d\n", meadowId)
	return nil
}

func UpdateMeadowTreeIds(meadowId int, treeId int64, shouldDelete bool) error {
	// Get current meadow
	meadow := FindOneMeadowById(meadowId)

	fmt.Printf("Current TreeIds for meadow %d: %v\n", meadowId, meadow.TreeIds)

	// Check whether to add or remove the tree ID
	if shouldDelete {
		newTreeIds := make([]int, 0, len(meadow.TreeIds))
		for _, id := range meadow.TreeIds {
			if id != int(treeId) {
				newTreeIds = append(newTreeIds, id)
			}
		}
		meadow.TreeIds = newTreeIds
	} else {
		meadow.TreeIds = append(meadow.TreeIds, int(treeId))
	}

	fmt.Printf("New TreeIds for meadow %d: %v\n", meadowId, meadow.TreeIds)

	// Value() method will automatically be called for TreeIds
	result, err := db.Exec("UPDATE Meadow SET TreeIds = ? WHERE ID = ?", meadow.TreeIds, meadowId)
	if err != nil {
		fmt.Printf("ERROR executing UPDATE: %v\n", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("UPDATE affected %d rows\n", rowsAffected)

	if rowsAffected == 0 {
		return fmt.Errorf("no meadow found with ID %d", meadowId)
	}

	fmt.Printf("Successfully updated meadow %d with tree ID: %d\n", meadowId, treeId)
	return nil
}

func InsertOneTree(tree models.Tree) int64 {
	result, err := db.Exec("INSERT INTO Tree (PlantDate, MeadowId, Position, Type) VALUES (?, ?, ?, ?)",
		tree.PlantDate, tree.MeadowId, tree.Position, tree.Type)
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

func FindOneTreeById(treeId int) models.Tree {
	var tree models.Tree

	row := db.QueryRow("SELECT ID, PlantDate, MeadowId, Position, Type FROM Tree WHERE ID = ?", treeId)
	if err := row.Scan(&tree.ID, &tree.PlantDate, &tree.MeadowId, &tree.Position, &tree.Type); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No tree found with ID:", treeId)
			return tree
		}
		fmt.Printf("Type of tree.PlantDate: %T\n", tree.PlantDate)
		panic(err)
	}
	fmt.Println("Found tree:", tree)
	return tree
}

// Deletes only the tree from the database, does not update meadow's TreeIds
func deleteTreeOnly(treeId int) error {
	result, err := db.Exec("DELETE FROM Tree WHERE ID = ?", treeId)
	if err != nil {
		return fmt.Errorf("failed to delete tree: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no tree found with ID %d", treeId)
	}

	fmt.Printf("Deleted tree with ID: %d\n", treeId)
	return nil
}

// Deletes the tree and updates the meadow's TreeIds accordingly
func DeleteOneTree(treeId int) error {
	// First, get the tree to know which meadow it belongs to
	tree := FindOneTreeById(treeId)
	if tree.ID == 0 {
		return fmt.Errorf("tree with ID %d not found", treeId)
	}

	meadowId := tree.MeadowId

	// Delete the tree from the database
	if err := deleteTreeOnly(treeId); err != nil {
		return err
	}

	// Remove tree ID from meadow's TreeIds
	if err := UpdateMeadowTreeIds(meadowId, int64(treeId), true); err != nil {
		fmt.Printf("Warning: Tree deleted but failed to update meadow: %v\n", err)
		return err
	}

	return nil
}

func Disconnect() {
	if err := db.Close(); err != nil {
		log.Fatal("Error closing database connection:", err)
	} else {
		fmt.Println("Database connection closed.")
	}
}
