<template>
  <b-container>
    <vue-title title="Paypack | Feedbacks" />
    <b-card no-body class="nav-controls">
      <b-tabs pills card vertical lazy>
        <b-tab title="FeedBacks" active @click="loadData">
          <div v-for="(feedback,index) in feedbacks" :key="index">
            <feedback :feedback="feedback" v-if="!state.loading" />
          </div>
          <div v-if="state.loading" class="text-center m-5">
            <b-spinner small></b-spinner>
            <strong>Loading...</strong>
          </div>
        </b-tab>
        <!-- <b-tab title="Send FeedBack">
          <add-feedback />
        </b-tab>-->
      </b-tabs>
    </b-card>
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
        .get(this.endpoint + "/feedback?offset=0&limit=10")
        .then(res => {
          this.feedbacks = res.data.Messages.sort((a, b) => {
            return new Date(b.update_at) - new Date(a.update_at);
          });
          this.state.loading = false;
        })
        .catch(err => {
          if (navigator.onLine) {
            const error = isNullOrUndefined(err.response)
              ? "an error occured"
              : err.response.data.error || err.response.data;
            this.$snotify.error(error);
          } else {
            this.$snotify.error("Please connect to the internet");
          }
          this.state.loading = false;
        });
    }
  }
};
</script>

<style>
</style>