import Vue from "vue";
import router from "./router";
import store from "./store";
import App from "./App";
import i18n from "./i18n";
import Noty from 'noty';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';

const defaultNoty = {
  layout: "topRight",
  type: "info",
  timeout: 5000,
  progressBar: true
};

Vue.prototype.$noty = function (opts) {
  new Noty(Object.assign({}, defaultNoty, opts)).show();
};

Vue.prototype.$showSuccess = function (message) {
  new Noty(Object.assign({}, defaultNoty, {
    text: message,
    type: 'success'
  })).show();
};

Vue.prototype.$showError = function (error) {
  let n = new Noty(Object.assign({}, defaultNoty, {
    text: error,
    type: 'error',
    timeout: null,
    buttons: [
      Noty.button(i18n.t('buttons.reportIssue'), '', function () {
        window.open('https://github.com/k-box/k-link-registry/issues');
      }),
      Noty.button(i18n.t('buttons.close'), '', function () {
        n.close();
      })
    ]
  }));

  n.show();
};

Vue.component('font-awesome-icon', FontAwesomeIcon);

console.log(
  "%cThank you for using the K-Link-Registry! %c😊",
  // Nice big comic-sans-like font, because comic-sans
  // is never wrong.
  "font: 2em cursive; color: #dd4814;",
  // we use serif for the emoji, since it is more
  // likely to contain the graphical variant.
  "font: 3em serif;"
);

new Vue({
  el: "#app",
  store,
  router,
  i18n,
  template: "<App/>",
  components: {
    App
  }
});