# LoL Live Match Visualization - Complete Setup

## 🎯 How It Works

```
League of Legends Live API
    ↓
lol-tracker (parser)
    ↓ (автоматически отправляет /ingest)
bobo server (receiver)
    ↓ (WebSocket broadcast)
Browser Dashboard (real-time visualization)
```

---

## 🚀 Quick Start

### 1. Start Visualization Server

```bash
cd bobo
go build -v
./test.exe  # or .\test.exe on Windows

# Output: LoL Live Match Visualization Server запущен на http://localhost:8080
```

### 2. Start Tracker (in another terminal)

```bash
cd lol-trakcer
go run ./cmd/tracker/main.go

# Output: League Tracker started
#         Server: http://localhost:8080
```

The tracker automatically sends data to `/ingest` endpoint every 3 seconds.

### 3. Open Dashboard

```
http://localhost:8080
```

---

## 📊 What You See

- **Real-time Stats**: KDA, CS, Gold, Level, Damage, Armor, MR
- **Live Charts**: 
  - 💰 Gold progression
  - 🌾 CS growth
  - ⚔️ Kills
  - 💀 Deaths
  - 🤝 Assists
- **Events Timeline**: All game events with timestamps
- **Player vs Enemy**: Side-by-side comparison

---

## 🔧 Architecture

### Backend: bobo/main.go
- WebSocket server on `ws://localhost:8080/ws`
- REST endpoint `/ingest` (OpenAPI compliant)
- Broadcasts to all connected browsers
- Caches last payload for new connections

### Parser: lol-tracker/cmd/tracker/main.go
- Polls League client API every 3s
- Analyzes player & opponent data
- Sends ServerPayload to `/ingest` automatically
- No manual steps needed!

### Frontend: bobo/index.html
- Connects via WebSocket
- Real-time stat updates
- Chart.js for visualization
- No page jumps or auto-scroll (fixed!)
- Responsive design

---

## 📡 Data Flow

**Tracker sends JSON payload to `POST /ingest`:**

```json
{
  "timestamp": 1716656220000,
  "player": {
    "summonerName": "YourName",
    "champion": "Garen",
    "kills": 5,
    "deaths": 1,
    "assists": 8,
    "cs": 180,
    "totalGold": 5500,
    "level": 10,
    "attackDamage": 145,
    "armor": 55,
    "magicResist": 35,
    "gameTime": 600,
    ...
  },
  "opponent": { /* same structure */ },
  "events": [
    { "type": "ChampionKill", "time": 120, "killer": "You", "victim": "Enemy" }
  ]
}
```

---

## ⚙️ Configuration

Set environment variables (optional):

```bash
# Default is localhost:8080
export LOL_SERVER_URL=http://localhost:8080

# Default is https://127.0.0.1:2999
export LOL_LEAGUE_URL=https://127.0.0.1:2999

# Default is 3s
export LOL_POLL_INTERVAL=3s
```

---

## 🐛 Troubleshooting

### "Connection refused" on tracker
- Make sure bobo server is running first
- Check `LOL_SERVER_URL` environment variable
- Verify port 8080 is not blocked

### "Waiting for game..."
- Launch a League of Legends game first
- Tracker waits for game to start before sending data
- Check that LoL client API is accessible at `https://127.0.0.1:2999`

### Events not showing
- Events come from LoL API as they happen
- Refresh page if connected before game started

### Charts not updating
- WebSocket must stay connected
- Check browser console (F12) for errors
- Dashboard receives data every 3 seconds (configurable)

---

## 🎨 UI Features

✅ **Fixed Header** - Always visible
✅ **No Auto-Scroll** - Page stays in place when updating
✅ **Smooth Charts** - No animation stutters
✅ **Color Coded** - Blue = You, Red = Enemy
✅ **Dark Theme** - Eye-friendly for long sessions
✅ **Responsive** - Works on mobile/tablet
✅ **Real-time** - WebSocket updates (instant)

---

## 📝 Files Changed

### bobo/
- ✅ `main.go` - Rewritten for OpenAPI compliance
- ✅ `index.html` - Fixed UI scroll, added 5 charts
- ✅ `README.md` - Full documentation

### lol-trakcer/
- ✅ `internal/transport/sender.go` - Now sends to `/ingest`
- ✅ `internal/config/config.go` - Updated default URL
- ✅ All other files - Unchanged, no breaking changes

---

## 🔄 Data Update Cycle

1. **Game Event Happens** (e.g., kill, CS earned)
2. **Tracker Polls API** (every 3 seconds)
3. **Tracker Sends Data** to `/ingest` 
4. **Server Receives** and caches
5. **Server Broadcasts** to all WebSocket clients
6. **Browser Updates** stats & charts in real-time

---

## ✨ No Manual Scripts Needed!

❌ No need to run test scripts
❌ No need to manually send JSON
❌ No need to click "Update" buttons

✅ Just run tracker, it handles everything!

---

## 🎓 OpenAPI Compliance

All data structures match `contracts/lol-tracker.openapi.yaml`:
- ✅ ServerPayload format
- ✅ PlayerAnalytics fields
- ✅ EventAnalytics structure
- ✅ camelCase JSON naming

---

## 🚦 Status Indicators

**Terminal Output:**

Tracker:
```
League Tracker started
Server: http://localhost:8080
[sends data every 3s]
```

Server:
```
LoL Live Match Visualization Server запущен на http://localhost:8080
```

Browser Console:
```
✓ Connected to server
[updates arriving...]
```

---

**You're all set! Enjoy real-time League of Legends analytics! 🎮**
