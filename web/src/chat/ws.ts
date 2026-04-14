let ws: WebSocket | null = null;

export const getWs = () => {
	if (!ws || ws.readyState === WebSocket.CLOSED) {
		ws = new WebSocket("ws://localhost:8080/ws");
	}
	return ws;
}

export const closeWs = () => {
	ws?.close();
	ws = null;
}
