# Users Microservice - Complete Authentication System

Microservicio completo de autenticación y gestión de usuarios desarrollado en Go, con registro de usuarios, verificación por email y gestión de roles.

## 📌 Características

- **Registro de usuarios** con validación de datos
- **Verificación de email** con código de 6 dígitos
- **Autenticación JWT** con tokens seguros
- **Gestión de roles** (usuario regular y administrador)
- **Login seguro** con hash SHA-256
- **Reenvío de código** de verificación
- **Emails de bienvenida** automáticos
- **API RESTful** con Gin Framework
- **Base de datos MySQL** con GORM
- **Despliegue fácil** con Docker Compose

---

## 🛠️ Tecnologías

- **Backend:** Go 1.24.1
- **Framework:** Gin
- **ORM:** GORM
- **Base de Datos:** MySQL 8.0
- **Autenticación:** JWT (golang-jwt/jwt)
- **Email:** SMTP nativo de Go
- **Contenedores:** Docker & Docker Compose

---

## 🚀 Cómo levantar el servicio

### Opción 1: Con Docker Compose (Recomendado)

1. Asegúrate de tener Docker y Docker Compose instalados
2. Clona el repositorio
3. En la raíz del proyecto, ejecuta:

```bash
docker-compose up --build
```

El servicio estará disponible en: **http://localhost:8080**

La base de datos MySQL estará en el puerto **3306**.

**Nota:** Por defecto, el servicio funcionará sin configuración SMTP. Los códigos de verificación se mostrarán en la consola para desarrollo.

### Opción 2: Ejecutar localmente

1. Instala Go 1.24.1 o superior
2. Instala MySQL 8.0 y crea una base de datos llamada `users_db`
3. Configura las variables de entorno:

```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=appuser
export DB_PASS=1234
export DB_NAME=users_db
```

4. Navega al directorio `backend` y ejecuta:

```bash
cd backend
go mod download
go run main.go
```

---

## 📧 Configuración de Email (Opcional)

Para habilitar el envío de emails de verificación, configura las siguientes variables de entorno:

```bash
export SMTP_HOST=smtp.gmail.com
export SMTP_PORT=587
export SMTP_USER=your-email@gmail.com
export SMTP_PASS=your-app-password
export SMTP_FROM=your-email@gmail.com
```

### Configurar Gmail para SMTP:

1. Habilita la verificación en dos pasos en tu cuenta de Gmail
2. Genera una "Contraseña de aplicación" en la configuración de seguridad
3. Usa esa contraseña como `SMTP_PASS`

**Nota:** Si no configuras SMTP, el sistema funcionará normalmente pero los códigos de verificación se mostrarán en la consola del servidor en lugar de enviarse por email (ideal para desarrollo).

---

## 📡 Endpoints disponibles

### 🔓 Endpoints Públicos (sin autenticación)

#### 1. Registro de usuario
```http
POST /users/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securePassword123",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response (201 Created):**
```json
{
  "message": "User registered successfully. Please check your email for verification code.",
  "email": "user@example.com"
}
```

**Validaciones:**
- Email válido (formato)
- Password mínimo 6 caracteres
- Campos first_name y last_name requeridos

---

#### 2. Verificar email
```http
POST /users/verify-email
Content-Type: application/json

