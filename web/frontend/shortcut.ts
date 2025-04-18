const addHandler = (el: Element, handlers: Record<string, EventListener>) => {
  const eventAndHandler = el.getAttribute('data-handler');
  if (!eventAndHandler) {
    return;
  }

  const [event, handlerName] = eventAndHandler.split(':');
  if (!handlerName || !event) {
    return;
  }

  const handler = handlers[handlerName] ? handlers[handlerName] : () => {
  };
  el.addEventListener(event, handler);
};

const addBindElements = (currentElement: Element) => {
  const bindValue = currentElement.getAttribute('data-bind');
  const bindAction = currentElement.getAttribute('data-bind-action');
  if (!bindValue || !bindAction) {
    return;
  }

  const bindRules = bindValue.split(',');

  const state: {
    [key: string]: any;
  } = {};

  for (const rule of bindRules) {
    const [cid, expr] = rule.split(':');
    if (!expr || !cid) {
      continue;
    }

    const elm = document.querySelector(`[data-cid='${cid}']`);
    if (!elm) {
      continue;
    }

    const trimmedExpr = expr.trim();
    // Not supported expression
    if (trimmedExpr.length !== 2) {
      continue;
    }

    const exprExec = (input: string) => {
      if (trimmedExpr.indexOf('>') !== -1) {
        return input.trim().length > parseInt(trimmedExpr[1]!);
      } else {
        return input.trim().length < parseInt(trimmedExpr[1]!);
      }
    };

    state[cid] = {
      elm,
      exprExec,
      bindAction
    };
  }

  for (const item in state) {
    state[item].elm.addEventListener('input', () => {
      const res: boolean[] = [];

      for (const cid in state) {
        res.push(state[cid].exprExec(state[cid].elm.value));
      }

      res.every((num) => num)
        ? currentElement.removeAttribute(bindAction)
        : currentElement.setAttribute(bindAction, bindAction);
    });
  }
};

export const shortcut = () => {
  const elements = Array.from(document.querySelectorAll('[data-cid]'));
  const handlerNames = elements.map(
    (e) => e.getAttribute('data-handler')?.split(':')[1]
  );

  let handlers: Record<string, EventListener> | null = null;

  return {
    addHandler(handlersArg: Record<string, EventListener>) {
      const isAllHandlersPresent = Object.keys(handlersArg).every(
        (handlerName) => handlerNames.includes(handlerName)
      );

      if (!isAllHandlersPresent) {
        throw new Error('Not all handlers present, did you remember to add data-handler attribute to your elements?');
      }

      handlers = handlersArg;
    },
    getElement(cid: string) {
      const elm = document.querySelector(`[data-cid='${cid}']`);
      if (!elm) {
        throw new Error(`Not found element ${cid}`);
      }

      return elm;
    },
    init() {
      for (const el of elements) {
        if (el instanceof HTMLFormElement) {
          el.addEventListener('submit', async (e) => {
            e.preventDefault();

            const action = el.getAttribute('action');
            const method = el.getAttribute('data-method');
            const successMessage = el.getAttribute('data-success-message');
            if (!method || !action) {
              return;
            }

            const methods = ['POST', 'PUT', 'DELETE', 'PATCH'];
            const upperCasedMethod = method.toUpperCase();
            if (!methods.includes(upperCasedMethod)) {
              return;
            }

            let body: Record<string, any> = {};
            for (const elm of el.querySelectorAll('[name]')) {
              if (elm instanceof HTMLInputElement) {
                body[elm.name] = elm.value;
              }
            }

            const res = await fetch(action, {
              method: upperCasedMethod,
              body: JSON.stringify(body),
              credentials: 'include',
              headers: {
                'Content-Type': 'application/json'
              }
            });

            if (!res.ok) {
              const message = await res.text();
              console.error('error: ' + message);
              return;
            }

            const contentType = res.headers.get('Content-Type');
            if (contentType && contentType?.length > 0) {
              location.href = '/';
            }

            if (successMessage) {
              const successElm = document.createElement('p');
              successElm.innerText = successMessage;
              el.appendChild(successElm);
              setTimeout(() => {
                successElm.remove();
              }, 2000);
            }
          });
        }

        if (handlers) {
          addHandler(el, handlers);
        }
        addBindElements(el);
      }
    }
  };
};
