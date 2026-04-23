# 🐕 Dogs API Manager

A full-stack web application for managing dog breeds with CRUD operations, persistent storage, and a modern web interface.

## 🚀 Features

- **Full CRUD Operations**: Create, Read, Update, and Delete dog breeds
- **Persistent Storage**: SQLite database for data persistence
- **Search Functionality**: Search dogs by breed name or sub-breeds
- **Modern UI**: Clean, responsive web interface built with vanilla JavaScript
- **RESTful API**: Well-structured REST API endpoints
- **Real-time Stats**: Live statistics showing total breeds and sub-breeds

## 🛠️ Tech Stack

### Backend
- **Go 1.21**: High-performance backend
- **SQLite**: Embedded database for persistence
- **Gorilla Mux**: HTTP router and URL matcher
- **CORS**: Cross-origin resource sharing support

### Frontend
- **Vanilla JavaScript**: No frameworks, pure JS
- **CSS3**: Modern styling with gradients and animations
- **Responsive Design**: Works on all screen sizes

## 📋 Prerequisites

- Go 1.21 or higher
- Git

## 🔧 Installation & Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd dogs-api
```

### 2. Install Dependencies

```bash
cd backend
go mod download
```

### 3. Run the Application

```bash
go run main.go
```

The application will start on `http://localhost:8080`

## 🌐 API Endpoints

### Get All Dogs
```http
GET /api/dogs
```

### Get Dog by ID
```http
GET /api/dogs/{id}
```

### Get Dog by Breed
```http
GET /api/dogs/breed/{breed}
```

### Search Dogs
```http
GET /api/dogs/search?query={search_term}
```

### Create Dog
```http
POST /api/dogs
Content-Type: application/json

{
  "breed": "Golden Retriever",
  "subBreeds": "American, English"
}
```

### Update Dog
```http
PUT /api/dogs/{id}
Content-Type: application/json

{
  "breed": "Golden Retriever",
  "subBreeds": "American, English, Canadian"
}
```

### Delete Dog (Soft Delete)
```http
DELETE /api/dogs/{id}
```

## 💾 Database

The application uses SQLite for data persistence. The database file (`dogs.db`) is automatically created on first run and initialized with data from `dogs.json`.

### Database Schema

```sql
CREATE TABLE dogs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    breed TEXT NOT NULL UNIQUE,
    sub_breeds TEXT,
    active BOOLEAN DEFAULT 1
);
```

## 🎨 Web Interface

### Access the UI
Open your browser and navigate to: **http://localhost:8080**

### Features
- **Dashboard**: View all dog breeds in a card grid layout
- **Statistics**: Real-time count of total breeds and sub-breeds
- **Search**: Instant search by breed name or sub-breeds
- **Add New Dog**: Create new dog breeds with optional sub-breeds
- **Edit**: Update existing dog breeds
- **Delete**: Soft-delete dog breeds (data persists but is hidden)

## 🚀 Deployment

### Render (Recommended - Free Tier)

The application includes a `render.yaml` configuration file for easy deployment with persistent storage.

1. Push your code to GitHub
2. Create a new Web Service on Render
3. Connect your GitHub repository
4. Render will automatically detect the Docker configuration
5. Deploy!

**Note**: The application uses a persistent disk for SQLite database storage, ensuring data persists across deployments.

### Docker Deployment

```bash
# Build the Docker image
docker build -t dogs-api .

# Run the container
docker run -p 8080:8080 -v $(pwd)/data:/data dogs-api

# Or use the deployment script
./deploy.sh
```

### Manual Deployment

```bash
# Build the application
cd backend
go build -o dogs-api main.go

# Run the binary
./dogs-api
```

## 📁 Project Structure

```
dogs-api/
├── backend/
│   ├── main.go           # Go backend application
│   ├── dogs.json         # Initial dog data
│   └── go.mod            # Go module dependencies
├── frontend/
│   ├── index.html        # Main HTML page
│   ├── style.css         # Styles
│   └── app.js            # Frontend JavaScript
├── render.yaml           # Render deployment config
├── README.md             # This file
└── .gitignore            # Git ignore rules
```

## 🧪 Testing the API

### Using cURL

```bash
# Get all dogs
curl http://localhost:8080/api/dogs

# Create a new dog
curl -X POST http://localhost:8080/api/dogs \
  -H "Content-Type: application/json" \
  -d '{"breed": "Test Breed", "subBreeds": "Sub1, Sub2"}'

# Update a dog
curl -X PUT http://localhost:8080/api/dogs/1 \
  -H "Content-Type: application/json" \
  -d '{"breed": "Updated Breed", "subBreeds": "New Sub"}'

# Delete a dog
curl -X DELETE http://localhost:8080/api/dogs/1

# Search dogs
curl http://localhost:8080/api/dogs/search?query=retriever
```

### Using the Web Interface

1. Open `http://localhost:8080` in your browser
2. Use the search box to find specific breeds
3. Click "Add New Dog" to create a new breed
4. Use the Edit/Delete buttons on each card to manage dogs

## 🔒 Security

- Input validation on all endpoints
- SQL injection prevention (parameterized queries)
- CORS enabled for development
- Soft delete for data safety

## 📝 License

This project is for demonstration purposes only.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## 📧 Contact

For questions or feedback, please contact the developer.
