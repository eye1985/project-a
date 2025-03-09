export type SocketController = {
  connect: (
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
  ) => void;
  disconnect: () => void;
  send: (message: string) => void;
};

export const initSocket = (): SocketController => {
  let ws: WebSocket;

  return {
    connect(
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
    ) {
      ws = new WebSocket(`ws://localhost:8080/ws?username=${username}`);

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
