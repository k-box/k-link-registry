<template>
  <div>
    <form @submit="submit" class="form-auth">
      <div v-if="wrong" class="notification is-warning">
        <strong>{{ $t('signup.invalid_request') }}</strong>
      </div>

      <h2 class="is-size-3 has-text-centered">{{ $t('signup.title') }}</h2>
      <input v-model="registration.name" name="name" type="text" class="input is-medium is-shadowless"
      :placeholder="$t('signup.name')" required>
      <input v-model="registration.email" name="email" type="email" class="input is-medium is-shadowless"
      :placeholder="$t('signup.email')" required autofocus>
      <button class="button is-medium is-fullwidth is-info" type="submit">{{ $t('signup.submit') }}</button>
    </form>
    <div class="has-text-centered">
      <router-link to="/auth/log-in" class="has-text-white">{{ $t('signup.login_link') }}</router-link>
    </div>
  </div>
</template>

<script>
import * as api from "@/utils/api";
import { mapState } from "vuex";

export default {
  name: "signup",
  props: ["dependencies"],
  data: function() {
    return {
      wrong: false,
      errors: [],
      success: false,
      registration: {
        name: "",
        email: ""
      }
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

      api
        .newRegistration(this.registration)
        .then(registrant => {
          this.success = true;
        })
        .catch(e => {
          this.errors.push(e);
        });
    }
  }
};
</script>