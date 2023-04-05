<template>
  <div class="message-container px-5 py-3">
    <div class="tabs-wrapper">
      <b-tabs content-class="primary-font" justified card pills active-nav-item-class="app-color text-white"
        nav-class="message-nav" @activate-tab="clear">
        <b-tab title="SECTOR" active v-if="isAdmin">
          <div class="sector-body">
            <h3 class="message-title primary-font">Send message to all payers in {{ activeSector ? activeSector :
              "Sector" }}</h3>
            <label for="textarea" class="message-label">Message:</label>
            <textarea name="message" id="message" placeholder="Write your message..." cols="30" rows="10"
              class="message-textarea" v-model="message" />
            <div class="w-100 d-flex mt-2">
              <b-button variant="info" :disabled="message ? false : true" class @click="sendToSector">
                {{ state.sending ? 'Sending ' : 'Send ' }}
                <i v-if="state.sending" class="fa fa-spinner fa-spin" />
                <i v-else class="fa fa-paper-plane" />
              </b-button>
            </div>
          </div>
        </b-tab>
        <b-tab title="CELL" :active="!isAdmin">
          <div class="cell-body">
            <h3 class="message-title primary-font">Send message to all payers in {{ select.cell ? select.cell : "Cell" }}
            </h3>
            <div class="control mb-4" v-if="isAdmin">
              <label class="message-label" for="select">Select a cell:</label>
              <b-select v-model="select.cell" :options="cellOptions">
                <template v-slot:first>
                  <option :value="null" disabled>select cell</option>
                </template>
              </b-select>
            </div>
            <label class="message-label" for="textarea">Message:</label>
            <textarea name="message" id="message" placeholder="Write your message..." cols="30" rows="10"
              class="message-textarea" v-model="message" />
            <div class="w-100 d-flex mt-2">
              <b-button variant="info" :disabled="message && select.cell ? false : true" class @click="sendToCell">
                {{ state.sending ? 'Sending ' : 'Send ' }}
                <b-spinner v-if="state.sending" variant="black" small />
                <i v-else class="fa fa-paper-plane" />
              </b-button>
            </div>
          </div>
        </b-tab>
        <b-tab title="VILLAGE">
          <div class="village-body">
            <h3 class="message-title primary-font">Send message to all payers in {{ select.village ? select.village :
              "Village" }}</h3>
            <div class="control mb-4">
              <div class="cell-select">
                <label class="message-label" for="select">Select a cell:</label>
                <b-select v-model="select.cell" :options="cellOptions">
                  <template v-slot:first>
                    <option :value="null" disabled>select cell</option>
                  </template>
                </b-select>
              </div>
              <div class="village-select">
                <label class="message-label" for="select">Select a village:</label>
                <b-select v-model="select.village" :disabled="!select.cell" :options="villageOptions">
                  <template v-slot:first>
                    <option :value="null" disabled>select village</option>
                  </template>
                </b-select>
              </div>
            </div>
            <label class="message-label" for="textarea">Message:</label>
            <textarea name="message" id="message" placeholder="Write your message..." cols="30" rows="10"
              class="message-textarea" v-model="message" />
            <div class="w-100 d-flex mt-2">
              <b-button variant="info" :disabled="message && select.village ? false : true" @click="sendToVillage">
                {{ state.sending ? 'Sending ' : 'Send ' }}
                <i v-if="state.sending" class="fa fa-spinner fa-spin" />
                <i v-else class="fa fa-paper-plane" />
              </b-button>
            </div>
          </div>
        </b-tab>
      </b-tabs>
    </div>
  </div>
</template>

