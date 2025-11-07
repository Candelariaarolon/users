# API Documentation - Users Microservice

## Dominio

Este microservicio gestiona todo lo relacionado con **usuarios y autenticación** en el sistema. Proporciona funcionalidades para:

- Registro de usuarios con verificación por email
- Autenticación mediante JWT (Access Token y Refresh Token)
- Gestión de roles (usuarios regulares y administradores)
- Verificación de identidad mediante códigos de 6 dígitos enviados por correo
- Renovación de tokens de acceso

### Conceptos Clave del Dominio

**Usuario**: Entidad principal que representa a una persona en el sistema. Cada usuario tiene:
- Email único (utilizado como identificador de login)
- Nombre y apellido
- Estado de verificación (verificado/no verificado)
- Rol (usuario regular o administrador)

**Verificación de Email**: Proceso de confirmación de identidad mediante código de 6 dígitos enviado al email del usuario. Los códigos tienen una validez temporal.

**Autenticación JWT**: Sistema de autenticación basado en tokens:
- **Access Token**: Token de corta duración para acceder a recursos protegidos
- **Refresh Token**: Token de larga duración para renovar el access token sin requerir login

**Roles**:
- **Usuario Regular**: Acceso estándar al sistema
- **Administrador**: Permisos elevados para gestión de usuarios

---

## Base URL

```
http://localhost:8080
```

---

## CORS Configuration

El servicio acepta requests desde:
- `http://localhost:3000`

Headers permitidos:
- `Origin`, `Content-Type`, `Authorization`

---

## Endpoints

### 1. Registro de Usuario

**POST** `/users/register`

Crea una nueva cuenta de usuario y envía un código de verificación al email proporcionado.

**Acceso**: Público (no requiere autenticación)

**Request Body**:
```json
{
  "email": "usuario@ejemplo.com",
  "password": "mipassword123",
  "first_name": "Juan",
  "last_name": "Pérez"
}
```

**Validaciones**:
- `email`: requerido, formato email válido
- `password`: requerido, mínimo 6 caracteres
- `first_name`: requerido
- `last_name`: requerido

**Response (201 Created)**:
```json
{
  "message": "Usuario registrado exitosamente. Verifica tu email.",
  "email": "usuario@ejemplo.com"
}
```

**Errores**:
- `400 Bad Request`: Datos inválidos o email ya registrado

---

### 2. Verificar Email

**POST** `/users/verify-email`

Verifica el email del usuario mediante el código de 6 dígitos enviado al registrarse.

**Acceso**: Público

**Request Body**:
```json
{
  "email": "usuario@ejemplo.com",
  "code": "123456"
}
```

**Validaciones**:
- `email`: requerido, formato email válido
- `code`: requerido, exactamente 6 caracteres

**Response (200 OK)**:
```json
{
  "message": "Email verificado exitosamente",
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Errores**:
- `400 Bad Request`: Código inválido o expirado

---

### 3. Reenviar Código de Verificación

**POST** `/users/resend-code`

Reenvía un nuevo código de verificación al email del usuario.

**Acceso**: Público

**Request Body**:
```json
{
  "email": "usuario@ejemplo.com"
}
```

**Validaciones**:
- `email`: requerido, formato email válido

**Response (200 OK)**:
```json
{
  "message": "Verification code sent successfully"
}
```

**Errores**:
- `400 Bad Request`: Email no encontrado o ya verificado

---

### 4. Login

**POST** `/users/login`

Autentica a un usuario verificado y devuelve tokens de acceso.

**Acceso**: Público

**Request Body**:
```json
{
  "email": "usuario@ejemplo.com",
  "password": "mipassword123"
}
```

**Response (200 OK)**:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "name": "Juan",
  "surname": "Pérez"
}
```

**Errores**:
- `400 Bad Request`: Formato de request inválido
- `403 Forbidden`: Credenciales incorrectas o email no verificado

---

### 5. Renovar Access Token

**POST** `/users/refresh-token`

Genera un nuevo access token utilizando un refresh token válido.

**Acceso**: Público

**Request Body**:
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Validaciones**:
- `refresh_token`: requerido

**Response (200 OK)**:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Errores**:
- `400 Bad Request`: Request inválido
- `401 Unauthorized`: Refresh token inválido o expirado

---

### 6. Obtener Usuario por ID

**GET** `/users/:id`

Obtiene la información de un usuario específico por su ID.

**Acceso**: Protegido (requiere autenticación)

**Headers Requeridos**:
```
Authorization: Bearer <access_token>
```

**Path Parameters**:
- `id`: ID del usuario (número entero)

**Response (200 OK)**:
```json
{
  "id": 1,
  "email": "usuario@ejemplo.com",
  "first_name": "Juan",
  "last_name": "Pérez",
  "is_admin": false,
  "is_verified": true
}
```

