<template>
  <div>
    <h1>
      <router-link v-bind:to="{ name: 'Home' }">Video Concater</router-link>
    </h1>
    <router-view></router-view>
  </div>
</template>
