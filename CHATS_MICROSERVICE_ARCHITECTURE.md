# Arquitectura del Microservicio de Chats - UniChat

Sistema de gestión de chats educativos con IA para consultar documentos académicos.

---

## 📋 Tabla de Contenidos

1. [Overview del Sistema](#overview-del-sistema)
2. [Arquitectura General](#arquitectura-general)
3. [Stack Tecnológico](#stack-tecnológico)
4. [Modelos de Datos](#modelos-de-datos)
5. [API REST Endpoints](#api-rest-endpoints)
6. [Integración con Microservicios](#integración-con-microservicios)
7. [Message Queue (RabbitMQ)](#message-queue-rabbitmq)
8. [Búsqueda con Solr](#búsqueda-con-solr)
9. [Estructura del Proyecto](#estructura-del-proyecto)
10. [Flujos de Trabajo](#flujos-de-trabajo)
11. [Configuración y Deployment](#configuración-y-deployment)

---

## 🎯 Overview del Sistema

### Descripción General

UniChat es un sistema educativo con inteligencia artificial que permite a profesores subir documentos académicos y a estudiantes consultarlos mediante conversaciones naturales. El microservicio de chats es el componente central que gestiona:

- **Asignaturas**: Materias creadas por profesores
- **Chats individuales**: Conversaciones privadas de cada estudiante por asignatura
- **Mensajes**: Historial completo de preguntas y respuestas
- **Búsqueda**: Consulta en historial de conversaciones con Solr

### Casos de Uso Principales

#### Profesor (Admin):
1. Crear asignaturas
2. Ver estadísticas de uso
3. Gestionar asignaturas (editar, eliminar)

#### Estudiante:
1. Ver asignaturas disponibles
2. Crear chat en una asignatura
3. Enviar mensajes/preguntas
4. Recibir respuestas de la IA
5. Buscar en historial de conversaciones
6. Consultar chats anteriores

### Ecosistema de Microservicios

```
┌─────────────────┐      ┌──────────────────┐      ┌─────────────────┐
│   Microservicio │      │  Microservicio   │      │ Microservicio   │
│   de Usuarios   │◄────►│   de Chats       │◄────►│  de IA/Docs     │
│   (Auth/JWT)    │      │  (Este servicio) │      │  (Indexación)   │
└─────────────────┘      └──────────────────┘      └─────────────────┘
                                  │                          │
                                  │                          │
                                  ▼                          ▼
                         ┌─────────────┐          ┌──────────────┐
                         │   MongoDB   │          │   RabbitMQ   │
                         │  (Chats DB) │          │ (Async Msgs) │
                         └─────────────┘          └──────────────┘
                                  │
                                  ▼
                         ┌─────────────┐
                         │    Solr     │
                         │  (Search)   │
                         └─────────────┘
```

---

## 🏗️ Arquitectura General

### Diagrama de Componentes

```
┌─────────────────────────────────────────────────────────────┐
│                  Microservicio de Chats                     │
│                                                              │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐ │
│  │              │    │              │    │              │ │
│  │ Controllers  │───►│  Services    │───►│ Repositories │ │
│  │  (HTTP API)  │    │  (Business   │    │ (MongoDB)    │ │
│  │              │    │   Logic)     │    │              │ │
│  └──────────────┘    └──────────────┘    └──────────────┘ │
│         │                    │                    │         │
│         │                    │                    │         │
│  ┌──────▼────────┐    ┌─────▼──────┐    ┌───────▼──────┐ │
│  │               │    │            │    │              │ │
│  │  Middleware   │    │   Queue    │    │    Search    │ │
│  │  (Auth/CORS)  │    │ (RabbitMQ) │    │    (Solr)    │ │
│  │               │    │            │    │              │ │
│  └───────────────┘    └────────────┘    └──────────────┘ │
│                                                              │
└─────────────────────────────────────────────────────────────┘
         │                      │                      │
         │                      │                      │
         ▼                      ▼                      ▼
┌────────────────┐    ┌──────────────┐      ┌──────────────┐
│  Users Service │    │   AI Service │      │  Solr Server │
│   (JWT Auth)   │    │ (via Queue)  │      │   (Search)   │
└────────────────┘    └──────────────┘      └──────────────┘
```

### Capas de la Aplicación

1. **API Layer (Controllers)**
   - Manejo de requests HTTP
   - Validación de entrada
   - Serialización JSON

2. **Business Logic Layer (Services)**
   - Lógica de negocio
   - Orquestación de operaciones
   - Integración con otros servicios

3. **Data Access Layer (Repositories)**
   - Operaciones de base de datos
   - Queries a MongoDB
   - Manejo de transacciones

4. **Infraestructura Layer**
   - Middleware de autenticación
   - Cliente RabbitMQ
   - Cliente Solr
   - Configuración

---

## 🔧 Stack Tecnológico

### Backend
- **Go 1.24+**: Lenguaje principal
- **Gin Framework**: HTTP web framework
- **MongoDB Driver**: Cliente oficial de MongoDB para Go
- **RabbitMQ (AMQP)**: Message broker para comunicación asíncrona
- **Solr Client**: Cliente para Apache Solr

### Bases de Datos y Storage
- **MongoDB 7.0**: Base de datos principal (NoSQL)
- **Apache Solr 9.0**: Motor de búsqueda e indexación

### Infraestructura
- **Docker & Docker Compose**: Contenedores
- **RabbitMQ 3.12**: Message queue

### Dependencias Go Principales

```go
require (
    github.com/gin-gonic/gin v1.10.0
    github.com/gin-contrib/cors v1.7.5
    go.mongodb.org/mongo-driver v1.13.1
    github.com/streadway/amqp v1.1.0      // RabbitMQ
    github.com/rtt/Go-Solr v0.0.0         // Solr client
    github.com/golang-jwt/jwt/v5 v5.2.2   // JWT validation
    github.com/sirupsen/logrus v1.9.3     // Logging
)
```

---

## 📊 Modelos de Datos

### 1. Subject (Asignatura)

Representa una materia/asignatura creada por un profesor.

```go
type Subject struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name        string             `bson:"name" json:"name" binding:"required"`
    Description string             `bson:"description" json:"description"`
    ProfessorID int                `bson:"professor_id" json:"professor_id"` // ID del users service
    CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
    IsActive    bool               `bson:"is_active" json:"is_active"`

    // Metadata
    DocumentCount int              `bson:"document_count" json:"document_count"` // Referencia
    StudentCount  int              `bson:"student_count" json:"student_count"`   // Contador
}
```

**Índices MongoDB:**
```javascript
db.subjects.createIndex({ "professor_id": 1 })
db.subjects.createIndex({ "is_active": 1 })
db.subjects.createIndex({ "created_at": -1 })
```

---

### 2. Chat

Representa una conversación individual de un estudiante en una asignatura.

```go
type Chat struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    SubjectID primitive.ObjectID `bson:"subject_id" json:"subject_id" binding:"required"`
    UserID    int                `bson:"user_id" json:"user_id"`       // ID del users service
    Title     string             `bson:"title" json:"title"`           // Opcional: título del chat
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`

    // Metadata
    MessageCount int              `bson:"message_count" json:"message_count"`
    LastMessage  string           `bson:"last_message" json:"last_message"`
    LastActivity time.Time        `bson:"last_activity" json:"last_activity"`
}
```

**Índices MongoDB:**
```javascript
db.chats.createIndex({ "subject_id": 1, "user_id": 1 }, { unique: true })
db.chats.createIndex({ "user_id": 1, "updated_at": -1 })
db.chats.createIndex({ "subject_id": 1 })
```

**Nota:** El índice único garantiza que un estudiante solo tenga un chat por asignatura.

---

### 3. Message

Representa un mensaje individual en un chat (pregunta del usuario o respuesta de la IA).

```go
type Message struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    ChatID    primitive.ObjectID `bson:"chat_id" json:"chat_id" binding:"required"`
    Content   string             `bson:"content" json:"content" binding:"required"`
    Sender    string             `bson:"sender" json:"sender"`         // "user" o "ai"
    Timestamp time.Time          `bson:"timestamp" json:"timestamp"`

    // Metadata para respuestas de IA
    AIMetadata AIMetadata         `bson:"ai_metadata,omitempty" json:"ai_metadata,omitempty"`

    // Indexación
    Indexed   bool               `bson:"indexed" json:"indexed"`       // Si está en Solr
}

type AIMetadata struct {
    Model          string   `bson:"model" json:"model"`                     // Modelo de IA usado
    Confidence     float64  `bson:"confidence" json:"confidence"`           // Confianza de la respuesta
    SourceDocs     []string `bson:"source_docs" json:"source_docs"`         // Documentos fuente
    ProcessingTime int64    `bson:"processing_time_ms" json:"processing_time_ms"` // Tiempo en ms
}
```

**Índices MongoDB:**
```javascript
db.messages.createIndex({ "chat_id": 1, "timestamp": 1 })
db.messages.createIndex({ "timestamp": -1 })
db.messages.createIndex({ "indexed": 1 }) // Para sincronización con Solr
```

---

## 🌐 API REST Endpoints

### Base URL
```
http://localhost:8081/api/v1
```

---

### 📚 Asignaturas (Subjects)

#### 1. Crear Asignatura
```http
POST /subjects
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "name": "Introducción a la Programación",
  "description": "Curso básico de programación en Python"
}
```

**Response (201 Created):**
```json
{
  "id": "65a1b2c3d4e5f6g7h8i9j0k1",
  "name": "Introducción a la Programación",
  "description": "Curso básico de programación en Python",
  "professor_id": 5,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "is_active": true,
  "document_count": 0,
  "student_count": 0
}
```

**Permisos:** Solo administradores (profesores)

---

#### 2. Listar Asignaturas
```http
GET /subjects?page=1&limit=10&professor_id=5&active=true
Authorization: Bearer <token>
```

**Query Parameters:**
- `page`: Número de página (default: 1)
- `limit`: Resultados por página (default: 10, max: 50)
- `professor_id`: Filtrar por profesor (opcional)
- `active`: Filtrar por activas (true/false, opcional)

**Response (200 OK):**
```json
{
  "subjects": [
    {
      "id": "65a1b2c3d4e5f6g7h8i9j0k1",
      "name": "Introducción a la Programación",
      "description": "Curso básico de programación en Python",
      "professor_id": 5,
      "created_at": "2024-01-15T10:30:00Z",
      "is_active": true,
      "document_count": 15,
      "student_count": 45
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 23,
    "total_pages": 3
  }
}
```

---

#### 3. Obtener Asignatura por ID
```http
GET /subjects/:id
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "id": "65a1b2c3d4e5f6g7h8i9j0k1",
  "name": "Introducción a la Programación",
  "description": "Curso básico de programación en Python",
  "professor_id": 5,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "is_active": true,
  "document_count": 15,
  "student_count": 45
}
```

---

#### 4. Actualizar Asignatura
```http
PUT /subjects/:id
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "name": "Programación Avanzada en Python",
  "description": "Curso actualizado",
  "is_active": true
}
```

**Permisos:** Solo el profesor dueño o super admin

---

#### 5. Eliminar Asignatura
```http
DELETE /subjects/:id
Authorization: Bearer <admin_token>
```

**Response (200 OK):**
```json
{
  "message": "Subject deleted successfully"
}
```

**Permisos:** Solo el profesor dueño o super admin

---

### 💬 Chats

#### 1. Crear Chat en Asignatura
```http
POST /subjects/:subject_id/chats
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Consultas sobre variables" // Opcional
}
```

**Response (201 Created):**
```json
{
  "id": "65a1b2c3d4e5f6g7h8i9j0k2",
  "subject_id": "65a1b2c3d4e5f6g7h8i9j0k1",
  "user_id": 10,
  "title": "Consultas sobre variables",
  "created_at": "2024-01-16T14:20:00Z",
  "updated_at": "2024-01-16T14:20:00Z",
  "message_count": 0
}
```

**Nota:** Si el usuario ya tiene un chat en esta asignatura, retorna el existente (409 Conflict con el chat existente).

---

#### 2. Listar Mis Chats en Asignatura
```http
GET /subjects/:subject_id/chats
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "chats": [
    {
      "id": "65a1b2c3d4e5f6g7h8i9j0k2",
      "subject_id": "65a1b2c3d4e5f6g7h8i9j0k1",
      "user_id": 10,
      "title": "Consultas sobre variables",
      "created_at": "2024-01-16T14:20:00Z",
      "message_count": 5,
      "last_message": "¿Puedes explicarme las listas?",
      "last_activity": "2024-01-16T15:30:00Z"
    }
  ]
}
```

---

#### 3. Listar Todos Mis Chats (Todas las Asignaturas)
```http
GET /chats
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "chats": [
    {
      "id": "65a1b2c3d4e5f6g7h8i9j0k2",
      "subject_id": "65a1b2c3d4e5f6g7h8i9j0k1",
      "subject_name": "Introducción a la Programación",
      "user_id": 10,
      "title": "Consultas sobre variables",
      "message_count": 5,
      "last_activity": "2024-01-16T15:30:00Z"
    }
  ]
}
```

---

#### 4. Obtener Chat por ID
```http
GET /chats/:id
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "id": "65a1b2c3d4e5f6g7h8i9j0k2",
  "subject_id": "65a1b2c3d4e5f6g7h8i9j0k1",
  "subject_name": "Introducción a la Programación",
  "user_id": 10,
  "title": "Consultas sobre variables",
  "created_at": "2024-01-16T14:20:00Z",
  "message_count": 5,
  "last_activity": "2024-01-16T15:30:00Z"
}
```

**Permisos:** Solo el dueño del chat o admin

---

#### 5. Eliminar Chat
```http
DELETE /chats/:id
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "message": "Chat deleted successfully"
}
```

**Permisos:** Solo el dueño del chat o admin

---

### 📨 Mensajes (Messages)

#### 1. Enviar Mensaje
```http
POST /chats/:chat_id/messages
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "¿Qué es una variable en Python?"
}
```

**Response (201 Created):**
```json
{
  "message": {
    "id": "65a1b2c3d4e5f6g7h8i9j0k3",
    "chat_id": "65a1b2c3d4e5f6g7h8i9j0k2",
    "content": "¿Qué es una variable en Python?",
    "sender": "user",
    "timestamp": "2024-01-16T15:30:00Z"
  },
  "ai_response_pending": true
}
```

**Flujo:**
1. Se guarda el mensaje del usuario en MongoDB
2. Se publica en RabbitMQ para que la IA procese
3. La IA responde (asíncronamente) y se guarda la respuesta
4. El frontend puede hacer polling o usar WebSockets para obtener la respuesta

---

#### 2. Obtener Historial de Mensajes
```http
GET /chats/:chat_id/messages?page=1&limit=50
Authorization: Bearer <token>
```

**Query Parameters:**
- `page`: Número de página (default: 1)
- `limit`: Mensajes por página (default: 50, max: 100)
- `before`: Timestamp para paginación (obtener mensajes anteriores)

**Response (200 OK):**
```json
{
  "messages": [
    {
      "id": "65a1b2c3d4e5f6g7h8i9j0k3",
      "chat_id": "65a1b2c3d4e5f6g7h8i9j0k2",
      "content": "¿Qué es una variable en Python?",
      "sender": "user",
      "timestamp": "2024-01-16T15:30:00Z"
    },
    {
      "id": "65a1b2c3d4e5f6g7h8i9j0k4",
      "chat_id": "65a1b2c3d4e5f6g7h8i9j0k2",
      "content": "Una variable en Python es un contenedor...",
      "sender": "ai",
      "timestamp": "2024-01-16T15:30:05Z",
      "ai_metadata": {
        "model": "gpt-4",
        "confidence": 0.95,
        "source_docs": ["intro_python.pdf", "variables_chapter.pdf"],
        "processing_time_ms": 1250
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 10,
    "has_more": false
  }
}
```

---

#### 3. Buscar en Historial (Solr)
```http
GET /chats/:chat_id/messages/search?q=variable&page=1&limit=20
Authorization: Bearer <token>
```

**Query Parameters:**
- `q`: Texto a buscar (requerido)
- `page`: Número de página (default: 1)
- `limit`: Resultados por página (default: 20, max: 50)

**Response (200 OK):**
```json
{
  "results": [
    {
      "id": "65a1b2c3d4e5f6g7h8i9j0k3",
      "chat_id": "65a1b2c3d4e5f6g7h8i9j0k2",
      "content": "¿Qué es una <em>variable</em> en Python?",
      "sender": "user",
      "timestamp": "2024-01-16T15:30:00Z",
      "score": 0.95
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 3
  }
}
```

---

### 📊 Estadísticas (Admin)

#### 1. Estadísticas Generales
```http
GET /admin/stats
Authorization: Bearer <admin_token>
```

**Response (200 OK):**
```json
{
  "total_subjects": 45,
  "total_chats": 1250,
  "total_messages": 18500,
  "active_users_today": 78,
  "avg_messages_per_chat": 14.8
}
```

---

#### 2. Estadísticas por Asignatura
```http
GET /admin/subjects/:id/stats
Authorization: Bearer <admin_token>
```

**Response (200 OK):**
```json
{
  "subject_id": "65a1b2c3d4e5f6g7h8i9j0k1",
  "total_chats": 45,
  "total_messages": 680,
  "unique_students": 42,
  "avg_messages_per_student": 16.2,
  "most_active_hours": ["14:00", "15:00", "16:00"]
}
```

---

## 🔗 Integración con Microservicios

### 1. Microservicio de Usuarios (Autenticación)

#### Middleware de Autenticación

```go
// middleware/auth.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "Authorization token required"})
            c.Abort()
            return
        }

        // Validar JWT usando la misma lógica que el microservicio de usuarios
        claims, err := validateJWT(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // Guardar información del usuario en el contexto
        c.Set("user_id", claims.UserID)
        c.Set("is_admin", claims.IsAdmin)
        c.Next()
    }
}
```

#### Validación de Roles

```go
func AdminOnly() gin.HandlerFunc {
    return func(c *gin.Context) {
        isAdmin, exists := c.Get("is_admin")
        if !exists || !isAdmin.(bool) {
            c.JSON(403, gin.H{"error": "Admin access required"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

**Uso:**
```go
// Ruta protegida para usuarios autenticados
router.GET("/chats", AuthMiddleware(), chatController.GetMyChats)

// Ruta protegida solo para admins
router.POST("/subjects", AuthMiddleware(), AdminOnly(), subjectController.Create)
```

---

### 2. Microservicio de IA/Documentos (RabbitMQ)

#### Configuración de RabbitMQ

**Exchanges:**
- `unichat.questions` (direct): Para enviar preguntas a la IA
- `unichat.answers` (direct): Para recibir respuestas de la IA

**Queues:**
- `questions.queue`: Cola de preguntas pendientes
- `answers.queue`: Cola de respuestas de la IA

#### Publicar Pregunta (Producer)

```go
// queue/producer.go
type QuestionMessage struct {
    MessageID string    `json:"message_id"`
    ChatID    string    `json:"chat_id"`
    SubjectID string    `json:"subject_id"`
    Question  string    `json:"question"`
    UserID    int       `json:"user_id"`
    Timestamp time.Time `json:"timestamp"`
}

func PublishQuestion(msg QuestionMessage) error {
    body, err := json.Marshal(msg)
    if err != nil {
        return err
    }

    err = rabbitChannel.Publish(
        "unichat.questions", // exchange
        "question.route",    // routing key
        false,               // mandatory
        false,               // immediate
        amqp.Publishing{
            ContentType:  "application/json",
            Body:         body,
            DeliveryMode: amqp.Persistent,
        },
    )
    return err
}
```

#### Consumir Respuestas (Consumer)

```go
// queue/consumer.go
type AnswerMessage struct {
    MessageID      string   `json:"message_id"`
    ChatID         string   `json:"chat_id"`
    Answer         string   `json:"answer"`
    Model          string   `json:"model"`
    Confidence     float64  `json:"confidence"`
    SourceDocs     []string `json:"source_docs"`
    ProcessingTime int64    `json:"processing_time_ms"`
}

func ConsumeAnswers() {
    msgs, err := rabbitChannel.Consume(
        "answers.queue", // queue
        "",              // consumer
        false,           // auto-ack
        false,           // exclusive
        false,           // no-local
        false,           // no-wait
        nil,             // args
    )

    for msg := range msgs {
        var answer AnswerMessage
        json.Unmarshal(msg.Body, &answer)

        // Guardar respuesta en MongoDB
        saveAIResponse(answer)

        // Indexar en Solr
        indexMessageInSolr(answer)

        msg.Ack(false)
    }
}
```

---

### 3. Apache Solr (Búsqueda)

#### Schema de Solr

```xml
<!-- solr/messages_schema.xml -->
<schema name="messages" version="1.6">
  <field name="id" type="string" indexed="true" stored="true" required="true" />
  <field name="chat_id" type="string" indexed="true" stored="true" required="true" />
  <field name="subject_id" type="string" indexed="true" stored="true" />
  <field name="user_id" type="int" indexed="true" stored="true" />
  <field name="content" type="text_general" indexed="true" stored="true" required="true" />
  <field name="sender" type="string" indexed="true" stored="true" />
  <field name="timestamp" type="pdate" indexed="true" stored="true" />
  <field name="_text_" type="text_general" indexed="true" stored="false" multiValued="true"/>

  <copyField source="content" dest="_text_"/>

  <uniqueKey>id</uniqueKey>
</schema>
```

#### Cliente Solr en Go

```go
// search/solr_client.go
type SolrClient struct {
    baseURL string
    client  *http.Client
}

func (s *SolrClient) IndexMessage(msg Message) error {
    doc := map[string]interface{}{
        "id":         msg.ID.Hex(),
        "chat_id":    msg.ChatID.Hex(),
        "subject_id": msg.SubjectID,
        "user_id":    msg.UserID,
        "content":    msg.Content,
        "sender":     msg.Sender,
        "timestamp":  msg.Timestamp,
    }

    url := fmt.Sprintf("%s/solr/messages/update?commit=true", s.baseURL)
    // POST document to Solr
    // ...
}

func (s *SolrClient) SearchMessages(chatID string, query string, page int, limit int) ([]Message, error) {
    start := (page - 1) * limit

    url := fmt.Sprintf(
        "%s/solr/messages/select?q=content:%s AND chat_id:%s&start=%d&rows=%d&hl=true&hl.fl=content",
        s.baseURL, query, chatID, start, limit,
    )

    // GET from Solr and parse results
    // ...
}
```

---

## 📮 Message Queue (RabbitMQ)

### Arquitectura de Mensajería

```
┌──────────────┐                    ┌──────────────┐
│   Chats MS   │                    │   AI/Docs MS │
│              │                    │              │
│  ┌────────┐  │   Question Queue   │  ┌────────┐  │
│  │Producer├──┼───────────────────►│  │Consumer│  │
│  └────────┘  │                    │  └────┬───┘  │
│              │                    │       │      │
│  ┌────────┐  │   Answer Queue     │  ┌────▼───┐  │
│  │Consumer│◄─┼────────────────────┼──┤Producer│  │
│  └────────┘  │                    │  └────────┘  │
└──────────────┘                    └──────────────┘
```

### Configuración

```go
// config/rabbitmq.go
type RabbitMQConfig struct {
    URL              string
    QuestionExchange string
    AnswerExchange   string
    QuestionQueue    string
    AnswerQueue      string
}

func InitRabbitMQ(config RabbitMQConfig) (*amqp.Connection, *amqp.Channel, error) {
    conn, err := amqp.Dial(config.URL)
    if err != nil {
        return nil, nil, err
    }

    ch, err := conn.Channel()
    if err != nil {
        return nil, nil, err
    }

    // Declarar exchanges
    err = ch.ExchangeDeclare(
        config.QuestionExchange, // name
        "direct",                // type
        true,                    // durable
        false,                   // auto-deleted
        false,                   // internal
        false,                   // no-wait
        nil,                     // arguments
    )

    // Declarar queues
    _, err = ch.QueueDeclare(
        config.QuestionQueue, // name
        true,                 // durable
        false,                // delete when unused
        false,                // exclusive
        false,                // no-wait
        nil,                  // arguments
    )

    // Bind queue to exchange
    err = ch.QueueBind(
        config.QuestionQueue,    // queue name
        "question.route",        // routing key
        config.QuestionExchange, // exchange
        false,
        nil,
    )

    return conn, ch, nil
}
```

### Manejo de Errores y Reintentos

```go
// queue/retry.go
type RetryConfig struct {
    MaxRetries int
    RetryDelay time.Duration
}

func PublishWithRetry(msg QuestionMessage, config RetryConfig) error {
    var err error
    for i := 0; i < config.MaxRetries; i++ {
        err = PublishQuestion(msg)
        if err == nil {
            return nil
        }

        log.Printf("Retry %d/%d for message %s", i+1, config.MaxRetries, msg.MessageID)
        time.Sleep(config.RetryDelay)
    }

    // Si falla después de todos los reintentos, guardar en dead letter queue
    return fmt.Errorf("failed after %d retries: %w", config.MaxRetries, err)
}
```

---

## 🔍 Búsqueda con Solr

### Flujo de Indexación

```
┌─────────────┐
│   Mensaje   │
│  guardado   │
│  en MongoDB │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Worker    │
│  Background │  (Cada 5 segundos)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Buscar     │
│  mensajes   │
│  no         │
│  indexados  │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Indexar    │
│  en Solr    │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Actualizar │
│  flag       │
│  "indexed"  │
└─────────────┘
```

### Worker de Indexación

```go
// search/indexer.go
type Indexer struct {
    solrClient *SolrClient
    messageRepo *repositories.MessageRepository
}

func (idx *Indexer) StartIndexWorker(interval time.Duration) {
    ticker := time.NewTicker(interval)

    for range ticker.C {
        // Buscar mensajes no indexados
        messages, err := idx.messageRepo.FindUnindexed(100)
        if err != nil {
            log.Error("Error finding unindexed messages:", err)
            continue
        }

        // Indexar en lote
        for _, msg := range messages {
            err := idx.solrClient.IndexMessage(msg)
            if err != nil {
                log.Error("Error indexing message:", err)
                continue
            }

            // Marcar como indexado
            idx.messageRepo.MarkAsIndexed(msg.ID)
        }

        log.Infof("Indexed %d messages", len(messages))
    }
}
```

### Búsqueda con Highlighting

```go
// search/search_service.go
type SearchResult struct {
    Message     Message
    Highlights  []string
    Score       float64
}

func SearchInChat(chatID string, query string) ([]SearchResult, error) {
    results, err := solrClient.SearchMessages(chatID, query, 1, 20)
    if err != nil {
        return nil, err
    }

    // Procesar highlights de Solr
    searchResults := make([]SearchResult, len(results))
    for i, msg := range results {
        searchResults[i] = SearchResult{
            Message:    msg,
            Highlights: extractHighlights(msg),
            Score:      msg.Score,
        }
    }

    return searchResults, nil
}
```

---

## 📁 Estructura del Proyecto

```
chats-microservice/
├── main.go                     # Entry point
├── go.mod                      # Go dependencies
├── go.sum
├── Dockerfile                  # Container definition
├── docker-compose.yml          # Multi-container setup
│
├── config/                     # Configuration
│   ├── config.go              # Load env vars
│   ├── mongodb.go             # MongoDB connection
│   ├── rabbitmq.go            # RabbitMQ setup
│   └── solr.go                # Solr client setup
│
├── models/                     # Data models
│   ├── subject.go             # Subject model
│   ├── chat.go                # Chat model
│   ├── message.go             # Message model
│   └── pagination.go          # Pagination helpers
│
├── dto/                        # Data Transfer Objects
│   ├── subject_dto.go         # Subject request/response
│   ├── chat_dto.go            # Chat request/response
│   ├── message_dto.go         # Message request/response
│   └── stats_dto.go           # Statistics response
│
├── repositories/               # Data access layer
│   ├── subject_repository.go  # Subject DB operations
│   ├── chat_repository.go     # Chat DB operations
│   └── message_repository.go  # Message DB operations
│
├── services/                   # Business logic
│   ├── subject_service.go     # Subject business logic
│   ├── chat_service.go        # Chat business logic
│   ├── message_service.go     # Message business logic
│   └── stats_service.go       # Statistics logic
│
├── controllers/                # HTTP handlers
│   ├── subject_controller.go  # Subject endpoints
│   ├── chat_controller.go     # Chat endpoints
│   ├── message_controller.go  # Message endpoints
│   └── admin_controller.go    # Admin/stats endpoints
│
├── middleware/                 # HTTP middleware
│   ├── auth.go                # JWT authentication
│   ├── admin.go               # Admin role check
│   ├── cors.go                # CORS configuration
│   └── logger.go              # Request logging
│
├── queue/                      # Message queue
│   ├── producer.go            # Publish questions
│   ├── consumer.go            # Consume answers
│   └── models.go              # Queue message structs
│
├── search/                     # Solr integration
│   ├── solr_client.go         # Solr HTTP client
│   ├── indexer.go             # Background indexer
│   └── search_service.go      # Search operations
│
├── utils/                      # Utilities
│   ├── jwt.go                 # JWT validation
│   ├── errors.go              # Error handling
│   └── logger.go              # Logging setup
│
└── tests/                      # Tests
    ├── integration/
    ├── unit/
    └── mocks/
```

---

## 🔄 Flujos de Trabajo

### Flujo 1: Crear Asignatura (Profesor)

```
┌────────┐
│Profesor│
└───┬────┘
    │
    │ POST /subjects
    │ {name, description}
    ▼
┌──────────────┐
│  Controller  │
└──────┬───────┘
       │ Validate JWT (admin)
       ▼
┌──────────────┐
│   Service    │
│              │
│ 1. Validate  │
│ 2. Create    │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Repository  │
│              │
│ Insert into  │
│   MongoDB    │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Response   │
│  201 Created │
└──────────────┘
```

---

### Flujo 2: Enviar Mensaje con IA

```
┌─────────┐
│Estudiante│
└────┬────┘
     │
     │ POST /chats/:id/messages
     │ {content: "¿Qué es Python?"}
     ▼
┌────────────┐
│ Controller │
└─────┬──────┘
      │ Validate JWT
      │ Check chat ownership
      ▼
┌────────────┐
│  Service   │
│            │
│ 1. Save    │
│    user    │
│    message │
└─────┬──────┘
      │
      ├─────────────────────┐
      │                     │
      ▼                     ▼
┌────────────┐      ┌──────────────┐
│  MongoDB   │      │  RabbitMQ    │
│            │      │              │
│ Save msg   │      │ Publish      │
│ sender:    │      │ question     │
│ "user"     │      └──────┬───────┘
└─────┬──────┘             │
      │                    │
      │                    ▼
      │            ┌──────────────┐
      │            │  AI Service  │
      │            │              │
      │            │ Process      │
      │            │ question     │
      │            └──────┬───────┘
      │                   │
      │                   │ Publish answer
      │                   ▼
      │            ┌──────────────┐
      │            │  RabbitMQ    │
      │            │              │
      │            │ Answer queue │
      │            └──────┬───────┘
      │                   │
      ▼                   ▼
┌──────────────────────────┐
│    Answer Consumer       │
│                          │
│ 1. Receive AI response   │
│ 2. Save to MongoDB       │
│    sender: "ai"          │
│ 3. Index in Solr         │
└──────────────────────────┘
      │
      ▼
┌──────────────┐
│   Response   │
│  201 Created │
│  + pending   │
│    flag      │
└──────────────┘
```

---

### Flujo 3: Buscar en Historial

```
┌─────────┐
│Estudiante│
└────┬────┘
     │
     │ GET /chats/:id/messages/search?q=variable
     ▼
┌────────────┐
│ Controller │
└─────┬──────┘
      │ Validate JWT
      │ Check chat ownership
      ▼
┌────────────┐
│  Service   │
└─────┬──────┘
      │
      ▼
┌────────────┐
│Solr Client │
│            │
│ Query with │
│ highlighting│
└─────┬──────┘
      │
      ▼
┌────────────┐
│   Solr     │
│            │
│ Full-text  │
│   search   │
└─────┬──────┘
      │
      ▼
┌────────────┐
│  Process   │
│ highlights │
│  & scores  │
└─────┬──────┘
      │
      ▼
┌────────────┐
│  Response  │
│  200 OK    │
│  + results │
└────────────┘
```

---

## 🐳 Configuración y Deployment

### docker-compose.yml

```yaml
version: '3.9'

services:
  # MongoDB - Base de datos principal
  mongodb:
    image: mongo:7.0
    container_name: unichat-mongodb
    restart: unless-stopped
    environment:
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: adminpass
      MONGO_INITDB_DATABASE: unichat_chats
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db
      - ./mongo-init:/docker-entrypoint-initdb.d
    healthcheck:
      test: echo 'db.runCommand("ping").ok' | mongosh localhost:27017/test --quiet
      interval: 10s
      timeout: 5s
      retries: 5

  # RabbitMQ - Message broker
  rabbitmq:
    image: rabbitmq:3.12-management
    container_name: unichat-rabbitmq
    restart: unless-stopped
    environment:
      RABBITMQ_DEFAULT_USER: rabbitmq
      RABBITMQ_DEFAULT_PASS: rabbitmq
    ports:
      - "5672:5672"   # AMQP
      - "15672:15672" # Management UI
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
    healthcheck:
      test: rabbitmq-diagnostics -q ping
      interval: 10s
      timeout: 5s
      retries: 5

  # Apache Solr - Search engine
  solr:
    image: solr:9.0
    container_name: unichat-solr
    restart: unless-stopped
    ports:
      - "8983:8983"
    volumes:
      - solr_data:/var/solr
      - ./solr/messages_schema.xml:/opt/solr/messages_schema.xml
    command:
      - solr-precreate
      - messages
      - /opt/solr/messages_schema.xml
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8983/solr/messages/admin/ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Chats Microservice
  chats-backend:
    build:
      context: ./chats-microservice
      dockerfile: Dockerfile
    container_name: unichat-chats-backend
    restart: unless-stopped
    environment:
      # Server
      PORT: "8081"

      # MongoDB
      MONGODB_URI: "mongodb://admin:adminpass@mongodb:27017"
      MONGODB_DATABASE: "unichat_chats"

      # RabbitMQ
      RABBITMQ_URL: "amqp://rabbitmq:rabbitmq@rabbitmq:5672/"
      QUESTION_EXCHANGE: "unichat.questions"
      ANSWER_EXCHANGE: "unichat.answers"
      QUESTION_QUEUE: "questions.queue"
      ANSWER_QUEUE: "answers.queue"

      # Solr
      SOLR_URL: "http://solr:8983"
      SOLR_COLLECTION: "messages"

      # JWT (compartido con Users microservice)
      JWT_SECRET: "jwtSecret"
    ports:
      - "8081:8081"
    depends_on:
      mongodb:
        condition: service_healthy
      rabbitmq:
        condition: service_healthy
      solr:
        condition: service_healthy

volumes:
  mongodb_data:
  rabbitmq_data:
  solr_data:
```

---

### Dockerfile

```dockerfile
# Multi-stage build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o chats-service main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/chats-service .

EXPOSE 8081

CMD ["./chats-service"]
```

---

### Variables de Entorno (.env)

```bash
# Server Configuration
PORT=8081
GIN_MODE=release

# MongoDB
MONGODB_URI=mongodb://admin:adminpass@localhost:27017
MONGODB_DATABASE=unichat_chats
MONGODB_TIMEOUT=10s

# RabbitMQ
RABBITMQ_URL=amqp://rabbitmq:rabbitmq@localhost:5672/
QUESTION_EXCHANGE=unichat.questions
ANSWER_EXCHANGE=unichat.answers
QUESTION_QUEUE=questions.queue
ANSWER_QUEUE=answers.queue

# Solr
SOLR_URL=http://localhost:8983
SOLR_COLLECTION=messages
SOLR_BATCH_SIZE=100
SOLR_INDEX_INTERVAL=5s

# JWT (debe ser el mismo que Users microservice)
JWT_SECRET=jwtSecret

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

---

### main.go (Entry Point)

```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "chats-service/config"
    "chats-service/controllers"
    "chats-service/middleware"
    "chats-service/queue"
    "chats-service/repositories"
    "chats-service/search"
    "chats-service/services"

    "github.com/gin-gonic/gin"
)

func main() {
    // Load configuration
    cfg := config.LoadConfig()

    // Initialize MongoDB
    mongoClient, err := config.InitMongoDB(cfg.MongoDB)
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }
    defer mongoClient.Disconnect(context.Background())

    // Initialize RabbitMQ
    rabbitConn, rabbitCh, err := config.InitRabbitMQ(cfg.RabbitMQ)
    if err != nil {
        log.Fatal("Failed to connect to RabbitMQ:", err)
    }
    defer rabbitConn.Close()
    defer rabbitCh.Close()

    // Initialize Solr
    solrClient := search.NewSolrClient(cfg.Solr.URL, cfg.Solr.Collection)

    // Initialize repositories
    db := mongoClient.Database(cfg.MongoDB.Database)
    subjectRepo := repositories.NewSubjectRepository(db)
    chatRepo := repositories.NewChatRepository(db)
    messageRepo := repositories.NewMessageRepository(db)

    // Initialize services
    subjectService := services.NewSubjectService(subjectRepo)
    chatService := services.NewChatService(chatRepo, subjectRepo)
    messageService := services.NewMessageService(messageRepo, chatRepo, rabbitCh, solrClient)

    // Initialize controllers
    subjectCtrl := controllers.NewSubjectController(subjectService)
    chatCtrl := controllers.NewChatController(chatService)
    messageCtrl := controllers.NewMessageController(messageService)

    // Start background workers
    go queue.ConsumeAnswers(rabbitCh, messageService)

    indexer := search.NewIndexer(solrClient, messageRepo)
    go indexer.StartIndexWorker(cfg.Solr.IndexInterval)

    // Setup router
    router := setupRouter(subjectCtrl, chatCtrl, messageCtrl)

    // Start server
    srv := &http.Server{
        Addr:    ":" + cfg.Port,
        Handler: router,
    }

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal("Server error:", err)
        }
    }()

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Println("Server exited")
}

func setupRouter(
    subjectCtrl *controllers.SubjectController,
    chatCtrl *controllers.ChatController,
    messageCtrl *controllers.MessageController,
) *gin.Engine {
    router := gin.Default()

    router.Use(middleware.CORS())
    router.Use(middleware.Logger())

    api := router.Group("/api/v1")
    {
        // Public endpoints
        api.GET("/health", func(c *gin.Context) {
            c.JSON(200, gin.H{"status": "ok"})
        })

        // Protected endpoints
        auth := api.Group("")
        auth.Use(middleware.AuthMiddleware())
        {
            // Subjects
            subjects := auth.Group("/subjects")
            {
                subjects.GET("", subjectCtrl.List)
                subjects.GET("/:id", subjectCtrl.GetByID)
                subjects.POST("", middleware.AdminOnly(), subjectCtrl.Create)
                subjects.PUT("/:id", middleware.AdminOnly(), subjectCtrl.Update)
                subjects.DELETE("/:id", middleware.AdminOnly(), subjectCtrl.Delete)

                // Chats in subject
                subjects.POST("/:subject_id/chats", chatCtrl.CreateInSubject)
                subjects.GET("/:subject_id/chats", chatCtrl.ListInSubject)
            }

            // Chats
            chats := auth.Group("/chats")
            {
                chats.GET("", chatCtrl.ListMine)
                chats.GET("/:id", chatCtrl.GetByID)
                chats.DELETE("/:id", chatCtrl.Delete)

                // Messages in chat
                chats.POST("/:chat_id/messages", messageCtrl.Send)
                chats.GET("/:chat_id/messages", messageCtrl.List)
                chats.GET("/:chat_id/messages/search", messageCtrl.Search)
            }

            // Admin stats
            admin := auth.Group("/admin")
            admin.Use(middleware.AdminOnly())
            {
                admin.GET("/stats", adminCtrl.GetStats)
                admin.GET("/subjects/:id/stats", adminCtrl.GetSubjectStats)
            }
        }
    }

    return router
}
```

---

## 📊 Métricas y Monitoring (Opcional)

### Prometheus Metrics

```go
// middleware/metrics.go
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    httpDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration",
        },
        []string{"method", "endpoint"},
    )

    messagesProcessed = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "messages_processed_total",
            Help: "Total messages processed",
        },
    )
)

func PrometheusMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        c.Next()

        duration := time.Since(start).Seconds()

        httpRequestsTotal.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
            strconv.Itoa(c.Writer.Status()),
        ).Inc()

        httpDuration.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
        ).Observe(duration)
    }
}
```

---

## 🔐 Seguridad

### Best Practices Implementadas

1. **Autenticación JWT**: Tokens del microservicio de usuarios
2. **Autorización basada en roles**: Admin vs Usuario regular
3. **Validación de ownership**: Los usuarios solo acceden a sus propios chats
4. **CORS configurado**: Solo orígenes permitidos
5. **Rate limiting** (recomendado): Limitar requests por usuario
6. **Input validation**: Validación con Gin binding
7. **Secrets en env vars**: No hardcodear credenciales
8. **HTTPS en producción**: Usar reverse proxy (nginx)

---

## 📝 Próximos Pasos de Implementación

1. **Fase 1: Setup inicial**
   - Crear estructura de carpetas
   - Configurar MongoDB, RabbitMQ, Solr
   - Implementar modelos base

2. **Fase 2: CRUD básico**
   - Implementar Subjects CRUD
   - Implementar Chats CRUD
   - Implementar Messages CRUD

3. **Fase 3: Integración RabbitMQ**
   - Publisher de preguntas
   - Consumer de respuestas
   - Manejo de errores

4. **Fase 4: Búsqueda con Solr**
   - Indexer background worker
   - Search endpoint
   - Highlighting

5. **Fase 5: Testing y deployment**
   - Tests unitarios
   - Tests de integración
   - Docker compose completo

---

## 📚 Referencias

- **Go Gin Framework**: https://gin-gonic.com/docs/
- **MongoDB Go Driver**: https://www.mongodb.com/docs/drivers/go/current/
- **RabbitMQ Go Client**: https://github.com/streadway/amqp
- **Apache Solr**: https://solr.apache.org/guide/
- **JWT Go**: https://github.com/golang-jwt/jwt

---

**Versión:** 1.0
**Fecha:** 2024
**Autor:** Arquitectura UniChat
