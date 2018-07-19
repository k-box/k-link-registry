<template>
  <router-view></router-view>
</template>

<script>
import fontawesome from "@fortawesome/fontawesome";
import faUser from "@fortawesome/fontawesome-free-solid/faUser";
import faCode from "@fortawesome/fontawesome-free-solid/faCode";
import faList from "@fortawesome/fontawesome-free-solid/faList";
import faSignOutAlt from "@fortawesome/fontawesome-free-solid/faSignOutAlt";

fontawesome.library.add(faUser, faCode, faList, faSignOutAlt);

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
      }, 200);
    }
  }
};
</script>

<style lang="scss">
@import url("https://fonts.googleapis.com/css?family=Montserrat:400,700&amp;subset=cyrillic-ext,latin-ext");

$family-sans-serif: "Montserrat", sans-serif;

@import "../node_modules/bulma/bulma.sass";
@import "../node_modules/noty/lib/noty.css";
@import "../node_modules/noty/lib/themes/mint.css";

html {
  background-color: #fafafa;
  color: #333;
}
</style>
