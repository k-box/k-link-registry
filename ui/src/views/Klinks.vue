<template>
    <div class="box">
      <router-link :to="{name: 'Klink', params: {id: 'new'}}" class="is-pulled-right button is-success is-medium" >{{ $t('klinks.button_new') }}</router-link>
        <h3 class="title is-3">{{ $t('klinks.title') }}</h3>
        <table class="table is-hoverable is-fullwidth">
            <thead>
                <tr>
                    <th>Name</th>
                    <th>Website</th>
                    <th>Active</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="app in klinks" :key="app.id">
                    <td><router-link :to="{name: 'Klink', params: {id: app.id}}" >{{ app.name }}</router-link><br/><span>{{ app.id }}</span></td>
                    <td>{{ app.website }}</td>
                    <td>{{ app.active }}</td>
                    <td>
                      <router-link :to="{name: 'Klink', params: {id: app.id}}" class="button is-pulled-right is-primary is-small" >{{ $t('klinks.button_edit') }}</router-link>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<script>
import * as api from "@/utils/api";

export default {
  name: "klinks",
  props: ["dependencies"],
  data: function() {
    return {
      klinks: [],
      errors: []
    };
  },
  created() {
    this.fetchData();
  },
  methods: {
    fetchData() {
      api
        .getKlinks()
        .then(klinks => {
          this.klinks = klinks;
        })
        .catch(error => {
          this.errors.push(error);
        });
    }
  }
};
</script>