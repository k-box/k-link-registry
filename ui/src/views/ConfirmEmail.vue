<template>
  <div>
    <form @submit="submit" class="form-auth">
      <div v-if="wrong" class="notification is-warning">
        <strong>{{ $t('confirm_email.wrong_token') }}</strong>
      </div>

      <h2 class="is-size-3 has-text-centered">{{ $t('confirm_email.title') }}</h2>
      <template v-if="emailVerification.require_password">
        <input v-model="emailVerification.password" name="password" type="password" class="input is-medium is-shadowless"
        :placeholder="$t('confirm_email.password')" required>
      </template>
      <button class="button is-medium is-fullwidth is-info" type="submit">{{ $t('confirm_email.submit') }}</button>
    </form>
  </div>
</template>

<script>
import auth from "@/utils/auth";
import { mapState } from "vuex";

export default {
  name: "confirm email",
  props: ["dependencies"],
  data: function() {
    return {
      emailVerification: {
        password: ""
      }
    };
  },
  created() {
    this.fetchData();
  },
  methods: {
    fetchData() {
      api
        .getEmailVerification()
        .then(emailVerification => {
          this.emailVerification = emailVerification;
        })
        .catch(e => {});
    },
    submit() {}
  }
};
</script>