<template>
  <div class="navbar is-info is-fixed-top header">
    <div class="navbar-brand">
      <router-link class="navbar-item" :to="{ name: 'Applications'}" :title="$t('header.title')">
        <h2 class="title is-5 has-text-white">{{ $t('header.title') }}</h2>
      </router-link>
    </div>
    <div class="navbar-menu">
      <div class="navbar-end">
        <div class="navbar-item">
          <button @click="logout" class="button is-info" id="logout" :title="$t('header.logout')">
            <font-awesome-icon icon="sign-out-alt" />
            <span>{{ $t('header.logout') }}</span>
          </button>
        </div>

      </div>
    </div>
  </div>
</template>

<script>
import auth from "@/utils/auth";
import FontAwesomeIcon from "@fortawesome/vue-fontawesome";

export default {
  methods: {
    logout: auth.logout
  },
  components: {
    FontAwesomeIcon
  }
};
</script>

<style lang="scss">
.header {
  height: 4em;
  box-shadow: 0 0 10px #0006;
}
</style>
