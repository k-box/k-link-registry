<template>
  <div>
    <form @submit="submit" class="form-auth">
      <div v-if="wrong" class="notification is-warning">
        <strong>{{ $t('signup.invalid_request') }}</strong>
      </div>

      <h2 class="is-size-3 has-text-centered">{{ $t('signup.title') }}</h2>

      <p v-if="errors" class="notification is-danger">
        {{errors}}
      </p>

      <div v-if="$store.state.acceptNewUsers">
        <input v-model="registration.name" name="name" type="text" class="input is-medium is-shadowless"
        :placeholder="$t('signup.name')" required>
        <input v-model="registration.email" name="email" type="email" class="input is-medium is-shadowless"
        :placeholder="$t('signup.email')" required autofocus>
        <button class="button is-medium is-fullwidth is-info" :disabled="inProgress" type="submit">{{ $t('signup.submit') }}</button>
      </div>
      <div v-else>
        {{ $t('signup.registrationDisabled') }}
      </div>
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
      errors: null,
      success: false,
      inProgress: false,
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

      this.inProgress = true;

      api
        .newRegistration(this.registration)
        .then(registrant => {
          this.success = true;
          this.inProgress = false;
        })
        .catch(e => {
          
          if(e.response.data){
            this.errors = e.response.data.message;
          }
          else {
            this.errors = e.statusText;
          }

          this.inProgress = false;
        });
    }
  }
};
</script>