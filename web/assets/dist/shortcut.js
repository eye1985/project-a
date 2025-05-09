const addHandler = (el, handlers) => {
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
const addBindElements = (currentElement) => {
    const bindValue = currentElement.getAttribute('data-bind');
    const bindAction = currentElement.getAttribute('data-bind-action');
    if (!bindValue || !bindAction) {
        return;
    }
    const bindRules = bindValue.split(',');
    const state = {};
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
        const exprExec = (input) => {
            if (trimmedExpr.indexOf('>') !== -1) {
                return input.trim().length > parseInt(trimmedExpr[1]);
            }
            else {
                return input.trim().length < parseInt(trimmedExpr[1]);
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
            const res = [];
            for (const cid in state) {
                res.push(state[cid].exprExec(state[cid].elm.value));
            }
            res.every((num) => num)
                ? currentElement.removeAttribute(bindAction)
                : currentElement.setAttribute(bindAction, bindAction);
        });
    }
};
const createTemplateStore = () => {
    const templates = [];
    const cloneRefs = [];
    return {
        add(id, template) {
            if (!id) {
                throw new Error('id is required');
            }
            templates.push({
                id,
                template
            });
        },
        get(id) {
            return templates.filter(t => t.id === id)[0];
        },
        createClone(id) {
            const template = templates.filter(t => t.id === id)[0];
            if (!template) {
                return null;
            }
            const wrapper = document.createElement('div');
            const clone = template.template.content.cloneNode(true);
            const rand = crypto.randomUUID();
            wrapper.appendChild(clone);
            wrapper.setAttribute('data-wid', rand);
            cloneRefs.push({
                id: rand,
                element: wrapper
            });
            return wrapper;
        },
        remove(randId) {
            cloneRefs.find(clone => clone.id === randId)?.element.remove();
        }
    };
};
const attachActions = (elements, templates, handlers) => {
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
                let body = {};
                for (const elm of el.querySelectorAll('[name]')) {
                    if (elm instanceof HTMLInputElement) {
                        const isNumber = !isNaN(Number(elm.value));
                        if (isNumber) {
                            if (Number.isInteger(elm.value)) {
                                body[elm.name] = parseInt(elm.value);
                                continue;
                            }
                            body[elm.name] = parseFloat(elm.value);
                            continue;
                        }
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
        if (el instanceof HTMLTemplateElement) {
            templates.add(el.getAttribute('data-cid'), el);
        }
        if (handlers) {
            addHandler(el, handlers);
        }
        addBindElements(el);
    }
};
export const getElement = (cid) => {
    const elm = document.querySelector(`[data-cid='${cid}']`);
    if (!elm) {
        console.warn(`Not found element ${cid}`);
    }
    return elm;
};
// TODO create a delete cookie
export const getCookie = (name) => {
    const cookie = document.cookie
        .split(';')
        .map(cookie => cookie.trim())
        .find(cookie => cookie.startsWith(name + '='))
        ?.split('=');
    if (cookie && cookie.length > 1) {
        return {
            key: cookie[0],
            value: cookie[1]
        };
    }
    return null;
};
export const deleteAllCookies = () => {
    document.cookie = 'flash=; Max-Age=0; path=/;';
};
export const shortcut = () => {
    const templateStore = createTemplateStore();
    let elements = [];
    let handlerNames = [];
    let handlers = null;
    return {
        templateStore() {
            return templateStore;
        },
        addHandler(handlersArg) {
            const isAllHandlersPresent = Object.keys(handlersArg).every((handlerName) => handlerNames.includes(handlerName));
            if (!isAllHandlersPresent) {
                throw new Error('Not all handlers present, did you remember to add data-handler attribute to your elements?');
            }
            handlers = handlersArg;
        },
        appendScanElements(target) {
            let scanned = Array.from(target.querySelectorAll('[data-cid]'));
            let scannedHandlerNames = scanned.map((e) => e.getAttribute('data-handler')?.split(':')[1]);
            let appendHandlers = null;
            return {
                addHandler(handlersArg) {
                    const isAllHandlersPresent = Object.keys(handlersArg).every((handlerName) => scannedHandlerNames.includes(handlerName));
                    if (!isAllHandlersPresent) {
                        throw new Error('Not all handlers present, did you remember to add data-handler attribute to your elements?');
                    }
                    appendHandlers = handlersArg;
                    return {
                        setActions() {
                            attachActions(scanned, templateStore, appendHandlers);
                        }
                    };
                }
            };
        },
        scanElements() {
            elements = Array.from(document.querySelectorAll('[data-cid]'));
            handlerNames = elements.map((e) => e.getAttribute('data-handler')?.split(':')[1]);
        },
        setActions() {
            attachActions(elements, templateStore, handlers);
        }
    };
};
