import en from './i18n/en';
import zh from './i18n/zh';

const LOCALE_KEY = 'gtm_locale';

export type Locale = 'en' | 'zh';

const messages: Record<Locale, Record<string, string>> = { en, zh };

function createI18n() {
  let locale = $state<Locale>(
    (typeof localStorage !== 'undefined'
      ? (localStorage.getItem(LOCALE_KEY) as Locale)
      : null) || 'zh'
  );

  return {
    get locale() {
      return locale;
    },
    set locale(v: Locale) {
      locale = v;
      if (typeof localStorage !== 'undefined') {
        localStorage.setItem(LOCALE_KEY, v);
      }
    },
    t(key: string, params?: Record<string, string | number>): string {
      let msg = messages[locale]?.[key] ?? messages.en[key] ?? key;
      if (params) {
        for (const [k, v] of Object.entries(params)) {
          msg = msg.replaceAll(`{${k}}`, String(v));
        }
      }
      return msg;
    },
  };
}

export const i18n = createI18n();
export const t = (key: string, params?: Record<string, string | number>) => i18n.t(key, params);
