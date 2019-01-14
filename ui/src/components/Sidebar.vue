<template>
    <nav :class="{active}">
      <div>
        <router-link class="nav-link" :to="{ name: 'Applications'}" :title="$t('sidebar.applications')">
          <font-awesome-icon icon="code" />
          <span>{{ $t('sidebar.applications') }}</span>
        </router-link>
      </div>
      <div v-if="$store.state.user.role == 'ROLE_ADMIN' || $store.state.user.role == 'ROLE_OWNER'">
        <router-link class="nav-link" :to="{ name: 'Registrants'}" :title="$t('sidebar.registrants')">
          <font-awesome-icon icon="user" />
          <span>{{ $t('sidebar.registrants') }}<span>
        </router-link>
        <router-link class="nav-link" :to="{ name: 'Permissions'}" :title="$t('sidebar.permissions')">
          <font-awesome-icon icon="list" />
          <span>{{ $t('sidebar.permissions') }}</span>
        </router-link>
      </div>
    </nav>
</template>

<script>
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";

export default {
  components: {
    FontAwesomeIcon
  }
};
</script>


<style lang="scss">
nav {
  position: fixed;
  height: 100%;
  width: 15em;
  top: 4em;
  left: 0;
}

.nav-link {
  display: block;
  font-size: 1.1em;
  padding: 0.5em;
  padding-left: 1em;
  color: #333;
  text-decoration: none;
  width: 100%;
  background-color: transparent;
  border: none;
  text-align: left;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.nav-link:hover {
  background-color: rgba(0, 0, 0, 0.1);
}
</style>
