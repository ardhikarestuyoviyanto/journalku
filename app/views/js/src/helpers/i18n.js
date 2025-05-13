import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import LanguageDetector from 'i18next-browser-languagedetector'

import translationEN from '../locale/en.json'
import translationID from '../locale/id.json'

const resources = {
  en: { translation: translationEN },
  id: { translation: translationID },
}

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    fallbackLng: 'id', // default bahasa ke Indonesia
    detection: {
      order: ['localStorage', 'navigator'], // opsional
    },
    interpolation: {
      escapeValue: false,
    },
  })

export default i18n
