import { shortcut } from './shortcut.js';
import { initSocket } from './websocket.js';
const socket = initSocket();
const sc = shortcut();
sc.addHandler({
    openChat(evt) {
        const cid = evt.target.getAttribute('data-cid');
        if (!cid) {
            return;
        }
        const toUuid = cid.split('_')[1];
        const chatTemplate = document.getElementById('chatTemplate');
        if (!chatTemplate) {
            return;
        }
        const clone = chatTemplate.content.cloneNode(true);
        const target = document.getElementById('chatBody');
        if (!target) {
            return;
        }
        target.appendChild(clone);
        // Quick and dirty fix for now
        //const messages = target.querySelector('[data-cid=\'messages\']');
        const input = target.querySelector('[data-cid=\'messageInput\']');
        if (input) {
            const inputElm = input;
            inputElm?.addEventListener('keyup', (evt) => {
                if (evt.key === 'Enter' && toUuid) {
                    socket.send(JSON.stringify({
                        toUuid,
                        msg: inputElm.value
                    }));
                    inputElm.value = '';
                }
            });
        }
    }
    // handleInput(e) {
    //
    //   const messageInput = sc.getElement('messageInput') as HTMLInputElement;
    //
    //   if ((e as KeyboardEvent).key === 'Enter') {
    //     socket.send(messageInput.value);
    //     messageInput.value = '';
    //   }
    // }
});
sc.init();
export default {
    connect(wsUrl) {
        socket.connect(wsUrl, {
            onOpen(evt) {
                console.log(evt, 'event onopen');
            },
            onMessage(event) {
                console.log(event, 'event onmessage');
                const parsedSocketData = JSON.parse(event.data);
                let element;
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
