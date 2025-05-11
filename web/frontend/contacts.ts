import { addFromTarget, addHandler, CustomElement, getElement, isTemplate, state } from './shortcut.js';
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

  const toUuid = currentCustomElement.id.split('_')[1];
  set('toUuid', toUuid);
  if (!toUuid) {
    throw new Error('toUuid not found');
  }
  const history = sessionStorage.getItem(toUuid);
  if (history) {
    const messages = getElement('messages');
    if (!messages) {
      throw new Error('messages not found');
    }
    JSON.parse(history).forEach((data: SocketMessage) => {
      const isMyMessage = data.fromUuid === get('myUuid');
      insertMessage(data, messages.ref, isMyMessage);
    });
  }
});

addHandler('handleInput', (e) => {
  const event = e as KeyboardEvent;
  const toUuid = get('toUuid');
  const inputElm = event.currentTarget as HTMLInputElement;

  if (event.key === 'Enter' && event.shiftKey) {
    return;
  }

  if (event.key === 'Enter' && toUuid && inputElm.value.trim().length > 0) {
    socket.send(JSON.stringify({
      toUuid,
      msg: inputElm.value
    }));
    inputElm.value = '';
  }
});

const insertMessage = (data: SocketMessage, target: Element, isCurrentUser: boolean) => {
  const message = getElement('messageTemplate');
  if (!message) {
    throw new Error('messageTemplate not found');
  }

  if (isTemplate(message.ref)) {
    const date = message.ref.content.querySelector('.message-date') as HTMLDivElement;
    const from = message.ref.content.querySelector('.message-from') as HTMLDivElement;
    const text = message.ref.content.querySelector('.message-text') as HTMLDivElement;

    date.innerText = `${new Date(data.createdAt).toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: 'numeric',
      minute: 'numeric',
      second: 'numeric',
      hour12: false
    })}`;

    from.innerText = `${data.username}`;
    text.innerText = `${data.message}`;
  }

  const classNames = isCurrentUser ? ['message', 'me'] : ['message'];
  message.insertTemplateInto(target, {
    clearBeforeInsert: false,
    classNames
  });
};

const insertChatHistory = (toUuid: string, data: SocketMessage) => {
  const history = sessionStorage.getItem(toUuid);
  if (!history) {
    sessionStorage.setItem(toUuid, JSON.stringify([data]));
  } else {
    sessionStorage.setItem(toUuid, JSON.stringify([...JSON.parse(history), data]));
  }
};
export default {
  connect(wsUrl: string, myUuid: string) {
    if (!myUuid) {
      throw new Error('myUuid not found');
    }
    set('myUuid', myUuid);
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
                element = getElement(`isOnline_${data.fromUuid}`);
                if (element) {
                  element.ref.textContent = 'Online';
                }
                break;
              case 'quit':
                element = getElement(`isOnline_${data.fromUuid}`);
                if (element) {
                  element.ref.textContent = 'Offline';
                }
                break;
              case 'message':
                const toUuid = get('toUuid');
                const messages = getElement('messages');
                const isMyMessage = data.fromUuid === myUuid;
                const isSystemMsgToMe = data.username === 'System' && data.fromUuid === myUuid;

                const isMessageToThisUser = data.fromUuid === toUuid || data.toUuid === toUuid;
                if (data.username !== 'System' && isMessageToThisUser) {
                  insertChatHistory(toUuid, data);
                }

                if (!messages) {
                  insertChatHistory(data.fromUuid, data);
                  console.warn('cant find messages');
                  return;
                }

                if (isMessageToThisUser || isSystemMsgToMe) {
                  insertMessage(data, messages.ref, isMyMessage);
                }
                break;
            }
          }
        );

        const messages = getElement('messages');
        if (messages) {
          requestAnimationFrame(() => {
            messages.ref.parentElement?.scrollTo({
              top: messages.ref.parentElement?.scrollHeight,
              behavior: 'smooth'
            });
          });
        }
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