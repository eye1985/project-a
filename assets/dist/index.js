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
                socket = initSocket(userNameInput.value, {
                    onMessage: (event) => {
                        const messages = elements.find((element) => element.getAttribute('data-cid') === 'messages');
                        if (!messages) {
                            throw new Error('Div not found');
                        }
                        const newMessage = document.createElement('p');
                        newMessage.innerText = event.data;
                        messages.appendChild(newMessage);
                        messages.scrollTo(0, messages.scrollHeight);
                    },
                });
                socket.connect();
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
