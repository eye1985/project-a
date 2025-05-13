const PREFIX = 'sc';

// data-* actions
const ID = 'id';
const TEMPLATE_ID = 'template-id';
const TYPE = 'type';
const METHOD = 'method';
const FORM_ON_ERROR = 'form-onerror';
const FORM_ON_SUCCESS = 'form-onsuccess';
const HANDLER = 'handler';

type Handle = (
  e: Event,
  currentCustomElement: CustomElement,
  is: InternalState,
) => void;
type FormMethod = (data: any) => void;
type InternalState = {
  formMethods: Map<string, FormMethod>;
  elements: Map<string, CustomElement>;
  handlers: Map<string, Handle>;
  state: Map<string, any>;
};
export const store: InternalState = {
  formMethods: new Map(),
  elements: new Map(),
  handlers: new Map(),
  state: new Map(),
};

export const getElementsByType = (type: string) => {
  const result: CustomElement[] = [];
  store.elements.forEach((elm) => {
    elm.type === type && result.push(elm);
  });

  return result;
};

const createDataName = (action: string) => {
  return `data-${PREFIX}-${action}`;
};

const scanElements = (target: Element) => {
  return Array.from(target.querySelectorAll(`[${createDataName(ID)}]`));
};

const addToInternalState = (el: Element) => {
  const id = el.getAttribute(createDataName(ID));
  if (!id) {
    throw new Error(`${addToInternalState.name}: id is required`);
  }
  if (store.elements.get(id)) {
    throw new Error(`${addToInternalState.name}: id ${id} already exists`);
  }
  store.elements.set(id, new CustomElement(el));
};

export const syncHandle = () => {
  for (const key of store.handlers.keys()) {
    store.elements.forEach((elm) => {
      if (elm.handleName === key && !elm.isHandleApplied && elm.handleEvent) {
        const handle = store.handlers.get(key);
        if (!handle) {
          throw new Error(`${syncHandle.name}: handler ${key} not found`);
        }
        const newHandle = (e: Event) => {
          handle(e, elm, store);
        };
        elm.ref.addEventListener(elm.handleEvent, newHandle);
        elm.handlers.push(newHandle);
        elm.isHandleApplied = true;
      }
    });
  }
};

export const state = {
  get(id: string) {
    return store.state.get(id);
  },
  set(id: string, value: any) {
    store.state.set(id, value);
  },
};

export class CustomElement {
  id: string;
  type: string | null = null;
  handleName: string | null = null;
  handleEvent: string | null = null;
  isHandleApplied = false;
  isTemplate = false;
  ref: Element | HTMLTemplateElement;
  templateWrapperRef: Element | null = null;
  handlers: EventListenerOrEventListenerObject[];

  constructor(el: Element | HTMLTemplateElement) {
    const id = el.getAttribute(createDataName(ID));
    if (!id) {
      throw new Error(`${CustomElement.name}: id is required`);
    }

    const type = el.getAttribute(createDataName(TYPE));
    if (type) {
      this.type = type;
    }

    const handler = el.getAttribute(`${createDataName(HANDLER)}`);

    this.id = id;
    this.ref = el;
    this.handlers = [];
    if (handler) {
      const handleNameAndEvent = handler.split(':');
      if (handleNameAndEvent[0] && handleNameAndEvent[1]) {
        this.handleEvent = handleNameAndEvent[0];
        this.handleName = handleNameAndEvent[1];
      }
    }

    const hasAllAttrs =
      this.ref.getAttribute(`${createDataName(METHOD)}`) &&
      this.ref.getAttribute('action');

    if (this.ref instanceof HTMLFormElement && hasAllAttrs) {
      this.overrideSubmit(this.ref);
    } else if (this.ref instanceof HTMLTemplateElement) {
      this.isTemplate = true;
    }
  }

