import {
  addFormMethod,
  addFromTarget,
  addHandler,
  getElement,
} from './shortcut.js';

addFromTarget(document.body);
addFormMethod('inviteOnError', (errorMsg: string) => {
  (getElement('inviteError')?.ref as HTMLDivElement).innerText = errorMsg;
});

addFormMethod('inviteOnSuccess', () => {
  location.reload();
});

addHandler('goToChat', () => {
  location.href = '/chat';
});