<script>
import { Village } from "rwanda";
export default {
  data() {
    return {
      message: null,
      state: {
        sending: false
      },
      select: {
        cell: null,
        village: null
      }
    };
  },
  mounted() {
    if (this.isAdmin === false) {
      this.select.cell = this.activeCell;
    }
  },
  computed: {
    cellOptions() {
      const { province, district } = this.location;
      return this.$cells(province, district, this.activeSector);
    },
    villageOptions() {
      const cell = this.select.cell;
      const { province, district } = this.location;
      if (cell) {
        return this.$villages(province, district, this.activeSector, cell);
      } else {
        return [];
      }
    },
    activeSector() {
      return this.$store.getters.getActiveSector;
    },
    activeCell() {
      return this.$store.getters.getActiveCell;
    },
    location() {
      return this.$store.getters.location;
    },
    user() {
      return this.$store.getters.userDetails;
    },
    isAdmin() {
      if (this.user.role === "admin") return true;
      return false;
    }
  },
  watch: {
    "select.cell"() {
      handler: {
        this.select.village = null;
      }
    }
  },
  methods: {
    async sendToSector() {
      this.state.sending = true;
      let request = `/properties?sector=${this.activeSector}&offset=0&limit=&names=&phone=`;
      let message = this.message;
      let recipients = await this.getPhoneArray(request);
      if (recipients) {
        this.axios
          .post(`/notifications/send`, {
            message: message,
            recipients: recipients
          })
          .then(res => {
            this.clear();
            this.$snotify.info(
              `Message sent to all properties in ${this.activeSector} sector!`
            );
          })
          .catch(err => {
            console.log(err, err.response, err.request);
            this.$snotify.error("Failed to send message! Try again later");
          })
          .finally(() => {
            this.state.sending = false;
          });
      } else {
        this.state.sending = false;
        this.$snotify.error("Failed to retrieve recipients");
      }
    },
    async sendToCell() {
      this.state.sending = true;
      let cell = this.select.cell;
      let request = `/properties?cell=${cell}&offset=0&limit=&names=`;
      let message = this.message;
      let recipients = await this.getPhoneArray(request);
      if (recipients) {
        this.axios
          .post(`/notifications/send`, {
            message: message,
            recipients: recipients
          })
          .then(res => {
            this.clear();
            this.$snotify.info(
              `Message sent to all properties in ${cell} cell!`
            );
          })
          .catch(err => {
            console.log(err, err.request, err.response);
            this.$snotify.error("Failed to send message! Try again later");
          })
          .finally(() => {
            this.state.sending = false;
          });
      } else {
        this.state.sending = false;
        this.$snotify.error("Failed to retrieve recipients");
      }
    },
    async sendToVillage() {
      this.state.sending = true;
      let village = this.select.village;
      let request = `/properties?cell=${village}&offset=0&limit=&names=`;
      let message = this.message;
      let recipients = await this.getPhoneArray(request);
      if (recipients) {
        this.axios
          .post(`/notifications/send`, {
            message: message,
            recipients: recipients
          })
          .then(res => {
            this.clear();
            this.$snotify.info(
              `Message sent to all properties in ${village} village!`
            );
          })
          .catch(err => {
            console.log(err.response);
            this.$snotify.error("Failed to send message! try again later");
          })
          .finally(() => {
            this.state.sending = false;
          });
      } else {
        this.state.sending = false;
        this.$snotify.error("Failed to retrieve recipients!");
      }
    },
    async getPhoneArray(request) {
      let total = await this.$getTotal(request + "0");
      return this.axios
        .get(request + `${total}`)
        .then(res => {
          let phoneNumbers = res.data.Properties.map(item => item.owner.phone);
          phoneNumbers = phoneNumbers.filter(phone => {
            if (!/^07[2893]/.test(phone)) {
              return false;
            }
            if (phone.length !== 10) {
              return false;
            }
            return true;
          });
          return phoneNumbers;
        })
        .catch(err => {
          console.log(err, err.response, err.request);
          return null;
        });
    },
    clear() {
      this.message = null;
      this.state.sending = false;
      this.select.village = null;
      if (this.isAdmin) this.select.cell = null;
    }
  }
};
</script>

<style lang="scss">
.tabs-wrapper {
  border: 1px solid #dee2e6;
  border-radius: 3px;
  width: 100%;
  margin-top: 3rem;

  .message-nav {
    .nav-item {
      margin: 0 5px;

      .nav-link {
        border-radius: 2px;
        color: #017db3;
        font-size: 14px;

        &:hover {
          color: white !important;
          background-color: #017db3 !important;
        }
      }
    }
  }

  .message-title {
    text-align: center;
    font-size: 1.2rem;
    text-transform: uppercase;
    font-weight: bold;
    margin-bottom: 1.5rem;
    color: #3a3a3a;
  }

  .message-label {
    color: #5d5d5d;
    margin-bottom: 2px;
    font-size: 14px;
    text-transform: uppercase;
  }

  .message-textarea {
    width: 100%;
    max-width: 100%;
    min-width: 100%;
    resize: vertical;
    color: #535d5f;
    -webkit-appearance: none;
    -webkit-box-align: center;
  }

  .sector-body {
    padding: 0 2rem;

    textarea {
      border-radius: 2px;
      border: 1px solid #b9bec3;
      box-shadow: none;
      -webkit-text-fill-color: #2e3757;
      padding: 0.5rem;
    }

    button {
      border-radius: 2px;
      width: 120px;
      margin-left: auto;
      text-transform: uppercase;
      letter-spacing: 1px;
    }
  }

  .cell-body,
  .village-body {
    padding: 0 2rem;

    textarea,
    select {
      border-radius: 2px;
      border: 1px solid #b9bec3;
      box-shadow: none;
      padding: 0.5rem;
      -webkit-text-fill-color: #2e3757;
    }

    button {
      border-radius: 2px;
      width: 120px;
      margin-left: auto;
      text-transform: uppercase;
      letter-spacing: 1px;
    }
  }

  .village-body {
    .control {
      display: flex;

      &>div:nth-child(1) {
        flex: 1;
        padding-right: 0.5rem;
      }

      &>div:nth-child(2) {
        flex: 1;
        padding-left: 0.5rem;
      }
    }
  }
}
</style>