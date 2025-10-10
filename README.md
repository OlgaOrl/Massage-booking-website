# Massage Booking System

A modern web application for booking massage appointments online. This system allows customers to browse available massage services, view pricing and duration, and select appointment times from an interactive calendar.

## Features

### Story #1 ✅ - Browse and Select Appointments
- **Service Selection**: Browse available massage services with detailed information
- **Service Details**: View duration (in minutes) and pricing (in euros) for each service
- **Interactive Calendar**: Navigate through the next 30 days to find available appointments
- **Visual Availability**: Clear visual distinction between available and booked time slots
- **Smart Scheduling**: Time slots generated based on service duration (60min=hourly, 90min=1.5h intervals)
- **Responsive Design**: Mobile-friendly interface that works on all devices
- **Real-time Data**: Live availability checking with SQLite database

### Story #2 ✅ - Complete Booking with Contact Information
- **Temporary Reservations**: 10-minute slot reservation system to prevent double-booking
- **Contact Form**: Collect customer name, email, and phone number
- **Real-time Validation**: Instant field validation with clear error messages
- **Countdown Timer**: Visual timer showing reservation expiration
- **Secure Booking**: Transaction-based booking confirmation
- **Auto-cleanup**: Background job removes expired reservations
- **Form Validation**: Both frontend and backend validation for data integrity

## Tech Stack

- **Backend**: Go (Golang) with standard library `net/http`
- **Frontend**: Vanilla HTML, CSS, JavaScript (no frameworks)
- **Database**: SQLite for data persistence
- **Deployment**: Docker with docker-compose
- **Architecture**: RESTful API with JSON responses

