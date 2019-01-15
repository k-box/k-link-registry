<template>
  <router-view></router-view>
</template>

<script>
import { library } from "@fortawesome/fontawesome-svg-core";
import { fas } from '@fortawesome/free-solid-svg-icons';
import { faUser } from "@fortawesome/free-solid-svg-icons/faUser";
import { faCode } from "@fortawesome/free-solid-svg-icons/faCode";
import { faList } from "@fortawesome/free-solid-svg-icons/faList";
import { faSignOutAlt } from "@fortawesome/free-solid-svg-icons/faSignOutAlt";

library.add(faUser, faCode, faList, faSignOutAlt);

export default {
  name: "app",
  data() {
    return {
      loaded: false
    };
  },
  mounted() {
    // once everything is initialized, remove
    // loading animation.
    this.unload();
  },
  methods: {
    unload() {
      this.loaded = true;
      // Remove loading animation.
      let loading = document.getElementById("loading");
      loading.classList.add("done");
      // remove after fade out
      setTimeout(function() {
        loading.parentNode.removeChild(loading);
      }, 10);
    }
  }
};
</script>

<style lang="scss">

$family-sans-serif: system-ui,BlinkMacSystemFont,-apple-system,Segoe UI,Roboto,Oxygen,Ubuntu,Cantarell,Fira Sans,Droid Sans,Helvetica Neue,sans-serif;

@import "../node_modules/bulma/bulma.sass";
@import "../node_modules/noty/lib/noty.css";
@import "../node_modules/noty/lib/themes/mint.css";

html {
  background-color: #fafafa;
}
</style>
