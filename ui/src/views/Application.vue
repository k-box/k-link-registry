<template>
  <div class="box">
    <h3 v-if="application.id == 0" class="title is-3">New Application</h3>
    <h3 v-if="application.id != 0" class="title is-3">Edit Application: {{ application.name }}</h3>
    <form @submit="checkForm" action="#" method="post">
      <div class="field is-horizontal">
        <label for="name" class="field-label is-normal">Name</label>
        <div class="field-body">
          <div class="field">
            <div class="control is-expanded">
              <input id="name" v-model="application.name" type="text" class="input">
              <p v-if="errors.name" class="help is-danger">
                {{errors.name}}
              </p>
            </div>
          </div>
        </div>
      </div>

      <div class="field is-horizontal">
        <label for="owner" class="field-label is-normal">Owner</label>
        <div class="field-body">
          <div class="field is-narrow">
            <div class="control">
              <div class="select">
                <select id="owner" v-model="application.owner_id">
                  <option v-for="user in registrants" :key="user.id" :value="user.id">{{ user.name }}</option>
                </select>
              </div>
                <p v-if="errors.owner_id" class="help is-danger">
                  {{errors.owner_id}}
                </p>
            </div>
          </div>
        </div>
      </div>

      <div class="field is-horizontal">
        <label for="app_domain" class="field-label is-normal">App Domain</label>
        <div class="field-body">
          <div class="field">
            <div class="control is-expanded">
              <input id="domain" v-model="application.app_domain" type="text" class="input">
              <p v-if="errors.app_domain" class="help is-danger">
                {{errors.app_domain}}
              </p>
            </div>
          </div>
        </div>
      </div>

      <div class="field is-horizontal">
        <label for="token" class="field-label is-normal">Token</label>
        <div class="field-body">
          <div class="field">
            <div class="control is-expanded">
              <input id="token" v-model="application.token" type="text" class="input" readonly>
            </div>
            <p class="help" v-if="!application.token">The token will be autogenerated on application save</p>
            <p class="help" v-if="application.token">Use the token to authenticate your application against the K-Link Network and its hosted K-Links</p>
          </div>
        </div>
      </div>

      <div class="field is-horizontal">
        <label class="field-label is-normal">Permissions</label>
        <div class="field-body">
          <div class="field is-narrow">
            <p v-if="errors.permissions" class="help is-danger">
              {{errors.permissions}}
            </p>
            <div class="control" v-for="permission in permissions" :key="permission.name">
              <label :for="permission.name" class="checkbox">
                <input :id="permission.name" :value="permission.name" v-model="application.permissions" type="checkbox">
                {{ permission.name }}
              </label>
            </div>
          </div>
        </div>
      </div>

      <div class="field is-horizontal">
        <label class="field-label is-normal">K-Links</label>
        <div class="field-body">
          <div class="field is-narrow">
            <div class="control" v-for="klink in klinks" :key="klink.name">
              <label :for="klink.id" class="checkbox">
                <input :id="klink.id" :value="klink.id" v-model="application.klinks" type="checkbox">
                {{ klink.name }}
              </label>
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
                <input id="active" v-model="application.active" type="checkbox">
                Application can authenticate
              </label>
            </div>
          </div>
        </div>
      </div>

      <div class="field is-horizontal">
        <div class="field-label">
          <!-- Left empty for spacing -->
        </div>
        <div class="field-body">
          <div class="field">
            <div class="control">
              <button type="submit" v-if="!!application.id" class="button is-primary">Update</button>
              <button type="submit" v-else class="button is-primary">Create</button>
            </div>
          </div>
        </div>
      </div>
    </form>

    <section class="section" v-if="!!application.id">
      <div class="container">
        <h1 class="title is-size-4">Danger zone</h1>
        <div>
          <button @click="deleteApplication" class="button is-danger">Delete</button>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import * as api from "@/utils/api";
import store from "@/store";

const baseApplication = {
  id: 0,
  owner_id: 0,
  name: "",
  app_domain: "",
  token: "",
  active: true,
  permissions: [],
  klinks: []
};

export default {
  name: "application",
  props: ["dependencies"],
  data: function() {
    return {
      application: {},
      permissions: [],
      registrants: [],
      klinks: [],
      errors: {}
    };
  },
  created() {
    this.fetchData();
  },
  methods: {
    checkForm(event) {
      this.errors = {
        name: null,
        permissions: null,
        app_domain: null,
        owner_id: null,
      };
      var errored = false;

      if (!this.application.name) {
        this.errors.name = 'Application Name is required.';
        errored = true;
      }
      if (!this.application.owner_id) {
        this.errors.owner_id = 'Please select the owner of this application.';
        errored = true;
      }
      if (!this.application.permissions || (this.application.permissions && this.application.permissions.length === 0)) {
        this.errors.permissions = 'Please select at least one permission to assign to the application.';
        errored = true;
      }
      if (!this.application.app_domain) {
        this.errors.app_domain = 'Please insert the application domain.';
        errored = true;
      }

      event.preventDefault();

      if(!errored){

        if(this.application.id){
          this.updateApplication(event);
        }
        else {
          this.createApplication(event);
        }
      }

    },
    createApplication(event) {
      event.preventDefault();
      event.stopPropagation();

      api
        .newApplication(this.application)
        .then(() => {
          this.$showSuccess("Application created");
          this.$router.push({ name: "Applications" });
        })
        .catch(e => {
          this.$showError("Error creating the Application");
        });
    },
    updateApplication(event) {
      event.preventDefault();
      event.stopPropagation();

      api
        .updateApplication(this.application)
        .then(application => {
          this.$showSuccess("Application updated");
          this.$router.push({ name: "Applications" });
        })
        .catch(e => {
          this.$showError("Error updating the Application");
        });
    },
    deleteApplication(event) {
      event.preventDefault();
      event.stopPropagation();

      api
        .deleteApplication(this.application.id)
        .then(() => {
          this.$showSuccess("Application deleted");
          this.$router.push({ name: "Applications" });
        })
        .catch(e => {
          this.$showError("Error deleting the Application");
          console.log(e);
          this.errors.push(e);
        });
    },
    fetchData() {
      let applicationID = this.$route.params.id;

      // if we want to create a new user, load the default user instead of
      // querying the API
      if (applicationID === "new") {
        this.application = baseApplication;
        this.application.owner_id = store.state.user.id;
      } else {
        api
          .getApplication(applicationID)
          .then(application => {
            this.application = application;
          })
          .catch(e => {
            this.$showError("Error fetching Applications");
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
        .getKlinks()
        .then(klinks => {
          this.klinks = klinks;
        })
        .catch(e => {
          this.$showError("Error fetching K-Links");
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
