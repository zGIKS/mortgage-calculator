# Finanzas Backend - MiVivienda Mortgage Calculator

API REST para el cálculo de hipotecas MiVivienda utilizando el método francés, implementada con arquitectura hexagonal (DDD) y dos bounded contexts: **IAM** y **Mortgage**.

## 📋 Requisitos Previos

- **Go** 1.23 o superior
- **PostgreSQL** 12 o superior
- **Git**

## 🚀 Instalación y Configuración

### 1. Clonar el repositorio

```bash
git clone <repository-url>
cd Backend
```

### 2. Instalar dependencias

```bash
go mod download
```

### 3. Configurar variables de entorno

Crear un archivo `.env` en la raíz del proyecto con las siguientes variables:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=tu_password
DB_NAME=finanzas
DB_SSLMODE=disable

# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost

# Application
APP_ENV=development

# JWT Configuration
JWT_SECRET_KEY=tu-clave-secreta-jwt-cambiar-en-produccion
JWT_ISSUER=finanzas-backend
JWT_EXPIRATION_HRS=24
```

### 4. Crear la base de datos

```bash
# Conectarse a PostgreSQL
psql -U postgres

# Crear la base de datos
CREATE DATABASE finanzas;
\q
```

### 5. Generar documentación Swagger (opcional)

```bash
# Instalar swag CLI si no lo tienes
go install github.com/swaggo/swag/cmd/swag@latest

# Generar documentación
swag init -g cmd/api/main.go -o cmd/api/docs
```

## ▶️ Ejecutar el Proyecto

### Modo desarrollo

```bash
go run cmd/api/main.go
```

### Compilar y ejecutar

```bash
# Compilar
go build -o bin/api cmd/api/main.go

# Ejecutar
./bin/api
```

El servidor estará disponible en:
- **API**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html

## 📚 Documentación de la API

Una vez que el servidor esté corriendo, accede a la documentación interactiva de Swagger:

```
http://localhost:8080/swagger/index.html
```

## 🏗️ Estructura del Proyecto

```
Backend/
├── cmd/
│   └── api/
│       ├── main.go           # Punto de entrada de la aplicación
│       └── docs/             # Documentación Swagger generada
├── internal/
│   ├── iam/                  # Bounded Context: Identity & Access Management
│   │   ├── application/      # Servicios de aplicación y ACL
│   │   ├── domain/           # Entidades, value objects, servicios de dominio
│   │   ├── infrastructure/   # Implementaciones técnicas (DB, JWT)
│   │   └── interfaces/       # Controladores REST y contratos ACL
│   ├── mortgage/             # Bounded Context: Mortgage Calculator
│   │   ├── application/      # Servicios de aplicación y ACL
│   │   ├── domain/           # Entidades, value objects, servicios de dominio
│   │   ├── infrastructure/   # Implementaciones técnicas (DB)
│   │   └── interfaces/       # Controladores REST y middleware
│   └── shared/               # Código compartido
│       └── infrastructure/   # Config, persistencia compartida
├── .env                      # Variables de entorno
├── go.mod                    # Dependencias de Go
└── README.md                 # Este archivo
```

## 🔐 Autenticación

La API utiliza JWT (JSON Web Tokens) para autenticación. Para acceder a los endpoints protegidos:


## 🛠️ Tecnologías Utilizadas

- **Gin** - Framework web
- **GORM** - ORM para PostgreSQL
- **JWT** - Autenticación y autorización
- **Swagger** - Documentación de API
- **bcrypt** - Hash de contraseñas
- **godotenv** - Gestión de variables de entorno

## 📝 Arquitectura

El proyecto sigue los principios de:
- **Domain-Driven Design (DDD)**
- **Arquitectura Hexagonal**
- **Bounded Contexts**
- **Anti-Corruption Layer (ACL)** para comunicación entre contextos


