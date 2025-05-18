import { addFormMethod, scanFrom, addHandler, getElement } from './shortcut.js';
scanFrom(document.body);
addFormMethod('inviteOnError', (errorMsg) => {
    (getElement('inviteError')?.ref).innerText = errorMsg;
});
addFormMethod('inviteOnSuccess', () => {
    location.reload();
});
addHandler('goToChat', () => {
    location.href = '/chat';
});
