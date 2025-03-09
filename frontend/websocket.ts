export type SocketController = {
  connect: () => void;
  disconnect: () => void;
  send: (message: string) => void;
};

export const initSocket = (
  username: string,
  {
    onClose,
    onError,
    onOpen,
    onMessage,
  }: {
    onOpen?: () => void;
    onMessage?: (message: MessageEvent) => void;
    onClose?: (e: CloseEvent) => void;
    onError?: (e: Event) => void;
  },
): SocketController => {
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

    send: (message: string) => {
      if (!ws) {
        return;
      }
      ws.send(message);
    },
  };
};
