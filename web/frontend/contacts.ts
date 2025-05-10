import { addFromTarget, addHandler, CustomElement, getElement, state } from './shortcut.js';
import { initSocket, SocketMessage } from './websocket.js';

const { get, set } = state;
const socket = initSocket();
addFromTarget(document.body);

addHandler('openChat', (e, currentCustomElement, store) => {
  const buttons = store.getByType('chat-button');
  buttons.forEach(button => {
    button.ref.classList.remove('active');
  });

  const currentButton = (e.currentTarget as HTMLButtonElement);
  currentButton.classList.add('active');

  const chatBody = store.elements.get('chatBody');
  if (!chatBody) {
    throw new Error('chatBody not found');
  }

  const template = store.elements.get('chatTemplate');
  if (!template) {
    throw new Error('chatTemplate not found');
  }

  template.insertTemplateInto(chatBody, {
    clearBeforeInsert: true
  });
  set('toUuid', currentCustomElement.id.split('_')[1]);
});

addHandler('handleInput', (e) => {
  const event = e as KeyboardEvent;
  const toUuid = get('toUuid');

  if (event.key === 'Enter' && toUuid) {
    const inputElm = event.currentTarget as HTMLInputElement;
    socket.send(JSON.stringify({
      toUuid,
      msg: inputElm.value
    }));
    inputElm.value = '';
  }
});

export default {
  connect(wsUrl: string) {
    socket.connect(wsUrl, {
      onOpen(evt) {
        console.log(evt, 'event onopen');
      },
      onMessage(event) {
        const parsedSocketData: SocketMessage[] = JSON.parse(event.data);
        let element: CustomElement | null;

        parsedSocketData.forEach(
          (data: SocketMessage) => {
            switch (data.event) {
              case 'isOnline':
                const online = JSON.parse(data.message);
                for (const uuid of online) {
                  element = getElement(`isOnline_${uuid}`);
                  if (element) {
                    element.ref.textContent = 'Online';
                  }
                }
                break;
              case 'join':
                element = getElement(`isOnline_${data.uuid}`);
                if (element) {
                  element.ref.textContent = 'Online';
                }
                break;
              case 'quit':
                element = getElement(`isOnline_${data.uuid}`);
                if (element) {
                  element.ref.textContent = 'Offline';
                }
                break;
              case 'message':
                const messages = getElement('messages');
                if (!messages) {
                  return;
                }

                // TODO Use template or something
                const container = document.createElement('div');
                container.classList.add('message');
                const date = document.createElement('div');
                date.classList.add('message-date');
                const from = document.createElement('div');
                from.classList.add('message-from');
                const p = document.createElement('p');
                p.classList.add('message-text');
                from.innerText = `${data.username}`;
                p.innerText = `${data.message}`;
                date.innerText = `${new Date(data.createdAt).toLocaleString('en-US', {
                  month: 'short',
                  day: 'numeric',
                  hour: 'numeric',
                  minute: 'numeric',
                  second: 'numeric',
                  hour12: false
                })}`;

                container.appendChild(from);
                container.appendChild(p);
                container.appendChild(date);

                messages.ref.appendChild(container);

                messages.ref.scrollTo({
                  top: messages.ref.scrollHeight,
                  behavior: 'smooth'
                });
                break;
            }
          }
        );
      },

      onClose(evt) {
        console.log(evt, 'on close');
      },

      onError(evt) {
        console.error(evt, 'on error');
        const template = getElement('toast');
        if (!template) {
          throw new Error('toast not found');
        }

        const p = (template.ref as HTMLTemplateElement).content.querySelector('p');
        if (!p) {
          throw new Error('toast p not found');
        }
        p.innerText = 'Could not connect to server. Please try again later.';
        template.insertTemplateInto(document.body);

        setTimeout(() => {
          template.remove();
        }, 2000);
      }
    });
  }
};