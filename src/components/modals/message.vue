<template>
  <b-modal @hidden="modalClosed" v-model="show" size="sm" class="mt-3" hide-footer>
    <template v-slot:modal-title>Send Text Message</template>
    <div class="message-wrapper px-2">
      <b-row class="align-items-baseline">
        <p>Phone Number:</p>
        <b-badge v-for="(phone,i) in phones" :key="i" class="align-items-baseline d-flex p-1 m-1">
          <b-card-text class="m-0">{{phone}}</b-card-text>
          <!-- <b-button size="sm" class="ml-2 p-1 rounded-circle bg-white close-button">
            <i class="fa fa-times text-black-50" />
          </b-button>-->
        </b-badge>
      </b-row>
      <b-row class="mt-3">
        <textarea
          v-model="message"
          class="rounded w-100"
          rows="5"
          placeholder="Enter your message..."
        ></textarea>
      </b-row>
      <b-row>
        <b-button
          @click.prevent="send"
          size="sm"
          class="mt-3 ml-auto position-relative"
          variant="info"
          :disabled="message? false : true"
        >
          Send
          <div class="loading-spinner" v-show="state.sending">
            <b-spinner variant="black" small />
          </div>
          <i class="fa fa-paper-plane" />
        </b-button>
      </b-row>
    </div>
  </b-modal>
</template>

<script>
import loader from "../loader";
export default {
  name: "message-component",
  components: {
    loader
  },
  props: {
    phones: Array
  },
  data() {
    return {
      show: true,
      message: null,
      state: {
        sending: false
      }
    };
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
          this.$snotify.info("Message Sent!");
          this.$emit("sent");
        })
        .catch(err => {
          this.$snotify.error("Message not sent!");
        })
        .finally(() => {
          this.state.sending = false;
        });
    },
    modalClosed() {
      this.$emit("modal-closed");
      this.message = null;
      this.phones = [];
      this.state.sending = false;
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