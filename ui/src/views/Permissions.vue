<template>
    <div class="box">
        <h3 class="title is-3">{{ $t('permissions.title') }}</h3>
        <table class="table is-hoverable is-fullwidth">
            <thead>
                <tr>
                    <th>{{ $t('permissions.name') }}</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="permission in permissions" :key="permission.name">
                    <td>{{ permission.name }}</td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<script>
import * as api from "@/utils/api";

export default {
  name: "permissions",
  props: ["dependencies"],
  data: function() {
    return {
      permissions: [],
      errors: []
    };
  },
  created() {
    this.fetchData();
  },
  methods: {
    fetchData() {
      api
        .getPermissions()
        .then(permissions => {
          this.permissions = permissions;
        })
        .catch(e => {
          this.errors.push(e);
        });
    }
  }
};
</script>
