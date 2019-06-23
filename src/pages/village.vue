<template>
  <div class="village-wrapper">
    <div class="container1">
      <div class="container-title">
        <i class="fa fa-th-large"></i>
        <h1>Street Code</h1>
        <span class="fa fa-cog"></span>
      </div>
      <div class="container-body"></div>
    </div>
    <div class="container2">
      <div class="container-title">
        <i class="fa fa-th-large"></i>
        <h1>Street Code</h1>
        <span class="fa fa-cog"></span>
      </div>
      <div class="container-body"></div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      houses: [
        { housename: "house1", houseNo: "01", amount: "300" },
        { housename: "house2", houseNo: "02", amount: "1000" },
        { housename: "house3", houseNo: "03", amount: "2000" },
        { housename: "house4", houseNo: "04", amount: "4000" },
        { housename: "house5", houseNo: "05", amount: "5000" }
      ]
    };
  },
  mounted() {
    this.draw(this.houses);
  },
  methods: {
    draw(house) {
      var elem = document.querySelector(".container-body");
      var elem2 = document.querySelector(".container2 .container-body");
      house.forEach((element, i) => {
        let percent = (element.amount * 100) / 5000;
        let color = "";
        if (percent == 100) {
          color = "green";
        } else if (percent >= 30) {
          color = "orange";
        } else if (percent < 30) {
          color = "red";
        } else {
          color = "transparent";
        }
        let innerHtml = `<div class="house" style="border: 1px solid #ebebeb;height: 85px;width: 95%;margin-left: 12px;margin-top: 10px;display: grid;padding: 10px;grid-template-rows: 40px 50px;box-shadow: 0 2px 3px #888;">
            <div class="top-detail" style="grid-row: 1;font-size: 18px;letter-spacing: 1px;">
                ${element.houseNo}
                <span style="float: right;">
                ${element.housename}
                </span>
            </div>
        <div class="bottom-detail" style="grid-row: 2;">
            <div id="progress-container" style="float: left;width: 96%;border-radius: 30px;background: #ebebeb;margin-top: 11px;">
                <div id="progress-bar" style="width: ${percent}%;border-radius: 30px;height: 10px;background: ${color};">
                </div>
            </div>
        <i class="fa fa-ellipsis-v" style=" float: right;font-size: 20px;"></i>
        </div>
        </div>`;
        let div = document.createElement("div");
        let div2 = document.createElement("div");
        div.innerHTML = innerHtml;
        div2.innerHTML = innerHtml;
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
