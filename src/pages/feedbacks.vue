<template>
  <b-container class="p-4 feedback-page" fluid>
    <vue-title title="Paypack | Feedbacks" />
    <header>
      <h4 class="mb-0">FEEDBACKS</h4>
      <b-button @click="loadData" variant="info" class="app-color">
        Refresh
        <i class="fa fa-sync-alt" :class="{'fa-spin': state.loading}" />
      </b-button>
    </header>
    <vue-load v-if="state.loading" class="secondary-font" />
    <div class="feedbacks" v-if="!state.loading">
      <feedback v-for="(feedback,index) in feedbacks" :key="index" :feedback="feedback" />
    </div>
    <b-card v-if="!state.loading && !feedbacks" class="empty-feedbacks">
      <p>No feedbacks Available at the moment!</p>
    </b-card>
  </b-container>
</template>

<script>
// import addFeedback from "../components/feedbacks/addFeedback.vue";
import feedbackCard from "../components/feedbacks/feedbackCard.vue";
export default {
  name: "feedbacks",
  components: {
    feedback: feedbackCard
  },
  data() {
    return {
      state: {
        loading: true
      },
      feedbacks: null
    };
  },

  beforeMount() {
    this.loadData();
  },
  methods: {
    async loadData() {
      this.state.loading = true;
      const Total = await this.getTotal("/feedback?offset=0&limit=0");
      this.axios
        .get("/feedback?offset=0&limit=" + Total)
        .then(res => {
          this.feedback = res.data.Messages.sort((a, b) => {
            return new Date(b.update_at) - new Date(a.update_at);
          });
        })
        .catch(err => {
          this.feedbacks = null;
          const error = err.response
            ? err.response.data.error || err.response.data
            : null;
          if (error) this.$snotify.error(error);
        })
        .finally(() => {
          this.state.loading = false;
        });
    },
    getTotal(endpoint) {
      return this.axios
        .get(endpoint)
        .then(res => res.data.Total)
        .catch(err => 0);
    }
  }
};
</script>

<style lang="scss">
.feedback-page {
  min-width: 500px;
  header {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin: 1rem 0;
    h4 {
      font-size: 1.7rem;
      letter-spacing: 1.5px;
    }
  }
  .empty-feedbacks {
    text-align: center;
    padding: 3rem;
    border-radius: 2px;
    p {
      margin: 0 !important;
      font-size: 1.2rem;
      text-transform: capitalize;
    }
  }
}
</style>