  private overrideSubmit = (el: HTMLFormElement) => {
    const handler = async (e: SubmitEvent) => {
      e.preventDefault();

      const action = el.getAttribute('action');
      const method = el.getAttribute(`${createDataName(METHOD)}`);
      const formOnSuccess = el.getAttribute(
        `${createDataName(FORM_ON_SUCCESS)}`,
      );
      const formOnError = el.getAttribute(`${createDataName(FORM_ON_ERROR)}`);

      if (!method || !action) {
        console.warn(
          `${this.overrideSubmit.name}: action or method is missing, skipping`,
        );
        return;
      }

      const methods = ['POST', 'PUT', 'DELETE', 'PATCH'];
      const upperCasedMethod = method.toUpperCase();
      if (!methods.includes(upperCasedMethod)) {
        console.warn(
          `${this.overrideSubmit.name}: method ${method} is not supported, skipping`,
        );
        return;
      }

      let body: Record<string, any> = {};
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

      try {
        const res = await fetch(action, {
          method: upperCasedMethod,
          body: JSON.stringify(body),
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json',
          },
        });

        if (!res.ok) {
          const message = await res.text();
          if (!formOnError) {
            console.error('error: ' + message);
            return;
          }

          const errorFn = store.formMethods.get(formOnError);
          if (!errorFn) {
            return;
          }
          errorFn(message);
          return;
        }

        if (formOnSuccess) {
          const isJson = res.headers
            .get('content-type')
            ?.includes('application/json');
          const json = isJson ? await res.json() : null;
          const successFn = store.formMethods.get(formOnSuccess);
          if (!successFn) {
            throw new Error(
              `${this.overrideSubmit.name}: success handler ${formOnSuccess} not found`,
            );
          }
          successFn(json);
        }
      } catch (e) {
        const error = e as Error;
        console.error(error.message);
      }
    };
    el.addEventListener('submit', handler);
    this.handlers.push(handler as unknown as EventListener);
  };

  insertTemplateInto = (
    target: CustomElement | Element,
    options?: {
      clearBeforeInsert?: boolean;
      classNames?: string[];
    },
  ) => {
    const intoTarget = target instanceof CustomElement ? target.ref : target;

    if (!this.isTemplate) {
      throw new Error(
        `${this.insertTemplateInto.name}: element is not a template`,
      );
    }
    if (!intoTarget) {
      throw new Error(`${this.insertTemplateInto.name}: target is required`);
    }

    const wrapper = document.createElement('div');
    const template = this.ref as HTMLTemplateElement;
    const clone = template.content.cloneNode(true) as DocumentFragment;
    wrapper.appendChild(clone);
    wrapper.setAttribute(createDataName(TEMPLATE_ID), this.id);
    if (options && options.classNames) {
      options.classNames.forEach((className) => {
        wrapper.classList.add(className);
      });
    }
    this.templateWrapperRef = wrapper;

    if (options && options.clearBeforeInsert) {
      // TODO test this for normal Element, this might throw error
      deleteAllFromTarget(intoTarget);
      intoTarget.innerHTML = '';
    }
    intoTarget.appendChild(wrapper);
    addFromTarget(wrapper);
    syncHandle();
  };

  remove() {
    const ref = this.isTemplate ? this.templateWrapperRef : this.ref;
    if (ref) {
      this.handlers.forEach((handle) => {
        if (this.handleEvent) {
          ref.removeEventListener(this.handleEvent, handle);
        }
      });
      deleteAllFromTarget(ref);
      ref.remove();
    }
  }
}

export const addFromTarget = (target: Element) => {
  const elements = scanElements(target);
  elements.forEach((element) => {
    addToInternalState(element);
  });
};

export const deleteAllFromTarget = (target: Element) => {
  const elements = scanElements(target);
  elements.forEach((element) => {
    const id = element.getAttribute(createDataName(ID));
    if (!id) {
      throw new Error(
        `${deleteAllFromTarget.name}: id is not found, you somehow added a non custom element into the store`,
      );
    }
    const elm = store.elements.get(id);
    if (elm) {
      elm.handlers.forEach((handle) => {
        if (elm.handleEvent) {
          elm.ref.removeEventListener(elm.handleEvent, handle);
        }
      });
      elm.ref.remove();
    }

    store.elements.delete(id);
  });
};

export const isTemplate = (
  el: Element | HTMLTemplateElement,
): el is HTMLTemplateElement => {
  return el instanceof HTMLTemplateElement;
};

export const getElement = (id: string) => {
  const elm = store.elements.get(id);
  if (elm) {
    return elm;
  }

  return null;
};

export const addHandler = (handlerName: string, handle: Handle) => {
  const getHandle = store.handlers.get(handlerName);
  if (getHandle) {
    throw new Error(
      `${addHandler.name}: handler ${handlerName} already exists`,
    );
  }
  store.handlers.set(handlerName, handle);
  syncHandle();
};

export const addFormMethod = (methodName: string, method: FormMethod) => {
  const getMethod = store.formMethods.get(methodName);
  if (getMethod) {
    throw new Error(
      `${addFormMethod.name}: method ${methodName} already exists`,
    );
  }
  store.formMethods.set(methodName, method);
};

export const quickCreateNode = (html: string) => {
  const node = document.createElement('div');
  node.insertAdjacentHTML('beforeend', html);
  return node;
};

export const getCookie = (name: string) => {
  const cookie = document.cookie
    .split(';')
    .map((cookie) => cookie.trim())
    .find((cookie) => cookie.startsWith(name + '='))
    ?.split('=');

  if (cookie && cookie.length > 1) {
    return {
      key: cookie[0],
      value: cookie[1],
    };
  }

  return null;
};

export const deleteAllCookies = () => {
  document.cookie = 'flash=; Max-Age=0; path=/;';
};
