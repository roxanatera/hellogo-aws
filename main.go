package main

import (
    "fmt"
    "net/http"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

// Definimos la variable global Database
var Database = func() (db *gorm.DB) {
    err := godotenv.Load()
    if err != nil {
        panic("Error cargando el archivo .env")
    }

    
    dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_SERVER") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Error conectando a la base de datos")
    }
    fmt.Println("Conexión a MySQL exitosa")
    return db
}()

// Definimos el modelo Categoria
type Categoria struct {
    Id     uint   `json:"id"`
    Nombre string `gorm:"type:varchar(100)" json:"nombre"`
    Slug   string `gorm:"type:varchar(100)" json:"slug"`
}

// Especificamos el nombre de la tabla
func (Categoria) TableName() string {
    return "categorias"
}

// Función para ejecutar migraciones
func Migraciones() {
    err := Database.AutoMigrate(&Categoria{})
    if err != nil {
        fmt.Println("Error al realizar la migración:", err)
        panic(err)
    }
    fmt.Println("Migración completada con éxito")
}

func main() {
    gin.SetMode(gin.ReleaseMode)
    router := gin.Default()

    // Middleware CORS
    router.Use(corsMiddleware())

    // Ejecutamos las migraciones
    Migraciones()

    // Ruta para obtener categorías
    router.GET("/categorias", func(c *gin.Context) {
        var datos []Categoria
        result := Database.Order("id desc").Find(&datos)

        // Logs de depuración
        fmt.Println("🚀 Resultado de la consulta:", result.RowsAffected)
        fmt.Println("🚀 Error en la consulta:", result.Error)

        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener categorías"})
            return
        }

        if len(datos) == 0 {
            c.JSON(http.StatusOK, gin.H{"message": "No hay categorías disponibles"})
            return
        }

        c.JSON(http.StatusOK, datos)
    })

    // Iniciamos el servidor
    router.Run(":" + os.Getenv("PORT"))
}

// Middleware CORS
func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}