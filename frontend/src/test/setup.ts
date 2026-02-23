import '@testing-library/jest-dom';

// Mock window.matchMedia for Mantine components in jsdom
window.matchMedia = window.matchMedia || ((query: string) => ({
  matches: false,
  media: query,
  onchange: null,
  addListener: () => {},
  removeListener: () => {},
  addEventListener: () => {},
  removeEventListener: () => {},
  dispatchEvent: () => false,
})) as typeof window.matchMedia;

// Mock ResizeObserver
(globalThis as any).ResizeObserver = class ResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
};

// Mock getComputedStyle to return empty for Mantine
const origGetComputedStyle = window.getComputedStyle;
window.getComputedStyle = (elt: Element, pseudoElt?: string | null) => {
  try {
    return origGetComputedStyle(elt, pseudoElt);
  } catch {
    return {} as CSSStyleDeclaration;
  }
};
