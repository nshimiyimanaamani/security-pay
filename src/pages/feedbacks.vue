<template>
  <b-container>
    <vue-title title="Paypack | Feedbacks" />
    <b-card no-body class="nav-controls">
      <b-tabs pills card vertical lazy>
        <b-tab title="FeedBacks">hello you are going to list all feedbacks</b-tab>
        <b-tab title="Send FeedBack" active>
          <add-feedback />
        </b-tab>
      </b-tabs>
    </b-card>
  </b-container>
</template>

<script>
import addFeedback from "../components/feedbacks/addFeedback.vue";
export default {
  name: "feedbacks",
  components: {
    "add-feedback": addFeedback
  }
};
</script>

<style>
</style>