import { addFormMethod, addFromTarget, addHandler, getElement, } from './shortcut.js';
addFromTarget(document.body);
addFormMethod('inviteOnError', (errorMsg) => {
    (getElement('inviteError')?.ref).innerText = errorMsg;
});
addFormMethod('inviteOnSuccess', () => {
    location.reload();
});
addHandler('goToChat', () => {
    location.href = '/chat';
});
