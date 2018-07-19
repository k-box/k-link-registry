import Vue from "vue";
import Vuex from "vuex";
import mutations from "./mutations";
import getters from "./getters";

Vue.use(Vuex);

// shared state
const state = {
  baseURL: document
    .querySelector('meta[name="base"]')
    .getAttribute("content")
    .replace(/\/$/, ''), // remove trailing slash
  jwt: '',
  user: {
    id: 0,
    role: '',
  }
};

export default new Vuex.Store({
  strict: process.env.NODE_ENV !== "production",
  state,
  getters,
  mutations
});