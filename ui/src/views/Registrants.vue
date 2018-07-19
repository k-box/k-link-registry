<template>
    <div class="box">
        <router-link :to="{name: 'Registrant', params: {id: 'new'}}" class="is-pulled-right button is-success is-medium" >{{ $t('registrants.button_new') }}</router-link>
        <h3 class="title is-3">{{ $t('registrants.title') }}</h3>
        <table class="table is-hoverable is-fullwidth">
            <thead>
                <tr>
                    <th style="width: 24px"><!-- avatar --></th>
                    <th>{{ $t('registrants.name') }}</th>
                    <th>{{ $t('registrants.active') }}</th>
                    <th>{{ $t('registrants.actions') }}</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="registrant in registrants" :key="registrant.id">
                    <td>
                      <figure class="image is-24x24">
                        <img src="@/assets/img/24x24.png">
                      </figure>
                    </td>
                    <td>
                      <router-link :to="{name: 'Registrant', params: {id: registrant.id}}">
                        {{ registrant.name }}
                        </router-link>
                      </td>
                    <td>{{ registrant.active }}</td>
                    <td>
                      <router-link :to="{name: 'Registrant', params: {id: registrant.id}}" class="button is-pulled-right is-primary is-small" >{{ $t('registrants.button_edit') }}</router-link>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<script>
import * as api from "@/utils/api";

export default {
  name: "registrants",
  props: ["dependencies"],
  data: function() {
    return {
      registrants: [],
      errors: []
    };
  },
  created() {
    this.fetchData();
  },
  methods: {
    fetchData() {
      api
        .getRegistrants()
        .then(registrants => {
          this.registrants = registrants;
        })
        .catch(e => {
          this.errors.push(e);
        });
    }
  }
};
</script>
