<template>
    <div class="box">
      <h3 v-if="!registrant.id" class="title is-3">New Registrant</h3>
      <h3 v-if="!!registrant.id" class="title is-3">Edit Registrant: {{ registrant.name }}</h3>
      <form>
        <div class="field is-horizontal">
          <label for="name" class="field-label is-normal">Name</label>
          <div class="field-body">
            <div class="field">
              <div class="control is-expanded">
                <input id="name" v-model="registrant.name" type="text" class="input">
              </div>
            </div>
          </div>
        </div>

        <div class="field is-horizontal">
          <label for="email" class="field-label is-normal">Email</label>
          <div class="field-body">
            <div class="field">
              <div class="control is-expanded">
                <input id="email" v-model="registrant.email" type="text" class="input">
              </div>
            </div>
          </div>
        </div>

        <div v-if="!!registrant.id" class="field is-horizontal">
          <label for="password" class="field-label is-normal">Password</label>
          <div class="field-body">
            <div class="field">
              <div class="control is-expanded">
                <input id="password" placeholder="leave empty to keep" v-model="registrant.password" type="password" class="input">
              </div>
            </div>
          </div>
        </div>

        <div class="field is-horizontal">
          <label for="role" class="field-label is-normal">Role</label>
          <div class="field-body">
            <div class="field is-narrow">
              <div class="control">
                <div class="select is-primary">
                  <select id="role" v-model="registrant.role">
                    <option value="ROLE_USER">User</option>
                    <option value="ROLE_ADMIN">Admin</option>
                    <option value="ROLE_OWNER">Owner</option>
                  </select>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="field is-horizontal">
          <label for="active" class="field-label is-normal">Active</label>
          <div class="field-body">
            <div class="field is-narrow">
              <div class="control">
                <label class="checkbox">
                  <input id="active" v-model="registrant.active" type="checkbox">
                  User can log in
                </label>
              </div>
            </div>
          </div>
        </div>

        <template v-if="!!registrant.id">
          <button @click="deleteRegistrant" class="is-pulled-right button is-danger">Delete</button>
          <button @click="updateRegistrant" class="button is-primary">Update</button>
        </template>
        <template v-else>
          <button @click="createRegistrant" class="button is-primary">Create</button>
        </template>
        


      </form>
    </div>
</template>

<script>
import * as api from "@/utils/api";

const baseRegistrant = {
  id: 0,
  email: "",
  name: "",
  role: "ROLE_USER",
  active: true
};

export default {
  name: "registrant",
  props: ["dependencies"],
  data: function() {
    return {
      registrant: {},
      errors: []
    };
  },
  created() {
    this.fetchData();
  },
  methods: {
    createRegistrant(event) {
      event.preventDefault();
      event.stopPropagation();

      api
        .newRegistrant(this.registrant)
        .then(() => {
          this.$showSuccess("Registrant created");
          this.$router.push({ name: "Registrants" });
        })
        .catch(e => {
          this.$showError("Error creating Registrant");
        });
    },
    updateRegistrant(event) {
      event.preventDefault();
      event.stopPropagation();

      api
        .updateRegistrant(this.registrant)
        .then(registrant => {
          this.$showSuccess("Registrant updated");
          this.$router.push({ name: "Registrants" });
        })
        .catch(e => {
          this.$showError("Error updating the Registrant");
        });
    },
    deleteRegistrant(event) {
      event.preventDefault();
      event.stopPropagation();

      api
        .deleteRegistrant(this.registrant.id)
        .then(() => {
          this.$showSuccess("Registrant deleted");
          this.$router.push({ name: "Registrants" });
          //console.log(this.$router);
        })
        .catch(e => {
          this.$showError("Error deleting the Registrant");
          console.log(e);
          this.errors.push(e);
        });
    },
    fetchData() {
      let registrantID = this.$route.params.id;

      // if we want to create a new user, load the default user instead of
      // querying the API
      if (registrantID === "new") {
        this.registrant = baseRegistrant;
        return;
      }

      api
        .getRegistrant(registrantID)
        .then(registrant => {
          this.registrant = registrant;
        })
        .catch(e => {
          console.log(e);
          this.errors.push(e);
        });
    }
  }
};
</script>