## Quick Start

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Installation

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd massage-booking-website
   ```

2. **Start the application**:
   ```bash
   docker compose up
   ```

3. **Open your browser**:
   ```
   http://localhost:8080
   ```

The application will automatically:
- Initialize the SQLite database
- Create all necessary tables
- Seed sample data (4 massage types, 30 days of time slots)
- Start the web server on port 8080

## API Documentation

### GET /api/massage-types

Returns a list of all available massage services.

**Response**:
```json
[
  {
    "id": 1,
    "name": "Swedish Massage",
    "duration": 60,
    "price": 50.0
  },
  {
    "id": 2,
    "name": "Deep Tissue",
    "duration": 90,
    "price": 70.0
  }
]
```

### GET /api/slots

Returns time slots for a specific date and service.

**Parameters**:
- `date` (required): Date in YYYY-MM-DD format
- `service_id` (required): ID of the massage service

**Example**:
```
GET /api/slots?date=2025-10-15&service_id=1
```

**Response**:
```json
[
  {
    "id": 1,
    "date": "2025-10-15",
    "time": "09:00",
    "service_id": 1,
    "available": true
  },
  {
    "id": 2,
    "date": "2025-10-15",
    "time": "10:00",
    "service_id": 1,
    "available": true
  }
]
```

**Note**: Time slots are generated based on service duration:
- 45-minute services: slots every 45 minutes (09:00, 09:45, 10:30...)
- 60-minute services: slots every hour (09:00, 10:00, 11:00...)
- 90-minute services: slots every 1.5 hours (09:00, 10:30, 12:00...)

### POST /api/reservations

Creates a temporary 10-minute reservation for a time slot.

**Request Body**:
```json
{
  "slot_id": 123
}
```

**Response**:
```json
{
  "reservation_id": 456,
  "expires_at": "2025-10-15T10:10:00Z",
  "expires_in_seconds": 600
}
```

### DELETE /api/reservations/:id

Cancels a temporary reservation.

**Response**: 204 No Content

### POST /api/bookings

Creates a confirmed booking with customer contact information.

**Request Body**:
```json
{
  "reservation_id": 456,
  "client_name": "John Doe",
  "email": "john@example.com",
  "phone": "+372 5123 4567",
  "service_id": 1,
  "date": "2025-10-15",
  "time_slot": "10:00"
}
```

**Response**:
```json
{
  "id": 789,
  "client_name": "John Doe",
  "email": "john@example.com",
  "phone": "+372 5123 4567",
  "service_id": 1,
  "date": "2025-10-15",
  "time_slot": "10:00",
  "created_at": "2025-10-15T09:55:30Z"
}
```

**Validation Rules**:
- **Name**: Required, minimum 2 characters, letters and spaces only
- **Email**: Required, valid email format
- **Phone**: Required, valid phone number format

## User Interface

### Service Selection
- Grid layout displaying all available massage services
- Each service card shows name, duration, and price
- Click any service to proceed to calendar selection

### Calendar View
- Displays the next 30 days from today
- Navigate between months using arrow buttons
- Past dates are disabled and cannot be selected
- Click any future date to view available time slots

### Time Slot Selection
- Time slots from 09:00 to 18:00 in 15-minute intervals
- **Available slots**: Green background (#4CAF50), clickable
- **Booked slots**: Grey background (#9E9E9E), strikethrough text, not clickable
- **Selected slot**: Highlighted with border or darker shade
- **Tooltip**: Hover over booked slots shows "This time is unavailable"

## Development

### Project Structure

```
massage-booking/
├── backend/
│   ├── main.go                  # Main entry point, HTTP server setup
│   ├── handlers/
│   │   ├── massage_types.go     # GET /api/massage-types handler
│   │   └── slots.go             # GET /api/slots handler
│   ├── models/
│   │   ├── massage_type.go      # MassageType struct definition
│   │   └── slot.go              # TimeSlot struct definition
│   ├── database/
│   │   └── db.go                # SQLite connection, migrations, queries
│   └── static/
│       ├── index.html           # Main page HTML
│       ├── style.css            # CSS styles
│       └── app.js               # Frontend JavaScript logic
├── Dockerfile                   # Multi-stage Docker build
├── docker-compose.yml           # Container orchestration
├── go.mod                       # Go module dependencies
├── .gitignore                   # Git ignore rules
└── README.md                    # This documentation
```

### Running Locally (without Docker)

If you prefer to run the application locally without Docker:

1. **Install Go** (version 1.21 or later)
2. **Install dependencies**:
   ```bash
   go mod tidy
   ```
3. **Run the application**:
   ```bash
   go run backend/main.go
   ```
4. **Access the application**: http://localhost:8080

### Database Schema

The application uses SQLite with the following tables:

#### massage_types
```sql
CREATE TABLE massage_types (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    duration INTEGER NOT NULL,
    price REAL NOT NULL
);
```

#### time_slots
```sql
CREATE TABLE time_slots (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL,
    time TEXT NOT NULL,
    service_id INTEGER NOT NULL,
    available INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (service_id) REFERENCES massage_types(id)
);
```

#### bookings (for future stories)
```sql
CREATE TABLE bookings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    client_name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT NOT NULL,
    service_id INTEGER NOT NULL,
    date TEXT NOT NULL,
    time_slot TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (service_id) REFERENCES massage_types(id)
);
```

## Sample Data

The application comes pre-loaded with sample data:

### Massage Services
- **Swedish Massage**: 60 minutes, €50
- **Deep Tissue**: 90 minutes, €70
- **Hot Stone**: 60 minutes, €65
- **Sports Massage**: 45 minutes, €45

### Time Slots
- Available from 09:00 to 18:00 in 15-minute intervals
- Generated for the next 30 days
- Approximately 30% of slots are randomly marked as booked for demonstration

## Testing Checklist

Before deploying, verify that:

- ✅ `docker compose up` starts the application successfully
- ✅ Application is accessible at http://localhost:8080
- ✅ Massage services load and display correctly
- ✅ Clicking a service shows the calendar
- ✅ Available slots are displayed in green (#4CAF50)
- ✅ Booked slots are displayed in grey (#9E9E9E) with strikethrough
- ✅ Available slots can be selected (highlighted when clicked)
- ✅ Booked slots cannot be interacted with
- ✅ Tooltip "This time is unavailable" appears on booked slots
- ✅ Data persists after application restart

## Future Development

This is Story #1 of a 3-story project. Upcoming features include:

- **Story #2**: Contact information form for completing bookings
- **Story #3**: Booking confirmation page and email notifications

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/story-X-description`
3. Make your changes
4. Test thoroughly using the checklist above
5. Commit your changes: `git commit -m "feat: implement story X - description"`
6. Push to the branch: `git push origin feature/story-X-description`
7. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For questions or issues, please open an issue on the GitHub repository.