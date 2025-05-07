import { getElement, shortcut } from './shortcut.js';
import { initSocket } from './websocket.js';
const socket = initSocket();
const sc = shortcut();
sc.scanElements();
sc.addHandler({
    openChat(evt) {
        const target = document.getElementById('chatBody');
        if (target && target.children.length > 0) {
            target.innerHTML = '';
            return;
        }
        const cid = evt.target.getAttribute('data-cid');
        if (!cid) {
            return;
        }
        const toUuid = cid.split('_')[1];
        const chatTemplate = document.getElementById('chatTemplate');
        if (!chatTemplate) {
            return;
        }
        const clone = sc.getTemplateClone('chatTemplate');
        if (!target) {
            return;
        }
        if (!clone) {
            return;
        }
        target.appendChild(clone);
        const chatSc = sc.appendScanElements(target);
        chatSc.addHandler({
            handleInput(evt) {
                const event = evt;
                if (event.key === 'Enter' && toUuid) {
                    const inputElm = event.target;
                    socket.send(JSON.stringify({
                        toUuid,
                        msg: inputElm.value
                    }));
                    inputElm.value = '';
                }
            }
        }).setActions();
    }
});
sc.setActions();
export default {
    connect(wsUrl) {
        socket.connect(wsUrl, {
            onOpen(evt) {
                console.log(evt, 'event onopen');
            },
            onMessage(event) {
                const parsedSocketData = JSON.parse(event.data);
                let element;
                switch (parsedSocketData.event) {
                    case 'isOnline':
                        const online = JSON.parse(parsedSocketData.message);
                        for (const uuid of online) {
                            element = getElement(`isOnline_${uuid}`);
                            element.textContent = 'Online';
                        }
                        break;
                    case 'join':
                        element = getElement(`isOnline_${parsedSocketData.uuid}`);
                        element.textContent = 'Online';
                        break;
                    case 'quit':
                        element = getElement(`isOnline_${parsedSocketData.uuid}`);
                        element.textContent = 'Offline';
                        break;
                    case 'message':
                        const messages = getElement('messages');
                        if (!messages) {
                            return;
                        }
                        const container = document.createElement('div');
                        container.classList.add('message');
                        const date = document.createElement('div');
                        date.classList.add('message-date');
                        const from = document.createElement('div');
                        from.classList.add('message-from');
                        const p = document.createElement('p');
                        p.classList.add('message-text');
                        from.innerText = `${parsedSocketData.username}`;
                        p.innerText = `${parsedSocketData.message}`;
                        date.innerText = `${new Date(parsedSocketData.createdAt).toLocaleString('en-US', {
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
                        messages.appendChild(container);
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
