package main // MISMO PAQUETE que main.go

import (
	"context"
	"database/sql"

	// "os"
	"testing"
	// "time"
	sqlc "tp2/db" // generado por sqlc

	_ "github.com/lib/pq"
)

// setupTestDB crea una conexión a la base de datos de prueba
func setupTestDB(t *testing.T) *sql.DB {
	// Usa la MISMA cadena de conexión que en main.go pero preferiblemente a una BD de prueba
	connStr := "user=videla password='XYZ' dbname=tarjetasdb port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr) //creo que aca va db-1 en vez de postgres
	if err != nil {
		t.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	// Verificar la conexión
	err = db.Ping()
	if err != nil {
		t.Fatalf("Error al hacer ping a la base de datos: %v", err)
	}

	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
	return db
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
	// Borrar primero Tarjeta
	if _, err := db.Exec(`DELETE FROM Tarjeta`); err != nil {
		t.Fatalf("Error al limpiar Tarjeta: %v", err)
	}

	// Luego borrar Usuario
	if _, err := db.Exec(`DELETE FROM Usuario`); err != nil {
		t.Fatalf("Error al limpiar Usuario: %v", err)
	}

	// Luego borrar Tema
	if _, err := db.Exec(`DELETE FROM Tema`); err != nil {
		t.Fatalf("Error al limpiar Tema: %v", err)
	}
}

func TestUsuarioRepository_CRUD(t *testing.T) {

	db := setupTestDB(t)       // Configurar la base de datos de prueba
	defer db.Close()           // ← Garantiza que la conexión se CERRARÁ al final
	defer cleanupTestDB(t, db) // Limpiar después de la prueba

	// Instanciar el repositorio EXACTAMENTE como en main.go
	queries := sqlc.New(db) // ← Así como lo haces en main.go
	ctx := context.Background()

	t.Run("CreateUsuario", func(t *testing.T) {
		cleanupTestDB(t, db) // ← LIMPIAR AL INICIO de cada subtest

		// Crear usuario usando la MISMA función que en main.go
		createdUsuario, err := queries.CreateUsuario(ctx, sqlc.CreateUsuarioParams{
			NombreUsuario: "PRUEBA 1",
			Email:         "test1@example.com",
			Contrasena:    "securepassword",
		})

		if err != nil {
			t.Fatalf("CreateUsuario failed: %v", err)
		}

		if createdUsuario.IDUsuario == 0 {
			t.Error("Expected auto-generated IDUsuario")
		}

		if createdUsuario.NombreUsuario != "PRUEBA 1" {
			t.Errorf("Expected NombreUsuario 'PRUEBA 1', got '%s'", createdUsuario.NombreUsuario)
		}
	})

	t.Run("GetUsuarioByID", func(t *testing.T) {
		cleanupTestDB(t, db) // ← LIMPIAR AL INICIO de cada subtest

		// Primero crear un usuario
		createdUsuario, err := queries.CreateUsuario(ctx, sqlc.CreateUsuarioParams{
			NombreUsuario: "PRUEBA 2",
			Email:         "test2@example.com",
			Contrasena:    "securepassword",
		})
		if err != nil {
			t.Fatalf("Setup failed: %v", err)
		}

		// Luego obtenerlo por IDUsuario
		retrievedUsuario, err := queries.GetUsuarioById(ctx, createdUsuario.IDUsuario)
		if err != nil {
			t.Fatalf("GetUsuarioByID failed: %v", err)
		}

		if retrievedUsuario.NombreUsuario != "PRUEBA 2" {
			t.Errorf("Expected NombreUsuario 'PRUEBA 2', got '%s'", retrievedUsuario.NombreUsuario)
		}
	})

	t.Run("UpdateUsuario", func(t *testing.T) {
		cleanupTestDB(t, db) // ← LIMPIAR AL INICIO de cada subtest

		// Crea un usuario de prueba en la base de datos
		createdUsuario, err := queries.CreateUsuario(ctx, sqlc.CreateUsuarioParams{
			NombreUsuario: "PRUEBA 3",
			Email:         "test3@example.com",
		})
		// Si falla la creación, termina el subtest con error
		if err != nil {
			t.Fatalf("UpdateUsuario failed: %v", err)
		}

		// Actualiza el usuario creado, cambiando el nombre
		err = queries.UpdateUsuario(ctx, sqlc.UpdateUsuarioParams{
			IDUsuario:     createdUsuario.IDUsuario, // IDUsuario del usuario a actualizar
			NombreUsuario: "PRUEBA 3 ACTUALIZADO",   // Nuevo nombre
			Email:         createdUsuario.Email,     // Email permanece igual
		})
		if err != nil {
			// Si falla la actualización, termina el subtest con error
			t.Fatalf("GetUsuario failed: %v", err)
		}

		// Obtiene el usuario actualizado desde la base de datos
		updatedUsuario, err := queries.GetUsuarioById(ctx, createdUsuario.IDUsuario)
		if err != nil {
			t.Fatalf("GetUsuario failed: %v", err)
		}
		// Verifica que el nombre del usuario haya sido actualizado correctamente
		if updatedUsuario.NombreUsuario != "PRUEBA 3 ACTUALIZADO" {
			t.Errorf("Expected updated NombreUsuario, got '%s'", updatedUsuario.NombreUsuario)
		}
	})

	t.Run("DeleteUsuario", func(t *testing.T) {
		cleanupTestDB(t, db) // ← LIMPIAR AL INICIO de cada subtest

		// Primero crear un usuario
		createdUsuario, err := queries.CreateUsuario(ctx, sqlc.CreateUsuarioParams{
			NombreUsuario: "PRUEBA 4",
			Email:         "test4ejemplo.com",
		})
		if err != nil {
			t.Fatalf("Setup failed: %v", err)
		}

		// Luego eliminarlo
		err = queries.DeleteUsuario(ctx, createdUsuario.IDUsuario)
		if err != nil {
			t.Fatalf("DeleteUsuario failed: %v", err)
		}

		// Intentar obtener el usuario eliminado
		_, err = queries.GetUsuarioById(ctx, createdUsuario.IDUsuario)
		if err == nil {
			t.Fatal("Expected error when getting deleted Usuario, got nil")
		}
	})

	t.Run("ListUsuarios", func(t *testing.T) {
		cleanupTestDB(t, db) // ← LIMPIAR AL INICIO de cada subtest

		// Crear varios usuarios de prueba
		UsuariosToCreate := []sqlc.CreateUsuarioParams{
			{NombreUsuario: "Usuario A", Email: "a"},
			{NombreUsuario: "Usuario B", Email: "b"},
			{NombreUsuario: "Usuario C", Email: "c"},
		}

		// Inserta los usuarios en la base de datos
		for _, u := range UsuariosToCreate {
			_, err := queries.CreateUsuario(ctx, u)
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
		}

		// Llama a ListUsuarios para obtener todos los usuarios
		Usuarios, err := queries.ListUsuarios(ctx)
		if err != nil {
			t.Fatalf("ListUsuarios failed: %v", err)
		}

		// Verifica que el número de usuarios obtenidos sea correcto
		// No lo pide pero esta bueno probarlo
		if len(Usuarios) != len(UsuariosToCreate) {
			t.Errorf("Expected %d Usuarios, got %d", len(UsuariosToCreate), len(Usuarios))
		}

		// Verifica que los nombres y emails creados estén en la lista recuperada
		for i, u := range UsuariosToCreate {
			if Usuarios[i].NombreUsuario != u.NombreUsuario || Usuarios[i].Email != u.Email {
				t.Errorf("Expected Usuario %d to be %+v, got %+v", i, u, Usuarios[i])
			}
		}
	})
}

func TestTarjetaRepository_CRUD(t *testing.T) {

	db := setupTestDB(t)       // Configurar la base de datos de prueba
	defer db.Close()           // ← Garantiza que la conexión se CERRARÁ al final
	defer cleanupTestDB(t, db) // Limpiar después de la prueba

	// Instanciar el repositorio EXACTAMENTE como en main.go
	queries := sqlc.New(db) // ← Así como lo haces en main.go
	ctx := context.Background()

	t.Run("CreateTarjeta", func(t *testing.T) {
		cleanupTestDB(t, db) // ← LIMPIAR AL INICIO de cada subtest

		// Crear tarjeta usando la MISMA función que en main.go
		createdTarjeta, err := queries.CreateTarjeta(ctx, sqlc.CreateTarjetaParams{
			Pregunta:  "¿Cuál es la capital de Francia?",
			Respuesta: "París",
			OpcionA:   "Berlín",
			OpcionB:   "Madrid",
			OpcionC:   "París",
			IDTema:    1,
		})

		if err != nil {
			t.Fatalf("CreateTarjeta failed: %v", err)
		}

		if createdTarjeta.IDTarjeta == 0 {
			t.Error("Expected auto-generated IDTarjeta")
		}

		if createdTarjeta.Pregunta != "¿Cuál es la capital de Francia?" {
			t.Errorf("Expected pregunta '¿Cuál es la capital de Francia?', got '%s'", createdTarjeta.Pregunta)
		}
	})

	t.Run("GetTarjetaByID", func(t *testing.T) {
		cleanupTestDB(t, db) // ← LIMPIAR AL INICIO de cada subtest

		// Primero crear un Tarjeta
		createdTarjeta, err := queries.CreateTarjeta(ctx, sqlc.CreateTarjetaParams{
			Pregunta:  "¿Cuál es la capital de España?",
			Respuesta: "Madrid",
			OpcionA:   "Barcelona",
			OpcionB:   "Madrid",
			OpcionC:   "Valencia",
			IDTema:    1,
		})
		if err != nil {
			t.Fatalf("Setup failed: %v", err)
		}

		// Luego obtenerlo por IDUsuario
		retrievedTarjeta, err := queries.GetTarjetaById(ctx, createdTarjeta.IDTarjeta)
		if err != nil {
			t.Fatalf("GetTarjetaByID failed: %v", err)
		}

		if retrievedTarjeta.Pregunta != "¿Cuál es la capital de España?" {
			t.Errorf("Expected Pregunta '¿Cuál es la capital de España?', got '%s'", retrievedTarjeta.Pregunta)
		}
	})

	t.Run("UpdateTarjeta", func(t *testing.T) {
		cleanupTestDB(t, db) // ← LIMPIAR AL INICIO de cada subtest

		// Crea un Tarjeta de prueba en la base de datos
		createdTarjeta, err := queries.CreateTarjeta(ctx, sqlc.CreateTarjetaParams{
			Pregunta:  "¿Cuál es la capital de Italia?",
			Respuesta: "Roma",
			OpcionA:   "Milán",
			OpcionB:   "Roma",
			OpcionC:   "Nápoles",
			IDTema:    1,
		})
		// Si falla la creación, termina el subtest con error
		if err != nil {
			t.Fatalf("UpdateTarjeta failed: %v", err)
		}

		// Actualiza el Tarjeta creado, cambiando el nombre
		err = queries.UpdateTarjeta(ctx, sqlc.UpdateTarjetaParams{
			IDTarjeta: createdTarjeta.IDTarjeta,                       // IDTarjeta del Tarjeta a actualizar
			Pregunta:  "¿Cuál es la capital de Italia? (Actualizada)", // Nueva pregunta
			Respuesta: createdTarjeta.Respuesta,                       // Respuesta permanece igual
			OpcionA:   createdTarjeta.OpcionA,                         // OpcionA permanece igual
			OpcionB:   createdTarjeta.OpcionB,                         // OpcionB permanece igual
			OpcionC:   createdTarjeta.OpcionC,                         // OpcionC permanece igual
			IDTema:    createdTarjeta.IDTema,                          // IDTema permanece igual
		})
		if err != nil {
			// Si falla la actualización, termina el subtest con error
			t.Fatalf("GetTarjeta failed: %v", err)
		}

		// Obtiene el Tarjeta actualizado desde la base de datos
		updatedTarjeta, err := queries.GetTarjetaById(ctx, createdTarjeta.IDTarjeta)
		if err != nil {
			t.Fatalf("GetTarjeta failed: %v", err)
		}
		// Verifica que el nombre del Tarjeta haya sido actualizado correctamente
		if updatedTarjeta.Pregunta != "¿Cuál es la capital de Italia? (Actualizada)" {
			t.Errorf("Expected updated Pregunta, got '%s'", updatedTarjeta.Pregunta)
		}
	})

	t.Run("DeleteTarjeta", func(t *testing.T) {
		cleanupTestDB(t, db) // ← LIMPIAR AL INICIO de cada subtest

		// Primero crear un Tarjeta
		createdTarjeta, err := queries.CreateTarjeta(ctx, sqlc.CreateTarjetaParams{
			Pregunta:  "¿Cuál es la capital de Alemania?",
			Respuesta: "Berlín",
			OpcionA:   "Múnich",
			OpcionB:   "Berlín",
			OpcionC:   "Hamburgo",
			IDTema:    1,
		})
		if err != nil {
			t.Fatalf("Setup failed: %v", err)
		}

		// Luego eliminarlo
		err = queries.DeleteTarjeta(ctx, createdTarjeta.IDTarjeta)
		if err != nil {
			t.Fatalf("DeleteTarjeta failed: %v", err)
		}

		// Intentar obtener el Tarjeta eliminado
		_, err = queries.GetTarjetaById(ctx, createdTarjeta.IDTarjeta)
		if err == nil {
			t.Fatal("Expected error when getting deleted Tarjeta, got nil")
		}
	})

	t.Run("ListTarjetas", func(t *testing.T) {
		cleanupTestDB(t, db) // ← LIMPIAR AL INICIO de cada subtest

		// Crear varios Tarjetas de prueba
		TarjetasToCreate := []sqlc.CreateTarjetaParams{
			{Pregunta: "Pregunta A", Respuesta: "Respuesta A", OpcionA: "Opcion A1", OpcionB: "Opcion A2", OpcionC: "Opcion A3", IDTema: 1},
			{Pregunta: "Pregunta B", Respuesta: "Respuesta B", OpcionA: "Opcion B1", OpcionB: "Opcion B2", OpcionC: "Opcion B3", IDTema: 1},
			{Pregunta: "Pregunta C", Respuesta: "Respuesta C", OpcionA: "Opcion C1", OpcionB: "Opcion C2", OpcionC: "Opcion C3", IDTema: 1},
		}

		// Inserta los Tarjetas en la base de datos
		for _, u := range TarjetasToCreate {
			_, err := queries.CreateTarjeta(ctx, u)
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
		}

		// Llama a ListTarjetas para obtener todos los Tarjetas
		Tarjetas, err := queries.ListTarjetas(ctx)
		if err != nil {
			t.Fatalf("ListTarjetas failed: %v", err)
		}

		// Verifica que el número de Tarjetas obtenidos sea correcto
		// No lo pide pero esta bueno probarlo
		if len(Tarjetas) != len(TarjetasToCreate) {
			t.Errorf("Expected %d Tarjetas, got %d", len(TarjetasToCreate), len(Tarjetas))
		}

		tarjetasMap := make(map[string]sqlc.Tarjetum)
		for _, tarjeta := range Tarjetas {
			tarjetasMap[tarjeta.Pregunta] = tarjeta
		}

		for _, u := range TarjetasToCreate {
			tarjeta, ok := tarjetasMap[u.Pregunta]
			if !ok {
				t.Errorf("Tarjeta %s no encontrada", u.Pregunta) // <-- 't' sigue siendo *testing.T
				continue
			}
			if tarjeta.Respuesta != u.Respuesta || tarjeta.OpcionA != u.OpcionA || tarjeta.OpcionB != u.OpcionB || tarjeta.OpcionC != u.OpcionC || tarjeta.IDTema != u.IDTema {
				t.Errorf("Expected Tarjeta %+v, got %+v", u, tarjeta)
			}
		}

	})
}

func TestTemaRepository_CRUD(t *testing.T) {

	db := setupTestDB(t)       // Configurar la base de datos de prueba
	defer db.Close()           // ← Garantiza que la conexión se CERRARÁ al final
	defer cleanupTestDB(t, db) // Limpiar después de la prueba

	// Instanciar el repositorio EXACTAMENTE como en main.go
	queries := sqlc.New(db) // ← Así como lo haces en main.go
	ctx := context.Background()

	t.Run("CreateTema", func(t *testing.T) {
		cleanupTestDB(t, db) // ← LIMPIAR AL INICIO de cada subtest

		// Crear tema usando la MISMA función que en main.go
		createdTema, err := queries.CreateTema(ctx, "PRUEBA 1")

		if err != nil {
			t.Fatalf("CreateTema failed: %v", err)
		}

		if createdTema.IDTema == 0 {
			t.Error("Expected auto-generated IDTema")
		}

		if createdTema.NombreTema != "PRUEBA 1" {
			t.Errorf("Expected NombreTema 'PRUEBA 1', got '%s'", createdTema.NombreTema)
		}
	})

	t.Run("GetTemaByID", func(t *testing.T) {
		cleanupTestDB(t, db) // ← LIMPIAR AL INICIO de cada subtest

		// Primero crear un tema
		createdTema, err := queries.CreateTema(ctx, "PRUEBA 2")
		if err != nil {
			t.Fatalf("Setup failed: %v", err)
		}

		// Luego obtenerlo por IDTema
		retrievedTema, err := queries.GetTemaById(ctx, createdTema.IDTema)
		if err != nil {
			t.Fatalf("GetTemaByID failed: %v", err)
		}

		if retrievedTema.NombreTema != "PRUEBA 2" {
			t.Errorf("Expected NombreTema 'PRUEBA 2', got '%s'", retrievedTema.NombreTema)
		}
	})

	t.Run("UpdateTema", func(t *testing.T) {
		cleanupTestDB(t, db) // ← LIMPIAR AL INICIO de cada subtest

		// Crea un tema de prueba en la base de datos
		createdTema, err := queries.CreateTema(ctx, "PRUEBA 3")
		// Si falla la creación, termina el subtest con error
		if err != nil {
			t.Fatalf("UpdateTema failed: %v", err)
		}

		// Actualiza el tema creado, cambiando el nombre
		err = queries.UpdateTema(ctx, sqlc.UpdateTemaParams{
			IDTema:     createdTema.IDTema,     // IDTema del tema a actualizar
			NombreTema: "PRUEBA 3 ACTUALIZADO", // Nuevo nombre
		})
		if err != nil {
			// Si falla la actualización, termina el subtest con error
			t.Fatalf("GetTema failed: %v", err)
		}

		// Obtiene el tema actualizado desde la base de datos
		updatedTema, err := queries.GetTemaById(ctx, createdTema.IDTema)
		if err != nil {
			t.Fatalf("GetTema failed: %v", err)
		}
		// Verifica que el nombre del tema haya sido actualizado correctamente
		if updatedTema.NombreTema != "PRUEBA 3 ACTUALIZADO" {
			t.Errorf("Expected updated NombreTema, got '%s'", updatedTema.NombreTema)
		}
	})

	t.Run("DeleteTema", func(t *testing.T) {
		cleanupTestDB(t, db) // ← LIMPIAR AL INICIO de cada subtest

		// Primero crear un tema
		createdTema, err := queries.CreateTema(ctx, "PRUEBA 4")
		if err != nil {
			t.Fatalf("Setup failed: %v", err)
		}

		// Luego eliminarlo
		err = queries.DeleteTema(ctx, createdTema.IDTema)
		if err != nil {
			t.Fatalf("DeleteTema failed: %v", err)
		}

		// Intentar obtener el tema eliminado
		_, err = queries.GetTemaById(ctx, createdTema.IDTema)
		if err == nil {
			t.Fatal("Expected error when getting deleted Tema, got nil")
		}
	})

	t.Run("ListTemas", func(t *testing.T) {
		cleanupTestDB(t, db) // ← LIMPIAR AL INICIO de cada subtest

		// Crear varios temas de prueba
		TemasToCreate := []string{
			"Tema A", "Tema B", "Tema C",
		}

		// Inserta los temas en la base de datos
		for _, tema := range TemasToCreate {
			_, err := queries.CreateTema(ctx, tema)
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
		}

		// Llama a ListTemas para obtener todos los temas
		Temas, err := queries.ListTemas(ctx)
		if err != nil {
			t.Fatalf("ListTemas failed: %v", err)
		}

		// Verifica que el número de temas obtenidos sea correcto
		// No lo pide pero esta bueno probarlo
		if len(Temas) != len(TemasToCreate) {
			t.Errorf("Expected %d Temas, got %d", len(TemasToCreate), len(Temas))
		}

		// Verifica que los nombres y emails creados estén en la lista recuperada
		for i, tema := range TemasToCreate {
			if Temas[i].NombreTema != tema {
				t.Errorf("Expected Tema %d to be %+v, got %+v", i, t, Temas[i])
			}
		}
	})
}

/*# Ejecutar todos los subtests
go test -v
# Output: TestUsuarioRepository_CRUD/CreateUsuario
# TestUsuarioRepository_CRUD/GetUsuarioByID
# TestUsuarioRepository_CRUD/UpdateUsuario

# Ejecutar solo un subtest específico
go test -v -run "TestUsuarioRepository_CRUD/CreateUsuario"*/
