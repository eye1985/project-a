export const initSocket = (username, { onClose, onError, onOpen, onMessage, }) => {
    let ws = new WebSocket(`ws://localhost:8080/ws?username=${username}`);
    ws.onopen = () => {
        onOpen && onOpen();
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
    return {
        connect() {
            ws = new WebSocket(`ws://localhost:8080/ws?username=${username}`);
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
