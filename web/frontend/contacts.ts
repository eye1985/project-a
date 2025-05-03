import { shortcut } from './shortcut.js';
import { initSocket, SocketMessage } from './websocket.js';

const sc = shortcut();
sc.init();

const socket = initSocket();

export default {
  connect(wsUrl: string) {
    socket.connect(wsUrl, {
      onOpen(evt) {
        console.log(evt, 'event onopen');
      },
      onMessage(event) {
        const parsedSocketData: SocketMessage = JSON.parse(event.data);
        let element: Element;
        switch (parsedSocketData.event) {
          case 'isOnline':
            const online = JSON.parse(parsedSocketData.message);
            for (const uuid of online) {
              element = sc.getElement(`isOnline_${uuid}`);
              element.textContent = 'Online';
            }

            break;
          case 'join':
            element = sc.getElement(`isOnline_${parsedSocketData.uuid}`);
            element.textContent = 'Online';
            break;
          case 'quit':
            element = sc.getElement(`isOnline_${parsedSocketData.uuid}`);
            element.textContent = 'Offline';
            break;
        }
      },

      onClose(evt) {
        console.log(evt, 'on close');
      },

      onError(evt) {
        console.error(evt, 'on error');
      }
    });
  }
};