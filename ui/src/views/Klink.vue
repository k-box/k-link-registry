<template>
  <div class="box">
    <h3 v-if="klink.id == 0" class="title is-3">New K-Link</h3>
    <h3 v-if="klink.id != 0" class="title is-3">Edit K-Link: {{ klink.name }}</h3>
    <form>
      <div class="field is-horizontal">
        <label for="name" class="field-label is-normal">Name</label>
        <div class="field-body">
          <div class="field">
            <div class="control is-expanded">
              <input id="name" v-model="klink.name" type="text" class="input">
            </div>
          </div>
        </div>
      </div>

      <div class="field is-horizontal">
        <label for="owner" class="field-label is-normal">Manager</label>
        <div class="field-body">
          <div class="field is-narrow">
            <div class="control">
              <div class="select is-primary">
                <select id="owner" v-model="klink.manager_id">
                  <option v-for="user in registrants" :key="user.id" :value="user.id">{{ user.name | user.email }}</option>
                </select>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="field is-horizontal">
        <label for="app_domain" class="field-label is-normal">Website</label>
        <div class="field-body">
          <div class="field">
            <div class="control is-expanded">
              <input id="domain" v-model="klink.website" type="text" class="input">
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
                <input id="active" v-model="klink.active" type="checkbox">
                Make this K-Link active and accepting publications
              </label>
            </div>
          </div>
        </div>
      </div>

      <template v-if="!!klink.id">
          <button @click="deleteKlink" class="is-pulled-right button is-danger">Delete</button>
          <button @click="updateKlink" type="submit" class="button is-primary">Update</button>
        </template>
        <template v-else>
          <button @click="createKlink" type="submit" class="button is-primary">Create</button>
        </template>

        <div class="field is-horizontal" v-if="klink.identifier !== ''">
        <label for="identifier" class="field-label is-normal">Identifier</label>
        <div class="field-body">
          <div class="field">
            <div class="control is-expanded">
              <input id="identifier" placeholder="" v-model="klink.identifier" type="text" class="input" readonly>
            </div>
          </div>
        </div>
      </div>
    </form>
  </div>
</template>

<script>
import * as api from "@/utils/api";
import store from "@/store";

const baseKlink = {
  id: null,
  manager_id: null,
  name: "",
  website: "",
  description: "",
  active: false,
};

export default {
  name: "klink",
  props: ["dependencies"],
  data: function() {
    return {
      klink: {},
      permissions: [],
      registrants: [],
      errors: []
    };
  },
  created() {
    this.fetchData();
  },
  methods: {
    createKlink(event) {
      event.preventDefault();
      event.stopPropagation();

      api
        .newKlink(this.klink)
        .then(() => {
          this.$showSuccess("K-Link created");
          this.$router.push({ name: "Klinks" });
        })
        .catch(e => {
          this.$showError("Error creating the K-Link");
        });
    },
    updateKlink(event) {
      event.preventDefault();
      event.stopPropagation();

      api
        .updateKlink(this.klink)
        .then(klink => {
          this.$showSuccess("K-Link updated");
          this.$router.push({ name: "Klinks" });
        })
        .catch(e => {
          this.$showError("Error updating the K-Link");
        });
    },
    deleteKlink(event) {
      event.preventDefault();
      event.stopPropagation();

      api
        .deleteKlink(this.klink.id)
        .then(() => {
          this.$showSuccess("K-Link deleted");
          this.$router.push({ name: "Klinks" });
        })
        .catch(e => {
          this.$showError("Error deleting the K-Link");
          console.log(e);
          this.errors.push(e);
        });
    },
    fetchData() {
      let klinkID = this.$route.params.id;

      // if we want to create a new user, load the default user instead of
      // querying the API
      if (klinkID === "new") {
        this.klink = baseKlink;
        this.klink.manager_id = store.state.user.id;
      } else {
        api
          .getKlink(klinkID)
          .then(klink => {
            this.klink = klink;
          })
          .catch(e => {
            this.$showError("Error fetching K-Link");
            this.errors.push(e);
          });
      }

      api
        .getPermissions()
        .then(permissions => {
          this.permissions = permissions;
        })
        .catch(e => {
          this.$showError("Error fetching Permissions");
          this.errors.push(e);
        });

      api
        .getRegistrants()
        .then(registrants => {
          this.registrants = registrants;
        })
        .catch(e => {
          this.$showError("Error fetching Registrants");
          this.errors.push(e);
        });
    }
  }
};
</script>
