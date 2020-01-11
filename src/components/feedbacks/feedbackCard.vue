<template>
  <b-card
    :title="feedback.title"
    :sub-title="'By '+feedback.display_name+' ['+feedback.creator+']'"
    class="feedback-card my-2"
  >
    <hr />
    <b-card-text>{{feedback.body}}</b-card-text>
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
  }
};
</script>

<style>
.feedback-card {
  text-transform: capitalize;
}
.feedback-card .card-title {
  font-size: 20px;
}
.feedback-card .card-subtitle {
  text-transform: uppercase;
  font-size: 13px;
}
.feedback-card .card-text {
  font-size: 15px;
}
</style>