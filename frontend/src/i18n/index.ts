import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

import en from './locales/en.json';
import uz from './locales/uz.json';
import ru from './locales/ru.json';

const DEFAULT_LANG = 'en';
const STORAGE_KEY = 'i18nextLng';

const savedLang = localStorage.getItem(STORAGE_KEY) ?? DEFAULT_LANG;

i18n.use(initReactI18next).init({
  lng: savedLang,
  fallbackLng: DEFAULT_LANG,
  supportedLngs: ['uz', 'ru', 'en'],
  nonExplicitSupportedLngs: true,
  interpolation: { escapeValue: false },
  resources: {
    uz: { translation: uz },
    ru: { translation: ru },
    en: { translation: en },
  },
});

i18n.on('languageChanged', (lng) => {
  localStorage.setItem(STORAGE_KEY, lng);
});

export default i18n;
