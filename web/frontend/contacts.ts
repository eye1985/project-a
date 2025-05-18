import { addFormMethod, scanFrom, addHandler, getElement } from './shortcut.js';

scanFrom(document.body);
addFormMethod('inviteOnError', (errorMsg: string) => {
  (getElement('inviteError')?.ref as HTMLDivElement).innerText = errorMsg;
});

addFormMethod('inviteOnSuccess', () => {
  location.reload();
});

addHandler('goToChat', () => {
  location.href = '/chat';
});
