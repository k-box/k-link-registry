<template>
  <div>
    <form @submit="submit" class="form-auth">
      <div v-if="wrong" class="notification is-warning">
        <strong>{{ $t('login.wrong_credentials') }}</strong>
      </div>

      <h2 class="is-size-3 has-text-centered">{{ $t('login.title') }}</h2>
      <input v-model="email" name="email" type="email" class="input is-medium is-shadowless"
      :placeholder="$t('login.email')" required autofocus>
      <input v-model="password" name="password" type="password" class="input is-medium is-shadowless"
      :placeholder="$t('login.password')" required>
      <button class="button is-medium is-fullwidth is-info" type="submit">{{ $t('login.submit') }}</button>
      <div class="has-text-centered">
        <router-link to="/auth/reset-password" class="button is-text is-fullwidth">{{ $t('login.forgot_password_link') }}</router-link>
      </div>
    </form>
    <div class="has-text-centered" v-if="$store.state.acceptNewUsers">
      <router-link to="/auth/sign-up" class="has-text-white">{{ $t('login.sign_up_link') }}</router-link>
    </div>
  </div>
</template>

<script>
import auth from "@/utils/auth";
import { mapState } from "vuex";

export default {
  name: "login",
  props: ["dependencies"],
  data: function() {
    return {
      wrong: false,
      email: "",
      password: "",
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

      let redirect = this.$route.query.redirect;
      if (redirect === "" || redirect === undefined || redirect === null) {
        redirect = "/applications";
      }

      auth
        .login(this.email, this.password)
        .then(() => {
          this.$showSuccess("Welcome back!");
          this.$router.push({ path: redirect });
        })
        .catch(e => {
          this.wrong = true;
          console.log(e);
        });
    }
  }
};
</script>