**Errores**:
- `400 Bad Request`: ID inválido
- `401 Unauthorized`: Token faltante o inválido
- `404 Not Found`: Usuario no encontrado

---

### 7. Verificar Token de Admin

**GET** `/users/admin`

Verifica que el token proporcionado pertenece a un usuario con rol de administrador.

**Acceso**: Solo administradores

**Headers Requeridos**:
```
Authorization: Bearer <access_token>
```

**Response (200 OK)**: Sin body (status 200 indica token válido)

**Errores**:
- `401 Unauthorized`: Token inválido o usuario no es administrador

---

### 8. Promover Usuario a Admin

**POST** `/users/promote-admin`

Otorga permisos de administrador a un usuario existente.

**Acceso**: Solo administradores

**Headers Requeridos**:
```
Authorization: Bearer <access_token>
```

**Request Body**:
```json
{
  "user_id": 5
}
```

**Validaciones**:
- `user_id`: requerido, número entero

**Response (200 OK)**:
```json
{
  "message": "User promoted to admin successfully"
}
```

**Errores**:
- `400 Bad Request`: Request inválido o usuario no encontrado
- `401 Unauthorized`: Token inválido o usuario no es administrador

---

## Modelos de Datos

### UserDto

Representa la información completa de un usuario en el sistema.

```typescript
{
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  is_admin: boolean;
  is_verified: boolean;
}
```

**Nota**: El campo `password` nunca se devuelve en las respuestas de la API.

---

## Autenticación

### Flujo de Autenticación

1. **Registro**: Usuario se registra → Recibe código por email
2. **Verificación**: Usuario ingresa código → Recibe access token y refresh token
3. **Acceso**: Usuario usa access token en header `Authorization: Bearer <token>`
4. **Renovación**: Cuando access token expira, usar refresh token para obtener uno nuevo

### Uso de Tokens

Para endpoints protegidos, incluir el access token en el header de cada request:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

### Tipos de Acceso

- **Público**: No requiere token
- **Protegido**: Requiere access token válido
- **Admin**: Requiere access token de usuario con rol administrador

---

## Códigos de Estado HTTP

| Código | Significado |
|--------|-------------|
| 200 | Operación exitosa |
| 201 | Recurso creado exitosamente |
| 400 | Request inválido (datos malformados o validación fallida) |
| 401 | No autenticado (token inválido o faltante) |
| 403 | No autorizado (credenciales incorrectas) |
| 404 | Recurso no encontrado |

---

## Flujos de Usuario Comunes

### Flujo de Registro Completo

```
1. POST /users/register
   → Usuario recibe email con código

2. POST /users/verify-email
   → Usuario obtiene tokens

3. Usuario puede comenzar a usar la aplicación
```

### Flujo de Login

```
1. POST /users/login
   → Usuario obtiene tokens

2. Usar access_token para requests protegidos

3. Cuando access_token expire:
   POST /users/refresh-token
   → Obtener nuevos tokens
```

### Flujo de Reenvío de Código

```
Si el usuario no recibió el código:

1. POST /users/resend-code
   → Nuevo código enviado

2. POST /users/verify-email
   → Verificar con el nuevo código
```

---

## Notas Importantes

- Los códigos de verificación tienen **6 dígitos** y **expiran** después de un tiempo determinado
- Los **access tokens** tienen corta duración (usar refresh token para renovar)
- Los **refresh tokens** tienen larga duración pero eventualmente expiran
- El email es **único** por usuario y se usa como identificador de login
- Las contraseñas deben tener **mínimo 6 caracteres**
- Los usuarios deben **verificar su email** antes de poder hacer login
- Solo los **administradores** pueden promover otros usuarios a admin

---

## Ejemplos de Integración

### Registro y Verificación (JavaScript/TypeScript)

```javascript
// 1. Registro
const registerResponse = await fetch('http://localhost:8080/users/register', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'password123',
    first_name: 'Juan',
    last_name: 'Pérez'
  })
});

// 2. Verificación
const verifyResponse = await fetch('http://localhost:8080/users/verify-email', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'user@example.com',
    code: '123456'
  })
});

const { access_token, refresh_token } = await verifyResponse.json();
// Guardar tokens en localStorage o estado de la aplicación
```

### Request Autenticado

```javascript
const getUserResponse = await fetch('http://localhost:8080/users/1', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${access_token}`
  }
});

const userData = await getUserResponse.json();
```

### Renovar Token

```javascript
const refreshResponse = await fetch('http://localhost:8080/users/refresh-token', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    refresh_token: refresh_token
  })
});

const { access_token: newAccessToken } = await refreshResponse.json();
// Actualizar el access token almacenado
```
