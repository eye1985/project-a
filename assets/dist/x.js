export const shortcut = () => {
    const elements = Array.from(document.querySelectorAll('[data-cid]'));
    elements.forEach((el) => {
        const eventAndHandler = el.getAttribute('data-handler');
        if (eventAndHandler) {
            const split = eventAndHandler.split(':');
            const event = split[0];
            const handlerName = split[1];
            if (handlerName && event) {
                const handler = globalThis[handlerName];
                el.addEventListener(event, handler);
            }
        }
    });
    return {
        elements,
    };
};
