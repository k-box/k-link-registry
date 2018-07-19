import Vue from 'vue';
import VueI18n from 'vue-i18n';
import en from './en.yaml';
import de from './de.yaml';

Vue.use(VueI18n);

export function detectLocale () {
    let locale = (navigator.language || navigator.browserLangugae).toLowerCase();
    switch (true) {
      case /^en.*/i.test(locale):
        locale = 'en';
        break;
      case /^de.*/i.test(locale):
        locale = 'de';
        break;
      default:
        locale = 'en';
    }
  
    return locale;
  }

const i18n = new VueI18n({
    locale: detectLocale(),
    fallbackLocale: 'en',
    messages: {
      'en': en,
      'de': de
    }
  });
  
  export default i18n;