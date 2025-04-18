export type SocketMessage = {
  message: string;
  username: string;
  createdAt: Date;
};

export type SocketController = {
  connect: (
    username: string,
    channel: string,
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
      username: string,
      channel: string,
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

      try {
        ws = new WebSocket(
          `${wsUrl}/ws?username=${username}&channels=${channel}`
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
      } catch (error) {
        const err = error as Error;
        throw new Error(err.message);
      }
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
