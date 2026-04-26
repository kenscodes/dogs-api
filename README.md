# 🐕 Dogs API Manager

A full-stack web application for managing dog breeds with CRUD operations, persistent storage, and a modern web interface.

## 🚀 Features

- **Full CRUD Operations**: Create, Read, Update, and Delete dog breeds
- **Persistent Storage**: SQLite database for data persistence
- **Minimalist UI**: Clean, simple list view built with vanilla JavaScript
- **RESTful API**: Well-structured REST API endpoints
- **Timestamps**: Track when dogs were created and last updated
- **Sorted by Recency**: Recently modified dogs appear first

## 🛠️ Tech Stack

### Backend
- **Go 1.21**: High-performance backend
- **SQLite**: Embedded database for persistence
- **Gorilla Mux**: HTTP router and URL matcher
- **CORS**: Cross-origin resource sharing support

### Frontend
- **Vanilla JavaScript**: No frameworks, pure JS
- **CSS3**: Minimalist, clean styling
- **Responsive Design**: Works on all screen sizes

## 📋 Prerequisites

- Go 1.21 or higher
- Git

## 🔧 Installation & Setup

### Step 1: Clone the Repository

Open your terminal and run:

```bash
git clone <repository-url>
cd dogs-api
```

### Step 2: Install Go Dependencies

Navigate to the backend directory and install the required Go packages:

```bash
cd backend
go mod download
```

### Step 3: Run the Application

Start the Go server:

```bash
go run main.go
```

You should see: `Server running on port 8080`

### Step 4: Access the Web Interface

Open your web browser and navigate to:

**http://localhost:8080**

You should see the Dogs API interface with a list of dog breeds.

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
    active BOOLEAN DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🎨 Web Interface

### Access the UI
Open your browser and navigate to: **http://localhost:8080**

### Features
- **List View**: View all dog breeds in a clean, scannable list
- **Add New Dog**: Create new dog breeds with optional sub-breeds
- **Edit**: Update existing dog breeds
- **Delete**: Soft-delete dog breeds (data persists but is hidden)
- **Timestamps**: See when each dog was last updated

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
```

### Using the Web Interface

#### How to Add a New Dog Breed

1. Click the **"+ Add"** button at the top of the page
2. A modal will appear - enter the breed name (e.g., "Golden Retriever")
3. Optionally add sub-breeds separated by commas (e.g., "American, English")
4. Click **"Save"** to add the dog
5. You'll see a success message and the new dog will appear in the list

#### How to Edit a Dog Breed

1. Find the dog breed you want to edit in the list
2. Click the **"✏️ Edit"** button next to it
3. Update the breed name or sub-breeds in the modal
4. Click **"Save"** to save your changes

#### How to Delete a Dog Breed

1. Find the dog breed you want to delete in the list
2. Click the **"🗑️ Delete"** button next to it
3. A confirmation modal will appear
4. Click **"Delete"** to confirm
5. The dog will be removed from the list (soft delete - data is preserved in database)

## ❓ Troubleshooting

### "command not found: go" error

If you see this error, Go is not installed on your system. Install Go from https://golang.org/dl/

### "port 8080 already in use" error

If port 8080 is already in use by another application, you can either:
- Stop the other application using port 8080
- Or change the port by setting the PORT environment variable:

```bash
PORT=3000 go run main.go
```

Then access the app at `http://localhost:3000`

### Page not loading / 404 error

Make sure you're running the server from the `backend` directory:

```bash
cd backend
go run main.go
```

### Database errors

The SQLite database (`dogs.db`) is automatically created on first run. If you encounter database errors, try deleting the `dogs.db` file and restarting the server - it will be recreated with fresh data.

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
