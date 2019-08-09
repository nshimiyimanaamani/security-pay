<template>
  <div class="table-container">
    <h4 class="title text-center">{{title}}</h4>
    <b-button @click.prevent="download" class="download btn-info">Download</b-button>
    <b-table id="data-table" bordered striped hover small :items="items" :fields="fields">
      <template slot="owner" slot-scope="data">{{ data.value.fname }} {{ data.value.lname }}</template>
    </b-table>
    <pulse-loader class="reports-loader" :loading="loading" :color="color" :size="size"></pulse-loader>
  </div>
</template>

<script>
import jsPDF from "jspdf";
import "jspdf-autotable";
export default {
  name: "reports",
  data() {
    return {
      title: "List of Users in Remera",
      loading: false,
      color: "#3db3fa",
      size: "12px",
      fields: {
        owner: {
          label: "Full name",
          sortable: true
        },
        "owner.phone": {
          label: "Phone number",
          sortable: false
        },
        sector: {
          label: "sector",
          sortable: true
        },
        cell: {
          label: "cell",
          sortable: true
        },
        village: {
          label: "village",
          sortable: true
        },
        due: {
          label: "amount",
          sortable: false
        }
      },
      items: []
    };
  },
  computed: {
    endpoint() {
      return this.$store.getters.getEndpoint;
    },
    activeSector() {
      return this.$store.getters.getActiveSector;
    }
  },
  mounted() {
    this.loadData();
  },
  methods: {
    loadData() {
      this.loading = true;
      this.axios
        .get(
          `${this.endpoint}/properties/sectors/${this.activeSector}?offset=1&limit=10`
        )
        .then(res => {
          this.items = new Array();
          for (const key in res.data.properties) {
            if (res.data.properties.hasOwnProperty(key)) {
              const element = res.data.properties[key];
              this.axios
                .get(`${this.endpoint}/properties/owners/${element.owner}`)
                .then(res => {
                  element.owner = res.data;
                  this.items = [...this.items, element];
                })
                .catch(err => {
                  console.log(err);
                });
            }
            this.loading = false;
          }
        })
        .catch(err => {
          console.log(err);
          this.loading = false;
        });
    },
    download() {
      if (this.items.length <= 0) {
        this.$snotify.error("No data available to download");
      } else if (this.items.length > 0) {
        const doc = new jsPDF();
        doc.autoTable({ html: "#data-table", useCss: true });
        doc.save(`${this.title}.pdf`);
      }
    }
  }
};
</script>

<style>
.table-container {
  padding: 40px;
}
.download {
  float: right;
  margin: 10px 0;
}
.reports-loader::after{
  display: block;
  font-weight: bold;
  content: 'Please wait...'
}
</style>
