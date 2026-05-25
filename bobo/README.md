# LoL Live Match Visualization Server

Beautiful real-time visualization of League of Legends live match analytics.

## Features

✨ **Real-time Dashboard**
- Live player and opponent statistics
- KDA, CS, Gold tracking
- Level, Attack Damage, Armor, Magic Resist display
- Team statistics overview

📊 **Multi-chart Visualization**
- 💰 **Gold Over Time** - Track economy progress
- 🌾 **Creep Score (CS)** - Monitor farm efficiency
- ⚔️ **Kills** - View kill progression
- 💀 **Deaths** - Track death statistics
- 🤝 **Assists** - Monitor team cooperation

📜 **Game Events Timeline**
- Real-time event logging
- Champion kills, turret destruction, and more
- Scrollable event history with timestamps

## Getting Started

### Requirements
- Go 1.16 or higher
- Modern web browser (Chrome, Firefox, Safari, Edge)

### Build & Run

```bash
cd bobo
go build -v
./bobo  # On Windows: .\bobo.exe
```

Server starts on `http://localhost:8080`

### API Endpoints

#### WebSocket
- **`ws://localhost:8080/ws`** - Subscribe to live updates

#### REST API
- **`POST /ingest`** - Receive match analytics (OpenAPI compliant)
- **`POST /api/update`** - Alternative endpoint (same as /ingest)
- **`GET /`** - Serve dashboard UI

### JSON Payload Format (OpenAPI Schema)

```json
{
  "timestamp": 1716656220000,
  "player": {
    "summonerName": "string",
    "champion": "string",
    "totalGold": 5500.0,
    "kills": 5,
    "deaths": 1,
    "assists": 8,
    "cs": 180,
    "level": 10,
    "gameTime": 600.0,
    ...
  },
  "opponent": { /* same structure */ },
  "events": [
    {
      "type": "ChampionKill",
      "time": 120.5,
      "killer": "PlayerName",
      "victim": "EnemyName"
    }
  ]
}
```

### Testing

Send test data:
```bash
$json = Get-Content test_data.json -Raw
Invoke-WebRequest -Uri "http://localhost:8080/ingest" `
  -Method POST -Body $json `
  -ContentType "application/json"
```

Then open `http://localhost:8080` in your browser.

## Architecture

### Backend (Go)
- **WebSocket Server** - Real-time data streaming to browsers
- **HTTP REST API** - Receive analytics payload from parser
- **Thread-safe** - Uses sync.Mutex for concurrent client management
- **State Management** - Maintains last payload for new connections

### Frontend (HTML/Chart.js)
- **Responsive Design** - Works on desktop and tablet
- **Tailwind CSS** - Modern, dark-themed UI
- **Chart.js** - Beautiful, performant charts
- **Real-time Updates** - WebSocket-driven UI refresh
- **Auto-scaling** - Charts show last 30 data points

## Configuration

Edit `main.go` to change:
- Server port (default: 8080)
- WebSocket connection parameters
- CORS settings

## UI Customization

The dashboard is fully responsive and themeable. Key color classes:
- **Player**: Blue theme (`from-blue-600 to-blue-500`)
- **Opponent**: Red theme (`from-red-600 to-red-500`)
- **Charts**: Dark theme with contrasting colors

## Performance Notes

- Charts maintain max 30 data points for smooth performance
- WebSocket clients cached for efficient broadcasting
- Minimal DOM updates using event-driven rendering
- No external dependencies except Chart.js and Tailwind CSS

## Development

### Add New Charts

In `index.html`, add a new canvas and Chart.js instance:

```javascript
const charts = {
  newChart: new Chart(ctx, {
    type: 'line',
    data: { /* ... */ },
    options: chartConfig
  })
};
```

### Extend Payload

Update Go structures in `main.go` to match OpenAPI schema in `../contracts/lol-tracker.openapi.yaml`

## License

Part of lol-analytics project