{
  "email": "user@example.com",
  "code": "123456"
}
```

**Response (200 OK):**
```json
{
  "message": "Email verified successfully. You can now log in.",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Errores comunes:**
- `"invalid verification code"` - Código incorrecto
- `"verification code expired"` - Código expirado (15 minutos)
- `"email already verified"` - Email ya verificado

---

#### 3. Reenviar código de verificación
```http
POST /users/resend-code
Content-Type: application/json

{
  "email": "user@example.com"
}
```

**Response (200 OK):**
```json
{
  "message": "Verification code sent successfully"
}
```

---

#### 4. Iniciar sesión
```http
POST /users/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "name": "John",
  "surname": "Doe"
}
```

**Notas:**
- Requiere que el email esté verificado
- El token JWT expira según configuración

---

### 🔒 Endpoints Protegidos (requieren autenticación)

#### 5. Obtener usuario por ID
```http
GET /users/:id
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "id": 1,
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "is_admin": false,
  "is_verified": true
}
```

---

### 👑 Endpoints de Administrador

#### 6. Verificar token de administrador
```http
GET /users/admin
Authorization: Bearer <admin_token>
```

**Response (200 OK):** Status 200 si el token es válido de admin

---

## 🔐 Sistema de Autenticación Completo

### Flujo de Registro y Verificación:

1. **Registro**: Usuario se registra con email, password, nombre y apellido
2. **Código enviado**: Se genera un código de 6 dígitos válido por 15 minutos
3. **Email enviado**: Se envía el código por email (o se muestra en consola)
4. **Verificación**: Usuario ingresa el código para verificar su email
5. **Cuenta activada**: Usuario puede hacer login normalmente
6. **Token JWT**: Se genera token JWT tras login exitoso

### Gestión de Roles:

El sistema tiene dos niveles de acceso:

1. **Usuario regular** (`IsAdmin: false`)
   - Puede registrarse
   - Debe verificar email
   - Acceso a endpoints protegidos básicos

2. **Administrador** (`IsAdmin: true`)
   - Todos los permisos de usuario regular
   - Acceso a endpoints administrativos
   - Se identifica con token JWT que incluye `IsAdmin: true`

### Seguridad:

- **Passwords**: Hasheados con SHA-256
- **Tokens**: JWT con firma HMAC
- **Códigos**: Aleatorios de 6 dígitos, expiran en 15 minutos
- **Verificación obligatoria**: No se puede hacer login sin verificar email
- **Email único**: No se permiten emails duplicados

---

## 📂 Estructura del proyecto

```
users/
├── backend/
│   ├── app/
│   │   ├── router.go           # Configuración de Gin
│   │   └── url_mappings.go     # Definición de rutas
│   ├── clients/user/
│   │   └── user_clients.go     # Operaciones de base de datos
│   ├── controllers/
│   │   └── user_controller.go  # Controladores HTTP
│   ├── db/
│   │   └── db.go               # Configuración de BD
│   ├── dto/
│   │   └── users_dto.go        # Data Transfer Objects
│   ├── model/
│   │   └── user_model.go       # Modelos de datos
│   ├── services/
│   │   └── user_servicies.go   # Lógica de negocio
│   └── utils/
│       ├── email.go            # Utilidades de email
│       ├── hash.go             # Hash de passwords
│       ├── jwt.go              # Manejo de JWT
│       ├── cors.go             # Configuración CORS
│       └── errors.go           # Manejo de errores
├── db/
│   └── init/                   # Scripts de inicialización
├── docker-compose.yml
└── README.md
```

---

## 🔧 Configuración

### Variables de entorno

#### Base de datos (requeridas):
- `DB_HOST`: Host de la base de datos (default: db)
- `DB_PORT`: Puerto de MySQL (default: 3306)
- `DB_USER`: Usuario de la base de datos (default: appuser)
- `DB_PASS`: Contraseña de la base de datos (default: 1234)
- `DB_NAME`: Nombre de la base de datos (default: users_db)

#### SMTP (opcionales):
- `SMTP_HOST`: Servidor SMTP (ej: smtp.gmail.com)
- `SMTP_PORT`: Puerto SMTP (ej: 587)
- `SMTP_USER`: Usuario de email
- `SMTP_PASS`: Contraseña o app password
- `SMTP_FROM`: Email del remitente

---

## 🗄️ Base de datos

### Tabla: `user_models`

| Campo | Tipo | Descripción |
|-------|------|-------------|
| id | INT | Clave primaria autoincremental |
| email | VARCHAR(100) | Email único |
| password_hash | LONGTEXT | Hash SHA-256 del password |
| first_name | VARCHAR(100) | Nombre |
| last_name | VARCHAR(100) | Apellido |
| is_admin | BOOLEAN | Rol de administrador |
| is_verified | BOOLEAN | Email verificado |
| created_at | TIMESTAMP | Fecha de creación |
| verification_code | VARCHAR(6) | Código de verificación |
| code_expires_at | TIMESTAMP | Expiración del código |

### Tabla: `verification_tokens`

| Campo | Tipo | Descripción |
|-------|------|-------------|
| id | INT | Clave primaria |
| user_id | INT | ID del usuario |
| token | VARCHAR(6) | Código de 6 dígitos |
| expires_at | TIMESTAMP | Fecha de expiración |
| created_at | TIMESTAMP | Fecha de creación |

---

## 🧪 Ejemplo de uso completo

### 1. Registrar nuevo usuario:

```bash
curl -X POST http://localhost:8080/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User"
  }'
```

### 2. Verificar email (con código recibido):

```bash
curl -X POST http://localhost:8080/users/verify-email \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "code": "123456"
  }'
```

### 3. Iniciar sesión:

```bash
curl -X POST http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 4. Usar el token para acceder a recursos protegidos:

```bash
curl -X GET http://localhost:8080/users/1 \
  -H "Authorization: Bearer <tu-token-jwt>"
```

---

## 📮 Colección de Postman

Hemos incluido una colección completa de Postman con todos los endpoints configurados y listos para usar.

### Cómo importar la colección:

1. Abre Postman
2. Click en "Import" (arriba a la izquierda)
3. Selecciona el archivo `Users_Microservice.postman_collection.json`
4. La colección se importará con todos los endpoints

### Características de la colección:

- ✅ **6 endpoints** completamente configurados
- ✅ **Ejemplos de request/response** en cada endpoint
- ✅ **Documentación integrada** en cada request
- ✅ **Variable automática de token**: El endpoint de login guarda automáticamente el token JWT
- ✅ **Autenticación configurada**: Los endpoints protegidos usan automáticamente el token guardado
- ✅ **Ejemplos de respuestas** exitosas incluidos

### Flujo recomendado de prueba:

1. **Register User** - Registra un nuevo usuario
2. Revisa la consola del servidor para obtener el código de verificación
3. **Verify Email** - Verifica el email con el código recibido
4. **Login** - Inicia sesión (guarda automáticamente el token)
5. **Get User By ID** - Obtiene información del usuario (usa el token automáticamente)
6. **Verify Admin Token** - Verifica permisos de admin (solo si eres admin)

### Variables de colección:

- `base_url`: http://localhost:8080 (cambiar para otros ambientes)
- `auth_token`: Se guarda automáticamente tras el login

---

## 📝 Notas de desarrollo

### Modo desarrollo (sin SMTP):

Si no configuras SMTP, el sistema mostrará los códigos de verificación en la consola:

```
=== VERIFICATION CODE FOR test@example.com ===
123456
===============================
```

Esto es útil para desarrollo y testing sin necesidad de configurar un servidor de email.

### Producción:

Para producción, asegúrate de:
1. Configurar correctamente las variables SMTP
2. Usar contraseñas seguras para la BD
3. Configurar un secreto JWT robusto
4. Habilitar HTTPS en el servidor
5. Ajustar las políticas de CORS según tus necesidades

---

## 🐛 Troubleshooting

### "User with email already exists"
- El email ya está registrado en el sistema
- Intenta con otro email o recupera tu cuenta

### "Please verify your email before logging in"
- Necesitas verificar tu email primero
- Revisa tu email o usa `/users/resend-code`

### "Verification code expired"
- El código es válido por 15 minutos
- Solicita un nuevo código con `/users/resend-code`

### Email no llega:
- Verifica la configuración SMTP
- Revisa la carpeta de spam
- En desarrollo, el código se muestra en la consola

---

## 📄 Licencia

Este es un proyecto académico desarrollado como microservicio de autenticación.

---

## 👥 Contribuir

Este es un microservicio educativo. Para mejoras o sugerencias, por favor abre un issue o pull request.
