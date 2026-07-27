# Zawya Reservation API

A backend system for a movie reservation platform, built with Go, Gin, and GORM.

## technologies
- Go (Golang)
- Gin framework
- GORM for database
- MySQL
- JWT Authentication

## setup
1. Install Go and MySQL
2. Install dependencies:`go mod tidy`
3. Set up the database and update `.env` with your DB credentials 
4. Run migrations: `go run cmd/api/main.go`
5. Start the server: `go run cmd/api/main.go`

## API Endpoints
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login and get JWT
- `GET /api/profile` - Get user profile (requires JWT)
- `GET /api/movies` - List all movies
- `GET /api/movies/:id` - Get movie by ID
- `GET /api/movies/:id/showtimes` - Get showtimes for a movie
- `GET /api/halls` - List all halls
- `GET /api/halls/:id` - Get hall with seats
- `GET /api/showtimes/:id` - Get showtime details
- `GET /api/showtimes/:id/seats` - Get seat availability for a showtime
- `POST /api/admin/movies` - Add new movie (admin only)
- `PUT /api/admin/movies/:id` - Update a movie (admin only)
- `DELETE /api/admin/movies/:id` - Delete a movie (admin only)
- `POST /api/admin/movies/:id/showtimes` - Create showtime for a movie (admin only)
- `POST /api/admin/halls` - Create hall (admin only)
- `PUT /api/admin/halls/:id` - Update a hall (admin only)
- `DELETE /api/admin/halls/:id` - Delete a hall (admin only)
- `POST /api/admin/halls/:id/seats` - Create seats for hall (admin only)
- `DELETE /api/admin/showtimes/:id` - Delete a showtime (admin only)
