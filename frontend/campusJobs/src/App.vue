<template>
  <div>
    <h2>hello world</h2>
    <ul>
      <li v-for="(id, index) in data" :key="index">{{ id }}</li>
    </ul>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import {fetchData} from './axios'

export default defineComponent({
  name: 'App',
  data() {
    return {
      data: []
    };
  },
  mounted() {
    fetchData()
    .then((responseData: any) => {
      this.data = responseData.data
    })
    .catch((error: any) => {
      console.log(error)
    });
  }
});
</script>