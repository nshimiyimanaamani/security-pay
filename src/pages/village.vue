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
          <div class="house" v-if="index % 2 == 0">
            <div class="top-detail">
              {{house.houseNo}}
              <span>{{house.housename}}</span>
            </div>
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
          </div>
          <div v-else></div>
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
          <div class="house" v-if="index % 2 != 0">
            <div class="top-detail">
              {{house.houseNo}}
              <span>{{house.housename}}</span>
            </div>
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
          </div>
          <div v-else></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
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
      green: "green",
      orange: "orange",
      red: "#d40d0de0"
    };
  },
  mounted() {
    console.log(this.houses.amount);
  },
  methods: {
    draw(house) {
      var elem = document.querySelector(".container-body");
      var elem2 = document.querySelector(".container2 .container-body");
      house.forEach((element, i) => {
        let percent = (element.amount * 100) / 5000;
        let visibility = "hidden";
        let color = "";
        if (percent == 100) {
          color = "green";
          visibility = "unset";
        } else if (percent >= 30) {
          color = "orange";
        } else if (percent < 30) {
          color = "#d40d0de0";
        } else {
          color = "transparent";
        }
        let innerHtml = `<div class="house" style="border: 1px solid #ebebeb;height: 85px;width: 95%;margin-left: 12px;margin-top: 10px;display: grid;padding: 6px;grid-template-rows: 40px 50px;box-shadow: 0 2px 3px #888;">
            <div class="top-detail" style="grid-row: 1;font-size: 18px;letter-spacing: 1px;">
                ${element.houseNo}
                <span style="float: right;">
                ${element.housename}
                </span>
            </div>
        <div class="bottom-detail" style="grid-row: 2;">
        <span style="float: right;padding-right:25px;color:green;visibility:${visibility};margin-top:-5px;padding-bottom:5px">completed</span>
            <div id="progress-container" style="float: left;width: 96%;border-radius: 30px;background: #ebebeb">
                <div id="progress-bar" style="width: ${percent}%;border-radius: 30px;height: 10px;background: ${color};">
                </div>
            </div>
        <i class="fa fa-ellipsis-v" style=" float: right;font-size: 20px;margin-top:-12px"></i>
        </div>
        </div>`;
        let div = document.createElement("div");
        let div2 = document.createElement("div");
        if (i % 2 == 0) {
          div.innerHTML = innerHtml;
        } else if (i % 2 != 0) {
          div2.innerHTML = innerHtml;
        }
        elem.append(div);
        elem2.append(div2);
      });
    }
  }
};
</script>

<style>
@import url("../assets/css/village.css");
</style>
