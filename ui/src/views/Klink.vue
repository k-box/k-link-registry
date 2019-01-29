<template>
  <div class="box">
    <h3 v-if="!klink.id" class="title is-3">New K-Link</h3>
    <h3 v-if="!!klink.id" class="title is-3">Edit K-Link: {{ klink.name }}</h3>
    <form @submit="checkForm" action="#" method="post">
      <div class="field is-horizontal">
        <label for="name" class="field-label is-normal">Name</label>
        <div class="field-body">
          <div class="field">
            <div class="control is-expanded">
              <input id="name" v-model="klink.name" type="text" class="input" v-bind:class="{ 'is-danger': errors.name}">
            </div>
            <p v-if="errors.name" class="help is-danger">
              {{errors.name}}
            </p>
          </div>
        </div>
      </div>

      <div class="field is-horizontal">
        <label for="manager" class="field-label is-normal">Manager</label>
        <div class="field-body">
          <div class="field is-narrow">
            <div class="control">
              <div class="select" v-bind:class="{ 'is-danger': errors.manager_id}">
                <select id="manager" v-model="klink.manager_id">
                  <option v-for="user in registrants" :key="user.id" :value="user.id">{{ user.name }}</option>
                </select>
              </div>
              <p v-if="errors.manager_id" class="help is-danger">
                {{errors.manager_id}}
              </p>
            </div>
          </div>
        </div>
      </div>

      <div class="field is-horizontal">
        <label for="website" class="field-label is-normal">Website</label>
        <div class="field-body">
          <div class="field">
            <div class="control is-expanded">
              <input id="website" v-model="klink.website" type="text" class="input">
            </div>
          </div>
        </div>
      </div>
      
      <div class="field is-horizontal">
        <label for="description" class="field-label is-normal">Description</label>
        <div class="field-body">
          <div class="field">
            <div class="control is-expanded">
              <textarea class="textarea" id="description" v-model="klink.description"></textarea>
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

      <div class="field is-horizontal">
        <div class="field-label">
          <!-- Left empty for spacing -->
        </div>
        <div class="field-body">
          <div class="field">
            <div class="control">
              <button type="submit" v-if="!!klink.id" class="button is-primary">Update</button>
              <button type="submit" v-else class="button is-primary">Create</button>
            </div>
          </div>
        </div>
      </div>

          

      <div class="field is-horizontal" v-if="!!klink.id">
        <label for="identifier" class="field-label is-normal">Identifier</label>
        <div class="field-body">
          <div class="field">
            {{klink.id}}
          </div>
        </div>
      </div>
    </form>

    <section class="section" v-if="!!klink.id">
      <div class="container">
        <h1 class="title is-size-4">Danger zone</h1>
        <div>
          <button @click="deleteKlink" class="button is-danger">Delete</button>
        </div>
      </div>
    </section>
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
      registrants: [],
      errors: []
    };
  },
  created() {
    this.fetchData();
  },
  methods: {
    checkForm(event) {

      this.errors = {
        name: null,
        manager_id: null,
      };
      var errored = false;

      if (!this.klink.name) {
        this.errors.name = 'K-Link Name required.';
        errored = true;
      }
      if (!this.klink.manager_id) {
        this.errors.manager_id = 'Please select a manager.';
        errored = true;
      }

      event.preventDefault();

      if(!errored){

        if(this.klink.id){
          this.updateKlink(event);
        }
        else {
          this.createKlink(event);
        }
      }

    },
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
