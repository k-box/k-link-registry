<template>
    <div class="box">
      <router-link :to="{name: 'Application', params: {id: 'new'}}" class="is-pulled-right button is-success is-medium" >{{ $t('applications.button_new') }}</router-link>
        <h3 class="title is-3">{{ $t('applications.title') }}</h3>
        <table class="table is-hoverable is-fullwidth">
            <thead>
                <tr>
                    <th>Name</th>
                    <th>Domain</th>
                    <th>Active</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="app in applications" :key="app.id">
                    <td><router-link :to="{name: 'Application', params: {id: app.id}}" >{{ app.name }}</router-link></td>
                    <td>{{ app.app_domain }}</td>
                    <td>{{ app.active }}</td>
                    <td>
                      <router-link :to="{name: 'Application', params: {id: app.id}}" class="button is-pulled-right is-primary is-small" >{{ $t('applications.button_edit') }}</router-link>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<script>
import * as api from "@/utils/api";

export default {
  name: "applications",
  props: ["dependencies"],
  data: function() {
    return {
      applications: [],
      errors: []
    };
  },
  created() {
    this.fetchData();
  },
  methods: {
    fetchData() {
      api
        .getApplications()
        .then(applications => {
          this.applications = applications;
        })
        .catch(error => {
          this.errors.push(error);
        });
    }
  }
};
</script>