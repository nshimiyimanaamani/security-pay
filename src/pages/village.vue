<template>
  <div class="village-wrapper">
    <div class="container1">
      <div class="container-title">
        <i class="fa fa-th-large"></i>
        <h1>Street Code</h1>
        <span class="fa fa-cog"></span>
      </div>
      <div class="container-body">
        <div v-for="(house,index) in houses" :key="index">
          <b-list-group-item v-if="index % 2 == 0">
            <b-card-body v-b-toggle.collapse="''+index">
              <p class>
                {{house.housename}}
                <span class="pull-right">{{house.houseNo}}</span>
              </p>
              <div class="bottom-detail">
                <span v-if="house.percentage == 100" class="completed">completed</span>
                <span
                  v-else-if="house.percentage != 100"
                  class="details"
                >{{house.amount}} /5 last months</span>
                <div id="progress-container">
                  <div
                    v-if="house.percentage == 100"
                    id="progress-bar"
                    :style="{width: house.percentage + '%', background:green}"
                  ></div>
                  <div
                    v-else-if="house.percentage<30"
                    id="progress-bar"
                    :style="{width: house.percentage + '%', background:orange}"
                  ></div>
                  <div
                    v-else-if="house.percentage>=30"
                    id="progress-bar"
                    :style="{width: house.percentage + '%', background:red}"
                  ></div>
                </div>
                <i class="fa fa-ellipsis-v"></i>
              </div>
            </b-card-body>

            <!-- Elements to collapse -->
            <b-collapse v-if="index %2 ==0" :id="'' + index">
              <b-list-group-item v-for="year in [2018 ,2019]" :key="year">
                <div class="historic-selector" v-b-toggle="''+year+index">
                  Historical {{year}}
                  <span class="fa fa-caret-up"></span>
                </div>

                <b-collapse :id="'' + year+index">
                  <b-card-text v-for="month in months" :key="month">
                    {{month}}
                    <span>{{house.amount}} Rwf</span>
                  </b-card-text>
                </b-collapse>
              </b-list-group-item>
              <b-list-group class="closeBtn">
                <i v-b-toggle="'' + index" class="fa fa-times-circle"></i>
              </b-list-group>
            </b-collapse>
          </b-list-group-item>
        </div>
      </div>
    </div>
    <div class="container2">
      <div class="container-title">
        <i class="fa fa-th-large"></i>
        <h1>Street Code</h1>
        <span class="fa fa-cog"></span>
      </div>
      <div class="container-body">
        <div v-for="(house,index) in houses" :key="index">
          <b-list-group-item v-if="index % 2 != 0">
            <b-card-body v-b-toggle.collapse="''+index">
              <p class>
                {{house.housename}}
                <span class="pull-right">{{house.houseNo}}</span>
              </p>
              <div class="bottom-detail">
                <span v-if="house.percentage == 100" class="completed">completed</span>
                <span
                  v-else-if="house.percentage != 100"
                  class="details"
                >{{house.amount}} /5 last months</span>
                <div id="progress-container">
                  <div
                    v-if="house.percentage == 100"
                    id="progress-bar"
                    :style="{width: house.percentage + '%', background:green}"
                  ></div>
                  <div
                    v-else-if="house.percentage<30"
                    id="progress-bar"
                    :style="{width: house.percentage + '%', background:orange}"
                  ></div>
                  <div
                    v-else-if="house.percentage>=30"
                    id="progress-bar"
                    :style="{width: house.percentage + '%', background:red}"
                  ></div>
                </div>
                <i class="fa fa-ellipsis-v"></i>
              </div>
            </b-card-body>

            <!-- Elements to collapse -->
            <b-collapse v-if="index %2 !=0" :id="'' + index">
              <b-list-group-item v-for="year in [2018 ,2019]" :key="year">
                <div class="historic-selector" v-b-toggle="''+year+index">
                  Historical {{year}}
                  <span class="fa fa-caret-up"></span>
                </div>

                <b-collapse :id="'' + year+index">
                  <b-card-text v-for="month in months" :key="month">
                    {{month}}
                    <span>{{house.amount}} Rwf</span>
                  </b-card-text>
                </b-collapse>
              </b-list-group-item>
              <b-list-group class="closeBtn">
                <i v-b-toggle="'' + index" class="fa fa-times-circle"></i>
              </b-list-group>
            </b-collapse>
          </b-list-group-item>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      months: [
        "January",
        "February",
        "March",
        "April",
        "May",
        "June",
        "July",
        "August",
        "September",
        "October",
        "November",
        "December"
      ],
      houses: [
        {
          housename: "house1",
          houseNo: "No 01",
          amount: "300",
          get percentage() {
            return (this.amount * 100) / 5000;
          }
        },
        {
          housename: "house2",
          houseNo: "No 02",
          amount: "1000",
          get percentage() {
            return (this.amount * 100) / 5000;
          }
        },
        {
          housename: "house3",
          houseNo: "No 03",
          amount: "2000",
          get percentage() {
            return (this.amount * 100) / 5000;
          }
        },
        {
          housename: "house4",
          houseNo: "No 04",
          amount: "4000",
          get percentage() {
            return (this.amount * 100) / 5000;
          }
        },
        {
          housename: "house5",
          houseNo: "No 05",
          amount: "5000",
          get percentage() {
            return (this.amount * 100) / 5000;
          }
        },
        {
          housename: "house6",
          houseNo: "No 06",
          amount: "3000",
          get percentage() {
            return (this.amount * 100) / 5000;
          }
        },
        {
          housename: "house7",
          houseNo: "No 07",
          amount: "4500",
          get percentage() {
            return (this.amount * 100) / 5000;
          }
        },
        {
          housename: "house8",
          houseNo: "No 08",
          amount: "1000",
          get percentage() {
            return (this.amount * 100) / 5000;
          }
        },
        {
          housename: "house9",
          houseNo: "No 09",
          amount: "4000",
          get percentage() {
            return (this.amount * 100) / 5000;
          }
        },
        {
          housename: "house10",
          houseNo: "No 10",
          amount: "5000",
          get percentage() {
            return (this.amount * 100) / 5000;
          }
        }
      ],
      green: "#50a031",
      orange: "#f0a700",
      red: "#f3573c"
    };
  }
};
</script>
<style>
@import url("../assets/css/village.css");
</style>
