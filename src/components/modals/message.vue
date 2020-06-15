<template>
  <b-modal
    @hidden="$emit('modal-closed')"
    v-model="state.show"
    content-class="secondary-font"
    body-class="px-4 py-3 message-wrapper"
    hide-footer
  >
    <template v-slot:modal-title>
      <h4 class="m-0">Send Text Message</h4>
    </template>
    <b-row class="align-items-end" no-gutters>
      <p class="m-0">Phone Number:</p>
      <b-badge v-for="(phone,i) in phones" :key="i" class="align-items-baseline d-flex p-2 ml-2">
        <b-card-text class="m-0">{{phone}}</b-card-text>
      </b-badge>
    </b-row>
    <b-row class="mt-4" no-gutters>
      <textarea
        v-model="message"
        class="w-100 p-2 rounded"
        rows="10"
        placeholder="Enter your message..."
      ></textarea>
    </b-row>
    <b-row class="justify-content-end" no-gutters>
      <b-button @click="send" class="mt-4 br-2 px-4" variant="info" :disabled="disableBtn">
        {{ state.sending ? 'Sending' : "Send" }}
        <i
          class="fa fa-spinner fa-spin"
          v-show="state.sending"
        />
        <i class="fa fa-paper-plane" v-show="!state.sending" />
      </b-button>
    </b-row>
  </b-modal>
</template>

<script>
export default {
  name: "message-component",
  props: {
    phones: Array
  },
  data() {
    return {
      message: null,
      state: {
        sending: false,
        show: true
      }
    };
  },
  computed: {
    disableBtn() {
      if (this.message === null || this.state.sending === true) return true;
      return false;
    }
  },
  methods: {
    send() {
      this.state.sending = true;
      this.axios
        .post("/notifications/send", {
          message: this.message,
          recipients: this.phones
        })
        .then(res => {
          this.$snotify.info(`Message sent to ${this.phones[0]}`);
          this.$emit("sent");
          this.show = false;
          this.state.sending = false;
        })
        .catch(err => {
          console.log(err, err.request, err.response);
          this.$snotify.error("Failed to send message! try again later");
          this.state.sending = false;
        });
    }
  }
};
</script>

<style lang="scss">
.message-wrapper {
  .close-button {
    width: 2em;
    height: 2em;
  }
  .loading-spinner {
    position: absolute;
    top: 0;
    left: 0;
    bottom: 0;
    right: 0;
    background: #0000008a;
  }
}
</style>