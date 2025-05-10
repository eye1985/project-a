export type SocketMessage = {
  fromUuid: string;
  toUuid: string;
  message: string;
  event: string;
  username: string;
  createdAt: Date;
};

export type SocketController = {
  connect: (
    wsUrl: string,
    {
      onClose,
      onError,
      onOpen,
      onMessage
    }: {
      onOpen?: (e: Event) => void;
      onMessage?: (message: MessageEvent) => void;
      onClose?: (e: CloseEvent) => void;
      onError?: (e: Event) => void;
    }
  ) => void;
  disconnect: () => void;
  send: (message: string) => void;
};

export const initSocket = (): SocketController => {
  let ws: WebSocket;

  return {
    connect(
      wsUrl,
      {
        onClose,
        onError,
        onOpen,
        onMessage
      }: {
        onOpen?: (e: Event) => void;
        onMessage?: (message: MessageEvent) => void;
        onClose?: (e: CloseEvent) => void;
        onError?: (e: Event) => void;
      }
    ) {
      ws = new WebSocket(
        `${wsUrl}/ws`
      );

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

    send: (message: string) => {
      if (!ws) {
        return;
      }
      ws.send(message);
    }
  };
};
