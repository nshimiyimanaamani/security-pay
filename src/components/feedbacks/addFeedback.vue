<template>
  <div>
    <b-card>
      <b-form @submit.prevent="send">
        <b-form-group label="Title:" label-for="input-1" id="input-group-1">
          <b-form-input v-model="title" required placeholder="Title..." />
        </b-form-group>
        <b-form-group label="Body:" label-for="input-2" id="input-group-2">
          <code>
            <b-textarea v-model="body" required placeholder="your Feedback..." />
          </code>
        </b-form-group>
        <b-button-group class="float-right" id="button-group-1">
          <b-button :disabled="state.sending" variant="info" type="submit">
            {{state.sending ? 'Sending' : "Send"}}
            <b-spinner v-show="state.sending" small type="grow"></b-spinner>
          </b-button>
        </b-button-group>
      </b-form>
    </b-card>
  </div>
</template>

<script>
import { isNull, isNullOrUndefined } from "util";
export default {
  name: "addFeedback",
  data() {
    return {
      state: { sending: false },
      title: null,
      body: null
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    user() {
      return this.$store.getters.userDetails;
    }
  },
  watch: {
    title() {
      handler: {
        this.title = this.title.trim();
      }
      immediate: true;
    },
    body() {
      handler: {
        this.body = this.body.trim();
      }
    }
  },
  methods: {
    send() {
      var creator;
      if (this.user.role == "dev") {
        creator = this.user.username;
      } else if (this.user.role == "min") {
        creator = this.user.telephone;
      }
      if (creator && this.title.length && this.body.length) {
        this.state.sending = true;
        this.axios
          .post("/feedback", {
            title: this.title,
            body: this.body,
            creator: creator
          })
          .then(res => {
            this.$snotify.info("FeedBack Added");
          })
          .catch(err => {
            const error = err.response
              ? err.response.data.error || err.response.data
              : null;
            if (error) this.$snotify.error(error);
          })
          .finally(() => {
            this.state.sending = false;
          });
      }
    }
  }
};
</script>

<style>
</style>