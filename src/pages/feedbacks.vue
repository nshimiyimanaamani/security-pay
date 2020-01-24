<template>
  <b-container class="max-width">
    <vue-title title="Paypack | Feedbacks" />
    <div class="d-flex">
      <h4 class="mb-0">FEEDBACKS</h4>&nbsp;&nbsp;
      <b-spinner class="align-self-center" v-if="state.loading" small></b-spinner>
    </div>

    <div v-for="(feedback,index) in feedbacks" :key="index">
      <feedback :feedback="feedback" v-if="!state.loading" />
    </div>
    <b-row v-if="!state.loading" class="justify-content-start mx-1 my-4">
      <b-card class="bg-transparent align-items-center w-50">
        <b-card-text>No feedbacks Available the moment!</b-card-text>
      </b-card>
    </b-row>
  </b-container>
</template>

<script>
// import addFeedback from "../components/feedbacks/addFeedback.vue";
import feedbackCard from "../components/feedbacks/feedbackCard.vue";
export default {
  name: "feedbacks",
  components: {
    // "add-feedback": addFeedback,
    feedback: feedbackCard
  },
  data() {
    return {
      state: {
        loading: false
      },
      feedbacks: null
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    }
  },
  mounted() {
    this.loadData();
  },
  methods: {
    loadData() {
      this.state.loading = true;
      this.axios
        .get(this.endpoint + "/feedback?offset=0&limit=1000")
        .then(res => {
          this.feedbacks = res.data.Messages.sort((a, b) => {
            return new Date(b.update_at) - new Date(a.update_at);
          });
        })
        .catch(err => {
          if (navigator.onLine) {
            const error = err.response
              ? err.response.data.error || err.response.data
              : "an error occured";
            this.$snotify.error(error);
          } else {
            this.$snotify.error("Please connect to the internet");
          }
        })
        .finally(() => {
          this.state.loading = false;
        });
    }
  }
};
</script>

<style>
</style>