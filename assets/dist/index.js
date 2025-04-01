import { initSocket } from './websocket.js';
const elements = Array.from(document.querySelectorAll('[data-cid]'));
let socket;
const getElement = (id) => {
    let element = elements.find((element) => element.getAttribute('data-cid') === id);
    if (!element) {
        return null;
    }
    return element;
};
elements.forEach((element) => {
    switch (element.getAttribute('data-cid')) {
        case 'connectToChatBtn':
            element.addEventListener('click', () => {
                const userNameInput = getElement('usernameInput');
                const connectButton = getElement('connectToChatBtn');
                const messageInput = getElement('messageInput');
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
                        const messages = elements.find((element) => element.getAttribute('data-cid') === 'messages');
                        if (!messages) {
                            throw new Error('Div not found');
                        }
                        const newMessage = document.createElement('p');
                        const { message, username, createdAt } = JSON.parse(event.data);
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
            const connectButton = getElement('connectToChatBtn');
            element.addEventListener('click', () => {
                socket.disconnect();
                connectButton.removeAttribute('disabled');
            });
            break;
        case 'messageInput':
            element.addEventListener('keypress', (e) => {
                if (e.key === 'Enter') {
                    const messageInput = getElement('messageInput');
                    socket.send(messageInput.value);
                    messageInput.value = '';
                }
            });
            break;
    }
});
