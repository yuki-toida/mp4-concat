<template>
  <div>
    <div v-if="uid">
      <div>
        <a v-bind:href="link" download="output.mp4">Download</a>
      </div>
      <div>
        <button v-on:click="init">Init</button>
      </div>
    </div>
    <div v-else>
      <div>
        <div v-for="key in inputs" v-bind:key="key">
          <input type="file" accept="video/*" v-on:change="change($event, key)">
        </div>
      </div>
      <div v-if="files.length == inputs">
        <button v-on:click="add">Add</button>
        <button v-if="1 < files.length" v-on:click="remove">Remove</button>
      </div>
      <div>
        <button v-on:click="upload">Concat</button>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

const domain = 'http://localhost:8080';

export default {
  data: function() {
    return {
      uid: null,
      inputs: 1,
      files: [],
    }
  },
  created: function() {
    axios.get(`${domain}/cookie`)
      .then((res) => {
        if (res.data) {
          this.uid = res.data;
        }
      });
  },
  computed: {
    link: function() {
      return `${domain}/static/video/${this.uid}.mp4`;
    }
  },
  methods: {
    init: function() {
      axios.delete(`${domain}/cookie`)
        .then(_ => this.uid = null);
    },
    change: function(event, key) {
      const file = event.target.files[0];
      this.files.splice(key - 1, 1, {key: key, value: file});
    },
    add: function() {
      this.inputs++;
    },
    remove: function() {
      this.files.splice(-1,1);
      this.inputs--;
    },
    upload: function() {
      let formData = new FormData();
      this.files.forEach(element => {
        formData.append(element.key, element.value);
      });

      axios.post(`${domain}/upload`, formData, {headers: {'Content-Type': 'multipart/form-data'}})
        .then((res) => this.uid = res.data)
        .catch((error) => alert(error));
    }
  }
}
</script>
