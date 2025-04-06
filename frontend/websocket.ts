export type SocketMessage = {
  message: string;
  username: string;
  createdAt: Date;
};

export type SocketController = {
  connect: (
    username: string,
    channel: string,
    url: string,
    {
      onClose,
      onError,
      onOpen,
      onMessage,
    }: {
      onOpen?: (e: Event) => void;
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
      channel: string,
      url,
      {
        onClose,
        onError,
        onOpen,
        onMessage,
      }: {
        onOpen?: (e: Event) => void;
        onMessage?: (message: MessageEvent) => void;
        onClose?: (e: CloseEvent) => void;
        onError?: (e: Event) => void;
      },
    ) {
      ws = new WebSocket(
        `ws://${url}/ws?username=${username}&channels=${channel}`,
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
    },
  };
};
