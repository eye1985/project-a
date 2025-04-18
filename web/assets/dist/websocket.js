export const initSocket = () => {
    let ws;
    return {
        connect(username, channel, wsUrl, { onClose, onError, onOpen, onMessage }) {
            try {
                ws = new WebSocket(`${wsUrl}/ws?username=${username}&channels=${channel}`);
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
            }
            catch (error) {
                const err = error;
                throw new Error(err.message);
            }
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
        }
    };
};
