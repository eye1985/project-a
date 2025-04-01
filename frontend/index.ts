import { initSocket, SocketController } from './websocket.js';

const elements = Array.from(document.querySelectorAll('[data-cid]'));

let socket: SocketController;

const getElement = (id: string) => {
  let element = elements.find(
    (element) => element.getAttribute('data-cid') === id,
  );

  if (!element) {
    return null;
  }

  return element;
};

type SocketMessage = {
  message: string;
  username: string;
  createdAt: Date;
};

elements.forEach((element) => {
  switch (element.getAttribute('data-cid')) {
    case 'connectToChatBtn':
      element.addEventListener('click', () => {
        const userNameInput = getElement('usernameInput') as HTMLInputElement;
        const connectButton = getElement(
          'connectToChatBtn',
        ) as HTMLButtonElement;
        const messageInput = getElement('messageInput') as HTMLInputElement;

        if (!userNameInput) {
          throw new Error('User name input not found');
        }

        if (!userNameInput.value) {
          console.error('No username');
          return;
        }

        socket = initSocket();
        socket.connect(userNameInput.value, {
          onOpen(evt) {
            console.log(evt, 'event');
          },
          onMessage(event) {
            const messages = elements.find(
              (element) => element.getAttribute('data-cid') === 'messages',
            ) as HTMLDivElement | undefined;

            if (!messages) {
              throw new Error('Div not found');
            }
            const newMessage = document.createElement('p');
            const { message, username, createdAt } = JSON.parse(
              event.data,
            ) as SocketMessage;

            const time = new Date(createdAt);
            const timeStamp = `${time.getDate() < 10 ? '0' + time.getDate() : time.getDate()}/${time.getMonth()}/${time.getFullYear()} ${time.getHours()}:${time.getMinutes()}:${time.getSeconds()}`;
            newMessage.innerText = `${timeStamp} - ${username}: ${message}`;

            messages.appendChild(newMessage);
            messages.scrollTo(0, messages.scrollHeight);
          },

          onClose() {
            connectButton.removeAttribute('disabled');
          },
        });
        connectButton.setAttribute('disabled', 'disabled');
        messageInput.removeAttribute('disabled');
      });
      break;
    case 'closeChatBtn':
      const connectButton = getElement('connectToChatBtn') as HTMLButtonElement;
      element.addEventListener('click', () => {
        socket.disconnect();
        connectButton.removeAttribute('disabled');
      });

      break;
    case 'messageInput':
      element.addEventListener('keypress', (e) => {
        if ((e as KeyboardEvent).key === 'Enter') {
          const messageInput = getElement('messageInput') as HTMLInputElement;
          socket.send(messageInput.value);
          messageInput.value = '';
        }
      });
      break;
  }
});
