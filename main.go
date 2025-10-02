package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	sqlc "tp2/db" // Así, no "TP2/db"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "user=videla password='XYZ' dbname=tarjetasdb port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()
	queries := sqlc.New(db)
	ctx := context.Background()
	createdUser, err := queries.CreateUsuario(ctx, // Create
		sqlc.CreateUsuarioParams{
			NombreUsuario: "John Doe",
			Email:         "john.doe@example.com",
			Contrasena:    "securepassword",
		})

	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}
	fmt.Printf("Created user: %+v\n", createdUser)
	user, err := queries.GetUsuarioById(ctx, createdUser.IDUsuario) // Read One
	if err != nil {
		log.Fatalf("failed to get user: %v", err)
	}
	fmt.Printf("Retrieved user: %+v\n", user)
	users, err := queries.ListUsuarios(ctx) // Read Many
	if err != nil {
		log.Fatalf("failed to list users: %v", err)
	}
	fmt.Printf("All users: %+v\n", users)
	err = queries.UpdateUsuario(ctx, sqlc.UpdateUsuarioParams{ // Update
		IDUsuario:     createdUser.IDUsuario,
		NombreUsuario: "Johnny Doe",
		Email:         "johnny.doe@example.com",
	})

	if err != nil {
		log.Fatalf("failed to update user: %v", err)
	}
	fmt.Println("User updated successfully")
	updatedUser, err := queries.GetUsuarioById(ctx, createdUser.IDUsuario)
	if err != nil {
		log.Fatalf("failed to get updated user: %v", err)
	}
	fmt.Printf("Updated user: %+v\n", updatedUser)

	err = queries.DeleteUsuario(ctx, createdUser.IDUsuario) // Delete
	if err != nil {
		log.Fatalf("failed to delete user: %v", err)
	}
	fmt.Println("User deleted successfully")
	_, err = queries.GetUsuarioById(ctx, createdUser.IDUsuario)
	if err == sql.ErrNoRows {
		fmt.Println("User not found after deletion")
	} else if err != nil {
		log.Fatalf("failed to get user after deletion: %v", err)
	}

}
