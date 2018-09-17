<template>
  <div>
    <form @submit="submit" class="form-auth">
      <div v-if="wrong" class="notification is-warning">
        <strong>{{ $t('email_verification.wrong_token') }}</strong>
      </div>

      <h2 class="is-size-3 has-text-centered">{{ $t('email_verification.title') }}</h2>
      <template v-if="email_verification.require_password || true">
        <p>{{ $t('email_verification.password_text') }}</p>
        <input v-model="confirmation.password" name="password" type="password" class="input is-medium is-shadowless"
        :placeholder="$t('email_verification.password')" required>
        <button @click="confirm" class="button is-medium is-fullwidth is-info" type="submit">{{ $t('email_verification.set_password_button') }}</button>
      </template>
      <template v-else>
        <p>{{ $t('email_verification.email_text') }}</p>
      <button @click="confirm" class="button is-medium is-fullwidth is-info" type="submit">{{ $t('email_verification.verify_email_button') }}</button>
      </template>
    </form>
  </div>
</template>

<script>
import * as api from "@/utils/api";
import auth from "@/utils/auth";
import { mapState } from "vuex";

export default {
  name: "confirm token",
  props: ["dependencies"],
  data: function() {
    return {
      email_verification: {
        require_password: true,
        display_name: ""
      },
      confirmation: {
        password: "",
        token: ""
      }
    };
  },
  created() {
    this.fetchData();
  },
  methods: {
    fetchData() {
      api
        .getEmailVerification(this.$route.params.token)
        .then(email_verification => {
          this.email_verification = email_verification;
        })
        .catch(e => {
          console.log(e);
        });
    },
    confirm(event) {
      event.preventDefault();
      event.stopPropagation();

      this.confirmation.token = this.$route.params.token;

      api
        .performEmailVerification(this.confirmation)
        .then(() => {
          this.$showSuccess("Confirmation successful");
          this.$router.push({ name: "Log in" });
        })
        .catch(e => {
          this.$showError("Confirmation is expired");
          this.$router.push({ name: "Log in" });
          console.log(e);
          this.errors.push(e);
        });
    }
  }
};
</script>