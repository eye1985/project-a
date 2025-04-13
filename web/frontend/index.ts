import type { SocketMessage } from './websocket.js';
import { initSocket, SocketController } from './websocket.js';
import { shortcut } from './shortcut.js';

export const init = (wsUrl: string, username: string) => {
  if (!wsUrl) {
    throw new Error('Invalid domain');
  }

  let socket: SocketController;
  const sc = shortcut();

  const channelInput = sc.getElement('channelInput') as HTMLInputElement;
  const connectButton = sc.getElement('connectToChatBtn') as HTMLButtonElement;
  const closeButton = sc.getElement('closeChatBtn') as HTMLButtonElement;
  const messageInput = sc.getElement('messageInput') as HTMLInputElement;
  const messages = sc.getElement('messages');

  sc.addHandler({
    connectWS: () => {
      socket = initSocket();
      socket.connect(
        username,
        channelInput.value.trim(),
        wsUrl,
        {
          onOpen(evt) {
            console.log(evt, 'event');
            closeButton.removeAttribute('disabled');
          },
          onMessage(event) {
            if (!messages) {
              throw new Error('Div not found');
            }
            const newMessage = document.createElement('p');
            const { message, username, createdAt } = JSON.parse(
              event.data
            ) as SocketMessage;

            const time = new Date(createdAt);
            const timeStamp = `${time.getDate() < 10 ? '0' + time.getDate() : time.getDate()}/${time.getMonth()}/${time.getFullYear()} ${time.getHours()}:${time.getMinutes()}:${time.getSeconds()}`;
            newMessage.innerText = `${timeStamp} - ${username}: ${message}`;

            messages.appendChild(newMessage);
            messages.scrollTo(0, messages.scrollHeight);
          },

          onClose(evt) {
            if (evt.code === 1008) {
              alert(evt.reason);
            }

            connectButton.removeAttribute('disabled');
            closeButton.setAttribute('disabled', 'disabled');
            messageInput.setAttribute('disabled', 'disabled');
          }
        }
      );
      connectButton.setAttribute('disabled', 'disabled');
      messageInput.removeAttribute('disabled');
    },
    closeWS: () => {
      const connectButton = sc.getElement(
        'connectToChatBtn'
      ) as HTMLButtonElement;
      socket.disconnect();
      connectButton.removeAttribute('disabled');
    },
    handleInput(e) {
      if ((e as KeyboardEvent).key === 'Enter') {
        socket.send(messageInput.value);
        messageInput.value = '';
      }
    }
  });
  sc.init();
};
