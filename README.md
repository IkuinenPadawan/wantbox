# Wantbox
A simple, lightweight wishlist application targeted for self-host audience. Track gift ideas, share wishlists, and never forget what your loved ones really want.

## Features
* Multi-user wishlist cards
* Add/edit/delete items with prices and URLs
* Comments/Notes - Add comments/notes to item
* SQLite database (no setup required)
* Docker support

## Requirements

* Go 1.24+ (for local development)
* Docker & Docker Compose (for containerized deployment)

## Quick Start

### Option 1: Docker Compose (Recommended)

```bash
# Clone the repository
git clone https://github.com/IkuinenPadawan/wantbox.git
cd wantbox

# Start the application
docker-compose up -d

# Visit http://localhost:8089
```

### Option 2: Run Locally

```bash
# Prerequisites: Go 1.24+
git clone https://github.com/IkuinenPadawan/wantbox.git
cd wantbox

# Install dependencies and run
go mod download
go run main.go

# Visit http://localhost:8089
```

### Option 3: Build Binary

```bash
# Prerequisites: Go 1.24+
git clone https://github.com/IkuinenPadawan/wantbox.git
cd wantbox

# Build and run
go build -o wantbox .
./wantbox

# Visit http://localhost:8089
```
## Tech Stack

**Backend**: Go 1.24 with Gin web framework
**Database**: SQLite (embedded, zero-config)
**Frontend**: HTML templates with responsive CSS
**Deployment**: Docker & Docker Compose

## Roadmap

**User Authentication** - Login system for secure access
**Categories/Tags** - Organize items by category
**Images** - Add images of wishlist items
**Styling Overhaul** - Uniform sleek design
**Mobile Friendly** - Responsive design
**Import/Export** - CSV import/export functionality


