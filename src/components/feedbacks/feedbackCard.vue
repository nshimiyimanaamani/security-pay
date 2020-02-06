<template>
  <b-card no-body class="feedback-card my-2 p-3">
    <b-card-title class="text-capitalize">{{feedback.title}}</b-card-title>
    <b-card-sub-title>
      <b-row class="px-3 justify-content-between font-13 text-uppercase">
        By {{feedback.display_name}} [{{feedback.creator}}]
        <span class="text-capitalize">{{date}}</span>
      </b-row>
    </b-card-sub-title>
    <hr />
    <b-card-text class="font-14 text-capitalize">{{feedback.body}}</b-card-text>
  </b-card>
</template>

<script>
export default {
  name: "feedback-card",
  props: {
    feedback: Object
  },
  data() {
    return {
      state: {
        loading: false
      }
    };
  },
  computed: {
    date() {
      const options = { year: "numeric", month: "long", day: "numeric" };
      return new Date(this.feedback.update_at).toLocaleDateString(
        "en-EN",
        options
      );
    }
  }
};
</script>

<style>
.feedback-card .card-title {
  font-size: 17px;
  font-weight: bold;
}
.feedback-card .card-text {
  font-size: 15px;
}
</style>