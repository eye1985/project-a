export const initSocket = () => {
    let ws;
    return {
        connect(username, channel, url, { onClose, onError, onOpen, onMessage, }) {
            ws = new WebSocket(`ws://${url}/ws?username=${username}&channels=${channel}`);
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
