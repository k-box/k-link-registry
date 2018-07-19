<template>
  <div>
    <form @submit="submit" class="form-auth">
      <div v-if="wrong" class="notification is-warning">
        <strong>{{ $t('reset.invalid_request') }}</strong>
      </div>

      <h2 class="is-size-3 has-text-centered">{{ $t('reset.title') }}</h2>
      <input v-model="email" name="email" type="email" class="input is-medium is-shadowless"
      :placeholder="$t('reset.email')" required autofocus>
      <button class="button is-medium is-fullwidth is-info" type="submit">{{ $t('reset.submit') }}</button>
    </form>
    <div class="has-text-centered">
      <router-link to="/auth/log-in" class="has-text-white">{{ $t('reset.login_link') }}</router-link>
    </div>
  </div>
</template>

<script>
import auth from "@/utils/auth";
import { mapState } from "vuex";

export default {
  name: "signup",
  props: ["dependencies"],
  data: function() {
    return {
      wrong: false,
      email: ""
    };
  },
  mounted() {
    if (this.dependencies) this.setup();
  },
  watch: {
    dependencies: function(val) {
      if (val) this.setup();
    }
  },
  methods: {
    submit(event) {
      event.preventDefault();
      event.stopPropagation();
    }
  }
};
</script>
