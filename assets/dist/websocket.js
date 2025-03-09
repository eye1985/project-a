export const initSocket = () => {
    let ws;
    return {
        connect(username, { onClose, onError, onOpen, onMessage, }) {
            ws = new WebSocket(`ws://localhost:8080/ws?username=${username}`);
            ws.onopen = (evt) => {
                onOpen && onOpen(evt);
            };
            ws.onmessage = (event) => {
                onMessage && onMessage(event);
            };
            ws.onclose = (event) => {
                onClose && onClose(event);
            };
            ws.onerror = (error) => {
                onError && onError(error);
            };
        },
        disconnect() {
            if (!ws) {
                return;
            }
            ws.close();
        },
        send: (message) => {
            if (!ws) {
                return;
            }
            ws.send(message);
        },
    };
};
