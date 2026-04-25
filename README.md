# Trust Me Bro It's Not Fake AI Agent

<p align="center">
  <img src="./docs/images/demo.gif" alt="Trust Me Bro Demo" width="800"/>
</p>

> **"100% Real AI, Not Fake At All"** — Our marketing department

---

## What's This?

A chat application where a **real person** on the other end pretends to be an AI agent. Revolutionary, right?

- **TUI Client** — Talk to your "AI agent" from the terminal
- **Web Dashboard** — Where the poor human sits, pretending to be AI
- **Zero hallucinations** — Just pure, unfiltered human intelligence (the human makes mistakes too)

---

## How It Works

```
┌─────────┐         ┌─────────────┐         ┌──────────┐
│   TUI   │────────▶│  RabbitMQ   │────────▶│  Backend │
│(Terminal)│         │  (Broker)   │         │  (Go)    │
└─────────┘         └─────────────┘         └────┬─────┘
                                                 │
                                                 ▼
┌─────────┐                               ┌──────────┐
│   Web   │◀──────────────────────────────│ WebSocket│
│(Dashboard)│         The Human           │          │
└─────────┘                               └──────────┘
```

The TUI user types → RabbitMQ queues it → Backend receives → Web dashboard shows message → **Human reads and replies** → Message goes back → TUI displays response.

It's like ChatGPT, but the "GPT" stands for **G**uy **P**retending **T**o be smart.

---

## ✨ Features

- **🖥️ Terminal Interface** — Chat from your terminal like a hacker in movies
- **🌐 Web Dashboard** — Clean UI for the "AI operator" to respond
- **🐇 Message Queue** — RabbitMQ handles the queue, humans handle the thinking
- **💾 Persistent Storage** — PostgreSQL remembers everything (including your mistakes)
- **⚡ Real-Time Updates** — WebSocket for the web side, RabbitMQ for the terminal side

---

## 🛠️ Tech Stack

| Layer | Technology |
|-------|------------|
| **TUI Client** | Go, Bubble Tea |
| **Backend API** | Go |
| **Web Dashboard** | React 19, TypeScript, Vite, Redux Toolkit |
| **UI Components** | Ant Design, Tailwind CSS |
| **Message Broker** | RabbitMQ 4 |
| **Database** | PostgreSQL 16 |
| **Real-time** | WebSocket, Socket.io |

---

## 📋 TODO

- [ ] **Pretty the Web UI** — Make the dashboard look less like a prototype
- [ ] **Pretty the TUI** — Terminal interface needs some love
- [ ] **Tool Calling** — Add function calling (like real AI agents do): files, terminal commands

---

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- A human willing to pretend they're AI

### Run Everything

```bash
docker compose up -d
```

### Open the Dashboard

```bash
# Go to this URL to be the "AI"
open http://localhost:5173
```

### Run the Terminal Chat

```bash
docker compose --profile tui run --rm --build tui
```

### Service URLs

| Service | URL | Purpose |
|---------|-----|---------|
| Web Dashboard | http://localhost:5173 | Where the human acts as "AI" |
| API | http://localhost:8080 | Backend server |
| RabbitMQ | http://localhost:15672 | Queue management (`guest/guest`) |
| PostgreSQL | localhost:5432 | Message storage (`admin/secret`) |

### Stop

```bash
docker compose down
```

---

## 📁 Project Structure

```
.
├── web/              # React dashboard (for the "AI operator")
├── backend/          # Go API server
├── tui/              # Terminal client (for the curious user)
├── shared/           # Shared Go packages
├── docs/images/      # Demo gifs
└── docker-compose.yml
```

---

## 🎭 The Philosophy

> *"Is it AI? We can't confirm or deny. But hey, the human on the other side needs a job."*

Built with Docker, slightly misleading naming, and a human who needed work.

**License**: MIT (you can pretend it's AI-generated if you want)
