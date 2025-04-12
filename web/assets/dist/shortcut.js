const addHandler = (el, handlers) => {
    const eventAndHandler = el.getAttribute('data-handler');
    if (!eventAndHandler) {
        return;
    }
    const [event, handlerName] = eventAndHandler.split(':');
    if (!handlerName || !event) {
        return;
    }
    const handler = handlers[handlerName] ? handlers[handlerName] : () => { };
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
            bindAction,
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
export const shortcut = () => {
    const elements = Array.from(document.querySelectorAll('[data-cid]'));
    const handlerNames = elements.map((e) => e.getAttribute('data-handler')?.split(':')[1]);
    let handlers = null;
    return {
        addHandler(handlersArg) {
            const isAllHandlersPresent = Object.keys(handlersArg).every((handlerName) => handlerNames.includes(handlerName));
            if (!isAllHandlersPresent) {
                throw new Error('Not all handlers present');
            }
            handlers = handlersArg;
        },
        getElement(cid) {
            const elm = document.querySelector(`[data-cid='${cid}']`);
            if (!elm) {
                throw new Error(`Not found element ${cid}`);
            }
            return elm;
        },
        init() {
            if (!handlers) {
                throw new Error('No handlers present');
            }
            for (const el of elements) {
                addHandler(el, handlers);
                addBindElements(el);
            }
        },
    };
};
