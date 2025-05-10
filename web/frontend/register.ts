import { addFormMethod, addFromTarget, getCookie, getElement } from './shortcut.js';

const createToast = (text: string, timeout?: number) => {
  const template = getElement('toast');
  if (template) {
    const p = (template.ref as HTMLTemplateElement).content.querySelector('p');
    if (!p) {
      throw new Error('p not found');
    }
    p.innerText = text;
    template.insertTemplateInto(document.body);
    if (timeout) {
      setTimeout(() => {
        template.remove();
      }, timeout);
    }
  }
};

addFromTarget(document.body);
addFormMethod('onerror', (error) => {
  createToast(error, 2000);
});

addFormMethod('onsuccess', (success) => {
  console.log(success);
  const layout = document.querySelector('.layout');
  if (!layout) {
    return;
  }

  const template = getElement('emailSent');
  if (template) {
    template.insertTemplateInto(layout, {
      clearBeforeInsert: true,
      classNames: ['layout-item']
    });
  }
});

const cookie = getCookie('flash');
if (cookie) {
  createToast(cookie.value + '', 2000);
}