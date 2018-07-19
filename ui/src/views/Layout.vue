<template>
    <div>
        <headerbar></headerbar>
        <sidebar></sidebar>

        <div class="main">
            <router-view></router-view>
        </div>

        <div class="modal">
            <div class="modal-background"></div>
            <div class="modal-content">
                <div class="box">
                    <h3 class="is-3">Test</h3>
                </div>
            </div>
            <button class="modal-close is-large" aria-label="close"></button>
        </div>
    </div>
</template>

<script>
import Sidebar from "@/components/Sidebar";
import Headerbar from "@/components/Headerbar";

export default {
  name: "layout",
  components: {
    Sidebar,
    Headerbar
  }
};
</script>


<style lang="scss">
.main {
  margin-left: 15em;
  width: auto;
  padding: 1em;
  margin-top: 4em;
}
</style>
